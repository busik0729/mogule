package client

import (
	"net/http"

	"../../helpers"
	"../../models/clientPack"
	"../../structs/appCxt"
)

func All(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}

	var AllowFilter = []string{"category_id"}

	if appContext.FS.FilterBy != nil {
		exists, _ := helpers.InArray(appContext.FS.FilterBy, AllowFilter)

		if !exists {
			helpers.LogToFile(helpers.Join("Not Exists: ", appContext.FS.FilterBy.(string), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
			resp.SendBadRequest(helpers.Message{"ERROR: Плохой запрос!"})
			return
		}
	}

	ccs, err, count := clientPack.GetAll(appContext)
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{err.Error()})
		return
	}

	res := make(helpers.ResMap)
	res[ccs.GetMapKey()] = ccs.Render(appContext.CurrentUser.GetRoleString())
	res["Count"] = count

	resp.SendResponse(res)
}
