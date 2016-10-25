package socket;

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/anshul35/ownit/Auth/JWT"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024*10,
	WriteBufferSize: 1024*10,
	CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	user, err := JWT.AuthenticateClientRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(err.Error()))
		return
	}

	//Upgrade to a web socket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}


	client := NewClient()
	client.ClientID = user.UserID
	client.hub = hub
	client.conn = conn
	client.User = user

	client.hub.register <- client

	go client.writePump()
	client.readPump()
}