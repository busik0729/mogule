package client

import (
	"encoding/json"
	"net/http"

	"../../helpers"
	"../../models/clientPack"
	"../../structs/appCxt"
)

func Update(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var u clientPack.Client
	b := appContext.RequestBody
	er := json.Unmarshal(b, &u)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	if _, err := clientPack.Update(&u); err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	m := make(helpers.ResMap)
	m[u.GetMapKey()] = u.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponseMap(m)
}
