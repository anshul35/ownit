package Router

import (
	"fmt"
	"log"
	"net/http"

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
	fmt.Println("routes table is now : ", routes)
}

func Server() {
	r := mux.NewRouter()
	//    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./templates/static/"))))

	for _, rt := range routes {
		r.HandleFunc(rt.Path, rt.Handler).Methods(rt.Method)
	}

	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":9090", nil))
}
