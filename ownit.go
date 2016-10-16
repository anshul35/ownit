package main

import (
	"github.com/anshul35/ownit/Auth"
	"github.com/anshul35/ownit/Router"
	"github.com/anshul35/ownit/API"
)

func main() {
	_ = Auth.RegisterMe
	_ = API.RegisterMe

	go Router.TCPServer()
	Router.StartServer()
}
