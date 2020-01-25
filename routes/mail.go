package routes

import (
	"../handlers/mail"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getMailRoutes() router.Routes {
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
			"MailCreate",
			"POST",
			"/mail/create",
			&context.ContextedHandler{&appCxt.AppContext{}, mail.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_mail_create",
			"",
			[]string{},
			[]string{"guest"},
		},
		// router.Route{
		// 	"BoardUpdate",
		// 	"POST",
		// 	"/board/update",
		// 	&context.ContextedHandler{&appCxt.AppContext{}, board.Update},
		// 	router.Middleware{Middle.Logger, Middle.Auth, Middle.ValidateHttpInput},
		// 	"post_board_update",
		// 	"",
		// 	[]string{},
		// },
		// router.Route{
		// 	"BoardDelete",
		// 	"DELETE",
		// 	"/board/delete/{id}",
		// 	&context.ContextedHandler{&appCxt.AppContext{}, board.Delete},
		// 	router.Middleware{Middle.Logger, Middle.Auth},
		// 	"",
		// 	"",
		// 	[]string{""},
		// },
	}
}
