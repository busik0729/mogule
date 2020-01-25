package routes

import (
	"../handlers/board"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getBoardRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"BoardGetAll",
			"GET",
			"/board/all",
			&context.ContextedHandler{&appCxt.AppContext{}, board.All},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"BoardGetOne",
			"GET",
			"/board/{id}",
			&context.ContextedHandler{&appCxt.AppContext{}, board.One},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"BoardCreate",
			"POST",
			"/board/create",
			&context.ContextedHandler{&appCxt.AppContext{}, board.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_board_create",
			"",
			[]string{""},
			[]string{"guest"},
		},
		router.Route{
			"BoardUpdate",
			"POST",
			"/board/update",
			&context.ContextedHandler{&appCxt.AppContext{}, board.Update},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_board_update",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"BoardDelete",
			"DELETE",
			"/board/delete/{id}",
			&context.ContextedHandler{&appCxt.AppContext{}, board.Delete},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{"admin", "leadership"},
			[]string{""},
		},
	}
}
