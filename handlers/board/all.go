package board

import (
	"net/http"

	"../../helpers"
	"../../models/boardPack"
	"../../structs/appCxt"
)

func All(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}

	ccs, err := boardPack.GetAll()
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Проекты не найдены!"})
		return
	}

	res := make(helpers.ResMap)
	res[ccs.GetMapKey()] = ccs.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponse(res)
}
