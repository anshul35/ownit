package Server;

import(
	"fmt"
	"net"
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

	//Get req if the requestID still valid?
	runReq, err := Models.GetRequestByID(requestID)
	//Not valid
	if err != nil {
		fmt.Println("No request found while reading tcp data. Error:",err)
		return
	}
	_ = runReq
	//Everything's fine. Read the command output now.
	for{
		n, err := conn.Read(buffer)

		//Either error in connection or EOF reached
		if err != nil{
			fmt.Println("Error in read : ",err)
			break
		}

		//Actual output of the command
		if (n != 0){
			//do something, for instance, 
			//send data to a message queue
			fmt.Println(string(buffer[:n]))
		}
	}
	return
}