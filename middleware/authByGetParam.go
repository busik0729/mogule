package middleware

import (
	"net/http"

	"../helpers"
	"../models/devicePack"
	"../models/userPack"
	"../structs/appCxt"
)

func AuthByGetParam(h http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := helpers.GetAccessTokenByGetParam(r)

		resp := helpers.Resp{RespObj: w}

		if len(accessToken) == 0 {
			helpers.LogToFile(helpers.Join("Invalid access token length", reqData.RouteInfo.Method, reqData.RouteInfo.Name))
			resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
			return
		}

		_, err := helpers.ValidateAccessToken(accessToken)
		if err != nil {
			helpers.LogToFile(helpers.Join(err.Error(), reqData.RouteInfo.Method, reqData.RouteInfo.Name))
			resp.SendUnauthorized()
			return
		}

		device := devicePack.Device{}

		deviceFromDB, errToken := devicePack.GetByAccessToken(accessToken)
		if errToken != nil {
			helpers.LogToFile(helpers.Join(errToken.Error(), reqData.RouteInfo.Method, reqData.RouteInfo.Name))
			resp.SendForbidden()
			return
		}
		device = deviceFromDB

		userFromDB, errUser := userPack.GetById(device.UserId)
		if errUser != nil {
			helpers.LogToFile(helpers.Join(errUser.Error(), reqData.RouteInfo.Method, reqData.RouteInfo.Name))
			resp.SendForbidden()
			return
		}

		reqData.CurrentUser = userFromDB
		reqData.CurrentDevice = device

		h.ServeHTTP(w, r)
	})
}
