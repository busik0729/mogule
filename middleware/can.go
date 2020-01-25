package middleware

import (
	"../helpers"
	"../structs/appCxt"
	"net/http"
)

func Can(h http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		resp := helpers.Resp{RespObj: w}
		user := reqData.CurrentUser
		routeInfo := reqData.RouteInfo

		userRole := user.GetRoleString()
		b, _ := helpers.InArray(userRole, routeInfo.CanRole)
		if !b {
			helpers.LogToFile(helpers.Join("Forbidden", reqData.RouteInfo.Method, reqData.RouteInfo.Name))
			resp.SendForbidden()
			return
		}

		h.ServeHTTP(w, r)
	})
}
