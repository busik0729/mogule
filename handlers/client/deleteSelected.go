package client

import (
	"net/http"

	"../../helpers"
	"../../models/clientPack"
	"../../structs/appCxt"

	"encoding/json"
)

type D struct {
	Ids clientPack.Clients `json:"ids"`
}

func DeleteSelected(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var u D
	b := appContext.RequestBody
	er := json.Unmarshal(b, &u)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	cls := u.Ids

	if _, err := clientPack.DeleteAll(&cls); err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	m := make(helpers.ResMap)
	m[cls.GetMapKey()] = cls.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponseMap(m)
}
