// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package socket;

import "fmt"

// hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client, 10),
		unregister: make(chan *Client, 10),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		fmt.Println("h:", h.register)
		select {
			case client := <-h.register:
				h.clients[client] = true
				go client.register()
			case client := <-h.unregister:
				if _, ok := h.clients[client]; ok {
					client.unregister()
					delete(h.clients, client)
				}
			}
	}
}