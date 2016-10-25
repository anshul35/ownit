package Router

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/anshul35/ownit/API/Server"
	"github.com/anshul35/ownit/Settings/Constants"
	"github.com/gorilla/mux"
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
}

func StartServer() {
	r := mux.NewRouter()
	fmt.Println("routes table is now : ", routes)
	for _, rt := range routes {
		if rt.Method != "" {
			r.HandleFunc(rt.Path, rt.Handler).Methods(rt.Method)
		} else {
			r.HandleFunc(rt.Path, rt.Handler)
		}
	}

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

func TCPServer() {
	l, err := net.Listen("tcp", Constants.TCPHost+":"+Constants.TCPPort)
	if err != nil {
		fmt.Println("Cannot start the tcp server. : ", err)
		return
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Cannot accept to the conn request : ", err)
			continue
		}
		go Server.TCPHandler(conn)
	}
}
