package middleware

import (
	"net/http"

	"../helpers"
	"../structs/appCxt"
)

func AuthorizeAction(inner http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !reqData.CurrentUser.Can(reqData.RouteInfo.CanRole) {
			resp := helpers.Resp{RespObj: w}
			resp.SendForbidden()
		}

		inner.ServeHTTP(w, r)
	})
}
