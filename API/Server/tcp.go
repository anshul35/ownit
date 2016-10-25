package Server

import (
	"fmt"
	"net"
	"encoding/json"
	"encoding/binary"
	"unsafe"
	"bytes"
	"io"

	"github.com/anshul35/ownit/Auth/JWT"
	"github.com/anshul35/ownit/Models"
)

func readMessage(conn net.Conn) ([]byte, error) {
	//Read # of bytes to be read
	var length int16
	b := make([]byte, int(unsafe.Sizeof(length)))
	_, err := io.ReadFull(conn, b)
	if err != nil {
		fmt.Println("Error while reading: ", err)
		return nil, err
	}

	err = binary.Read(bytes.NewBuffer(b), binary.LittleEndian, &length)
	if err != nil {
		return nil, err
	}
	
	intLength := int(length)
	b = make([]byte, intLength)
	_, err = io.ReadFull(conn, b)
	//Either error in connection or EOF reached
	if err != nil {
		fmt.Println("Error while reading: ", err)
		return nil, err
	}
	return b, nil
}

func TCPHandler(conn net.Conn) {

	defer func () {
		conn.Close()
	}()

	//Read Request Header
	b, err := readMessage(conn)
	if err != nil {
		fmt.Println("Cannot read buffer")
		return
	}

	type Header struct {
		Id string `json:"id"`
		ServerID string `json:"server_id"`
		Token []byte `json:token`
	}
	headerData := Header{}
	err = json.Unmarshal(b, &headerData)
	if err != nil {
		fmt.Println(err)
		conn.Close()
		return
	}

	requestID := headerData.Id
	serverID := headerData.ServerID
	token := headerData.Token

	//Get req if the requestID still valid?
	runReq, err := Models.GetRequestByID(requestID)
	if requestID != "PL" && requestID != "FL" {
	 	if err != nil {
			fmt.Println("No request found while reading tcp data. Error:", err)
			return
		}
	}

	//Get Senders Identity
	server, err := Models.GetServerByID(serverID)
	if err != nil {
		fmt.Println("No server found")
		return
	}

	//Do authentication
	err = JWT.AuthenticateServerToken(server, token)
	if err != nil {
		fmt.Println("Not authenticated ", err)
		return
	}

	//Everything's fine. Read the command output now.
	for {
		//Read Message
		b, err = readMessage(conn)
		if err != nil {
			if err != io.EOF{
				fmt.Println("Error in read ", err)
			}
			if requestID != "PL" && requestID != "FL" {
				close(runReq.OutChan)
			}
			return
		}

		switch requestID{
			case "PL":
					server.Processes <- b
			case "FL":
					server.FileList <- b
			default:
					runReq.OutChan <- b
		}
	}
	return
}
