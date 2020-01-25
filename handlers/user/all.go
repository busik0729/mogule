package user

import (
	"../../helpers"
	"../../models/userPack"
	"../../structs/appCxt"
	"net/http"
)

func All(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}

	ccs, err := userPack.GetAll()
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	res := make(helpers.ResMap)
	res[ccs.GetMapKey()] = ccs.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponse(res)
}
