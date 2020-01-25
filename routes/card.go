package routes

import (
	"../handlers/card"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getCardRoutes() router.Routes {
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
			"CardCreate",
			"POST",
			"/card/create",
			&context.ContextedHandler{&appCxt.AppContext{}, card.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_card_create",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"CardUpdate",
			"POST",
			"/card/update",
			&context.ContextedHandler{&appCxt.AppContext{}, card.Update},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_card_update",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"CardDelete",
			"DELETE",
			"/card/delete/{id}",
			&context.ContextedHandler{&appCxt.AppContext{}, card.Delete},
			router.Middleware{Middle.Logger, Middle.Can, Middle.Auth},
			"",
			"",
			[]string{"admin"},
			[]string{""},
		},
	}
}
