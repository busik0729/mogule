package categoryClient

import (
	"../../helpers"
	"../../models/categoryClientPack"
	"../../structs/appCxt"
	"net/http"

	"encoding/json"
)

func Create(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var u categoryClientPack.CategoryClient
	b := appContext.RequestBody
	er := json.Unmarshal(b, &u)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	if _, err := u.Create(); err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	m := make(helpers.ResMap)
	m[u.GetMapKey()] = u.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponseMap(m)
}
