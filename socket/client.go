// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket;

import (
	"fmt"
	"log"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"

	"github.com/anshul35/ownit/Models"
	"github.com/anshul35/ownit/Settings/Constants"

)

const (
	// Time allowed to write a message to the peer.
	writeWait = 50 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type MsgQueue struct {
	CO chan *Message
	PL chan *Message
	FL chan *Message
}

// Client is an middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	//The websocket connection.
	conn *websocket.Conn

	//User who initiated the request
	User *Models.User

	//Identification for client, available after auth
	ClientID string

	//Final Message sending queue to the browser client
	send chan *Message

	//Outbound Msg Queue
	SendQueue MsgQueue

	//Map of pending requests
	pending map[string]*Models.RunCommandRequest

	//Stop all child goroutines when true
	stop bool
}

func NewClient() *Client {
	client := Client{ 
		SendQueue: MsgQueue{
			CO:make(chan *Message, Constants.WSChannelSize),
			PL:make(chan *Message, Constants.WSChannelSize),
			FL:make(chan *Message, Constants.WSChannelSize)},
		stop: false,
		pending: make(map[string]*Models.RunCommandRequest),
		send: make(chan *Message, 100),
	}

	go client.ListenSendQueue()
	go client.ListenServers()
	go client.ListenRequests()

	return &client
}


//This will be used to add priority algorithm later
func (c *Client) ListenSendQueue() {
	for{
		if c.stop {
			break
		}
		select{
			case msg := <-c.SendQueue.PL:
				c.send<-msg
			case msg := <-c.SendQueue.FL:
				c.send<-msg
			case msg := <-c.SendQueue.CO:
				c.send<-msg
		}
	}
}


//Return WS message, corresponding raw message
//and error, if any
func (c *Client) AbsorbChannelIntoWSMessage(msg []byte, ch chan []byte, messageType string, serverID string) (*Message, []byte){
	n := len(ch)
	for i:=0;i<n;i++ {
		newM := <-ch
		msg = append(msg, newM...)
	}
	m := c.NewWSMessage(serverID, messageType, msg)
	return m, msg
}

func (c *Client) ListenServers() {
	for{
		for _, server := range c.User.Servers {
			server.Start()
			if c.stop {
				return
			}
			select{
			case msg := <-server.Processes:
				m, _:= c.AbsorbChannelIntoWSMessage(msg, server.Processes, TypeProcessList, server.ServerID)
				c.SendQueue.PL <- m
			case msg := <-server.FileList:
				m, _:= c.AbsorbChannelIntoWSMessage(msg, server.FileList, TypeFileList, server.ServerID)
				c.SendQueue.FL <- m
			default:
				//To enable non blocking on this server
				continue
			}
		}
	}
}

func (c *Client) ListenRequests() {
	for{
		for _, req := range c.pending{
			if c.stop {
				return
			}
			select {
			case msg, ok := <-req.OutChan:
				if !ok {
					//Request output is over, channel is closed 
					//by sender unregister it now
					c.unregisterRequest(req)
				} else {
					m, msg := c.AbsorbChannelIntoWSMessage(msg, req.OutChan, TypeCommandOutput, req.Server.ServerID)
					c.SendQueue.CO <- m
					req.Output = append(req.Output, msg...)
				}
			default:
				//for non blocking
				continue
			}
		}
	}
}

func (c *Client) registerRequest(r *Models.RunCommandRequest) {
	c.pending[r.RequestID] = r
	return
}

func (c *Client) unregisterRequest(r *Models.RunCommandRequest) {
	if _, ok := c.pending[r.RequestID]; ok {
		delete(c.pending, r.RequestID)
	}
	return
}

func (c *Client) NewWSMessage(serverID string, messageType string, body []byte) *Message {
	m := Message{
		Type:messageType,
		Body:body,
		ServerID:serverID,
		isInbound:false,
		client:c}
	return &m
}

func (c *Client) DecodeWSMessage(message []byte) *Message {
	var m Message
	err := json.Unmarshal(message, &m)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	m.isInbound = true
	m.client = c
	return &m
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {

	//Un-register client once there is nothing to read
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(
		func(string) error { 
			c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			//Break causes unregistry of client
			return
		}
		m := c.DecodeWSMessage(message)
		go m.Execute(false)
	}
}

// write writes a message with the given message type and payload.
func (c *Client) write(mt int, payload []byte) error {
	c.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return c.conn.WriteMessage(mt, payload)
}


// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		c.hub.unregister <- c
		ticker.Stop()
		c.conn.Close()
	}()
	for{
		if c.stop {
			break
		}
		select {
		case msg := <-c.send:
			encodedMsg, err := msg.encode()
			if err != nil {
				fmt.Println(err)
				return
			}
			err = c.write(websocket.TextMessage, encodedMsg)
			if err != nil{
				fmt.Println("Error while writing text message",err)
				return
			}
	
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				c.hub.unregister<-c
				fmt.Println(err)
				return
			}

		}
	}
}

//Init process of client creation.
func (c *Client) register() {
	for _, serv := range c.User.Servers {
		serv.Start()
	}
	return
}

func (c *Client) unregister() {
	for _, serv := range c.User.Servers {
		serv.Stop()
	}
	c.stop = true
	return
}