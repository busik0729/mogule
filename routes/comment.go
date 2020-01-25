package routes

import (
	"../handlers/comment"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getCommentRoutes() router.Routes {
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
			"CommentCreate",
			"POST",
			"/comment/create",
			&context.ContextedHandler{&appCxt.AppContext{}, comment.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_comment_create",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"CommentUpdate",
			"POST",
			"/comment/update",
			&context.ContextedHandler{&appCxt.AppContext{}, comment.Update},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_comment_update",
			"",
			[]string{},
			[]string{"guest"},
		},
	}
}
