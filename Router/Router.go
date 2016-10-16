package Router

import (
	"fmt"
	"log"
	"net/http"
	"net"

	"github.com/gorilla/mux"
	"github.com/anshul35/ownit/Settings/Constants"
	"github.com/anshul35/ownit/API/Server"
)

type Route struct {
	Method  string
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
}

var routes = make([]Route, 0)

func RegisterRoute(r Route) {
	fmt.Println("Registering Route : ", r)
	routes = append(routes, r)
	fmt.Println("routes table is now : ", routes)
}

func StartServer() {
	r := mux.NewRouter()
	for _, rt := range routes {
		r.HandleFunc(rt.Path, rt.Handler).Methods(rt.Method)
	}

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func TCPServer() {
	l, err := net.Listen("tcp", Constants.TCPHost + ":" + Constants.TCPPort)
	if err != nil {
		fmt.Println("Cannot start the tcp server. : ", err)
		return
	}

	defer l.Close()

	for{
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Cannot accept to the conn request : ", err)
			continue
		}
		go Server.TCPHandler(conn)
	}
}