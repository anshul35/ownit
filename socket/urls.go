package socket;

import (
	"fmt"
	"net/http"

	"github.com/anshul35/ownit/Router"
)

var RegisterMe = "socket"

const basePath = "/ws"

func init() {

	hub := newHub()
	go hub.run()


	r := Router.Route{
		Path:    basePath,
		Handler: func (w http.ResponseWriter, r *http.Request) {
			//Call WS handler from client.go file
			serveWs(hub, w, r)
		},
	}
	Router.RegisterRoute(r)
	fmt.Println("Registered socket")
}

