package routes

import (
	"../handlers/label"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getLabelRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"LabelCreate",
			"POST",
			"/label/create",
			&context.ContextedHandler{&appCxt.AppContext{}, label.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_label_create",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"LabelUpdate",
			"POST",
			"/label/update",
			&context.ContextedHandler{&appCxt.AppContext{}, label.Update},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_label_update",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"ListDelete",
			"DELETE",
			"/label/delete/{id}",
			&context.ContextedHandler{&appCxt.AppContext{}, label.Delete},
			router.Middleware{Middle.Logger, Middle.Can, Middle.Auth},
			"",
			"",
			[]string{"admin", "leadership"},
			[]string{""},
		},
	}
}
