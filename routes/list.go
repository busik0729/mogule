package routes

import (
	"../handlers/list"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getListRoutes() router.Routes {
	return router.Routes{
		// router.Route{
		// 	"BoardGetAll",
		// 	"GET",
		// 	"/board/all",
		// 	&context.ContextedHandler{&appCxt.AppContext{}, board.All},
		// 	router.Middleware{Middle.Logger},
		// 	"",
		// 	"",
		// 	[]string{},
		// },
		// router.Route{
		// 	"BoardGetOne",
		// 	"GET",
		// 	"/board/{id}",
		// 	&context.ContextedHandler{&appCxt.AppContext{}, board.One},
		// 	router.Middleware{Middle.Logger},
		// 	"",
		// 	"",
		// 	[]string{},
		// },
		router.Route{
			"ListCreate",
			"POST",
			"/list/create",
			&context.ContextedHandler{&appCxt.AppContext{}, list.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_list_create",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"ListUpdate",
			"POST",
			"/list/update",
			&context.ContextedHandler{&appCxt.AppContext{}, list.Update},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_list_update",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"ListDelete",
			"DELETE",
			"/list/delete/{id}",
			&context.ContextedHandler{&appCxt.AppContext{}, list.Delete},
			router.Middleware{Middle.Logger, Middle.Can, Middle.Auth},
			"",
			"",
			[]string{"admin", "leadership"},
			[]string{},
		},
	}
}
