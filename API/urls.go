package API

import (
	"github.com/anshul35/ownit/Router"
	"github.com/anshul35/ownit/API/Gadget"
	"github.com/anshul35/ownit/API/Server"
)

var RegisterMe = "API"

const basePath = "/api/v1"

func init() {
	r := Router.Route{
		Method:  "GET",
		Path:    basePath + Gadget.BasePath + "/all",
		Handler: Gadget.GadgetListHandler,
	}
	Router.RegisterRoute(r)

	r = Router.Route{
		Method: "POST",
		Path: basePath + Server.BasePath + "/register",
		Handler: Server.RegisterServerHandler,
	}
	Router.RegisterRoute(r)

	r = Router.Route{
		Method: "POST",
		Path: basePath + Server.BasePath + "/{serverID}/commands/add",
		Handler: Server.AddCommandHandler,
	}
	Router.RegisterRoute(r)

	r =Router.Route{
		Method: "GET",
		Path: basePath + Server.BasePath + "/{serverID}/commands/all",
		Handler: Server.ListCommandHandler,
	}
	Router.RegisterRoute(r)
}
