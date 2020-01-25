package routes

import (
	"../handlers/ws"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getWsRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"Ws",
			"GET",
			"/ws",
			&context.ContextedHandler{&appCxt.AppContext{}, ws.WsMainConnect},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.AuthByGetParam},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
	}
}
