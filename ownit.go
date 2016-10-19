package main

import (
	"flag"

	"github.com/anshul35/ownit/API"
	"github.com/anshul35/ownit/Auth"
	"github.com/anshul35/ownit/Router"
)

func init() {
	flag.Parse()
}

func main() {
	_ = Auth.RegisterMe
	_ = API.RegisterMe

	go Router.TCPServer()
	Router.StartServer()
}
