package middleware

import (
	"net/http"

	"../helpers"
	"../structs/appCxt"
)

func NotCan(h http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := helpers.Resp{RespObj: w}
		user := reqData.CurrentUser
		routeInfo := reqData.RouteInfo
		userRole := user.GetRoleString()
		b, _ := helpers.InArray(userRole, routeInfo.NotCanRole)
		if b {
			resp.SendForbidden()
			return
		}

		h.ServeHTTP(w, r)
	})
}
