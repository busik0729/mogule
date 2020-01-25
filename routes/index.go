package routes

import (
	"../handlers/index"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getIndexRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"Index",
			"GET",
			"/",
			&context.ContextedHandler{&appCxt.AppContext{}, index.Index},
			router.Middleware{Middle.Can, Middle.Auth},
			"",
			"",
			[]string{"admin", "client"},
			[]string{},
		},
	}
}
