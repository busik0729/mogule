package routes

import (
	"../handlers/user"
	Middle "../middleware"
	"../structs/appCxt"
	"../structs/context"
	"../structs/router"
)

func getUserRoutes() router.Routes {
	return router.Routes{
		router.Route{
			"UserIndex",
			"GET",
			"/users",
			&context.ContextedHandler{&appCxt.AppContext{}, user.Index},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.FS, Middle.Paginator},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"UserRoles",
			"GET",
			"/users/roles",
			&context.ContextedHandler{&appCxt.AppContext{}, user.Roles},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"UserAll",
			"GET",
			"/users/all",
			&context.ContextedHandler{&appCxt.AppContext{}, user.All},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"UserShow",
			"GET",
			"/users/{id}",
			&context.ContextedHandler{&appCxt.AppContext{}, user.Show},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth},
			"",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"UserSignUp",
			"POST",
			"/user/signup",
			&context.ContextedHandler{&appCxt.AppContext{}, user.Signup},
			router.Middleware{Middle.Logger},
			"post_user_create",
			"",
			[]string{},
			[]string{},
		},
		router.Route{
			"UserLogin",
			"POST",
			"/user/login",
			&context.ContextedHandler{&appCxt.AppContext{}, user.Login},
			router.Middleware{Middle.Logger},
			"post_user_login",
			"",
			[]string{},
			[]string{},
		},
		router.Route{
			"UserUpdate",
			"POST",
			"/users/update",
			&context.ContextedHandler{&appCxt.AppContext{}, user.Update},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"put_user_update",
			"",
			[]string{},
			[]string{"guest"},
		},
		router.Route{
			"UserCreate",
			"POST",
			"/users/create",
			&context.ContextedHandler{&appCxt.AppContext{}, user.Create},
			router.Middleware{Middle.Logger, Middle.NotCan, Middle.Auth, Middle.ValidateHttpInput},
			"post_user_create",
			"",
			[]string{"admin", "leadership"},
			[]string{},
		},
	}
}
