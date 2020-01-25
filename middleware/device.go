package middleware

import (
	"../helpers"
	"../structs/appCxt"
	"net/http"
)

func Device(h http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := helpers.Resp{RespObj: w}

		deviceArr, err := helpers.GetDeviceValues(r)
		if err != nil {
			helpers.LogToFile(helpers.Join(err.Error(), reqData.RouteInfo.Method, reqData.RouteInfo.Name))
			resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
			return
		}
		if len(deviceArr) != 6 {
			helpers.LogToFile(helpers.Join("Invalid device header", reqData.RouteInfo.Method, reqData.RouteInfo.Name))
			resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
			return
		}

		h.ServeHTTP(w, r)
	})
}
