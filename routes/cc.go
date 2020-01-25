package routes

import (
	"../handlers/categoryClient"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getCategoryClientRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"CategoryClientGetAll",
			"GET",
			"/category-client/all",
			&context.ContextedHandler{&appCxt.AppContext{}, categoryClient.All},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"CategoryClientCreate",
			"POST",
			"/category-client/create",
			&context.ContextedHandler{&appCxt.AppContext{}, categoryClient.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_category_client_create",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"CategoryClientUpdate",
			"POST",
			"/category-client/update",
			&context.ContextedHandler{&appCxt.AppContext{}, categoryClient.Update},
			router.Middleware{Middle.Logger, Middle.Can, Middle.Auth, Middle.ValidateHttpInput},
			"post_category_client_update",
			"",
			[]string{},
			[]string{"guest"},
		},
	}
}
