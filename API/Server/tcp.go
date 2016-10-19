package Server

import (
	"fmt"
	"net"

	"github.com/anshul35/ownit/Auth/JWT"
	"github.com/anshul35/ownit/Models"
)

func TCPHandler(conn net.Conn) {
	buffer := make([]byte, 1024)

	//Read first input from the buffer,
	//which should be the requestID of the command
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Cannot read buffer")
		return
	}
	requestID := string(buffer[:n])
	fmt.Println("Response for request id ", requestID)

	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Cannot read buffer")
		return
	}
	serverID := string(buffer[:n])
	fmt.Println("Response from server id ", serverID)

	n, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Cannot read buffer")
		return
	}
	token := buffer[:n]
	fmt.Println("Response token ", token)

	server, err := Models.GetServerByID(serverID)
	if err != nil {
		fmt.Println("No server found")
		return
	}
	err = JWT.AuthenticateServerToken(server, token)
	if err != nil {
		fmt.Println("Not authenticated ", err)
		return
	}

	//Get req if the requestID still valid?
	runReq, err := Models.GetRequestByID(requestID)
	//Not valid
	if err != nil {
		fmt.Println("No request found while reading tcp data. Error:", err)
		return
	}
	_ = runReq
	//Everything's fine. Read the command output now.
	for {
		n, err := conn.Read(buffer)

		//Either error in connection or EOF reached
		if err != nil {
			fmt.Println("Error in read : ", err)
			break
		}

		//Actual output of the command
		if n != 0 {
			//do something, for instance,
			//send data to a message queue
			fmt.Println(string(buffer[:n]))
		}
	}
	return
}
