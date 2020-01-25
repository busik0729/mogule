package routes

import (
	"../handlers/auth"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getAuthRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"RefreshToken",
			"GET",
			"/refresh-token",
			&context.ContextedHandler{&appCxt.AppContext{}, auth.RefreshToken},
			router.Middleware{Middle.Logger},
			"",
			"",
			[]string{},
			[]string{""},
		},
	}
}
