package main

import (
	"flag"

	_ "net/http/pprof"

	"github.com/anshul35/ownit/API"
	"github.com/anshul35/ownit/Auth"
	"github.com/anshul35/ownit/Router"
	"github.com/anshul35/ownit/socket"
)

func init() {
	flag.Parse()
}

func main() {
	_ = Auth.RegisterMe
	_ = API.RegisterMe
	_ = socket.RegisterMe

	go Router.TCPServer()
	Router.StartServer()
}
