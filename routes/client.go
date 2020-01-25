package routes

import (
	"../handlers/client"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getClientRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"ClientGetAll",
			"GET",
			"/client/all",
			&context.ContextedHandler{&appCxt.AppContext{}, client.All},
			router.Middleware{Middle.Logger, Middle.FS, Middle.Paginator, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"ClientCreate",
			"POST",
			"/client/create",
			&context.ContextedHandler{&appCxt.AppContext{}, client.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_client_create",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"ClientUpdate",
			"POST",
			"/client/update",
			&context.ContextedHandler{&appCxt.AppContext{}, client.Update},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_client_update",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"ClientDelete",
			"DELETE",
			"/client/delete/{id}",
			&context.ContextedHandler{&appCxt.AppContext{}, client.Delete},
			router.Middleware{Middle.Logger, Middle.Can, Middle.Auth},
			"",
			"",
			[]string{"admin", "leadership"},
			[]string{},
		},
		router.Route{
			"ClientDeleteAll",
			"POST",
			"/client/delete/all",
			&context.ContextedHandler{&appCxt.AppContext{}, client.DeleteSelected},
			router.Middleware{Middle.Logger, Middle.Can, Middle.Auth},
			"post_client_delete_all",
			"",
			[]string{"admin", "leadership"},
			[]string{},
		},
	}
}
