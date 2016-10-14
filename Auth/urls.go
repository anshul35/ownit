package Auth

import "github.com/anshul35/ownit/Router"

var RegisterMe = "Auth"

const basePath = "/login"

func init() {
	r := Router.Route{
		Method:  "GET",
		Path:    basePath,
		Handler: HandleLogin,
	}
	Router.RegisterRoute(r)

	r = Router.Route{
		Method:  "GET",
		Path:    basePath + "/facebook",
		Handler: HandleFacebookLogin,
	}
	Router.RegisterRoute(r)

	r = Router.Route{
		Method:  "GET",
		Path:    basePath + "/facebook/successfull",
		Handler: HandleFacebookCallback,
	}
	Router.RegisterRoute(r)
}
