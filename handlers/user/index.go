package user

import (
	"../../helpers"
	"../../models/userPack"
	"../../structs/appCxt"
	"net/http"
)

func Index(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var AllowFilter = []string{"role"}

	if appContext.FS.FilterBy != nil {
		exists, _ := helpers.InArray(appContext.FS.FilterBy, AllowFilter)

		if !exists {
			helpers.LogToFile(helpers.Join("Not Exists", appContext.FS.FilterBy.(string), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
			resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
			return
		}
	}

	ccs, err, count := userPack.GetAllWithFS(appContext.FS, appContext.Paginator)
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	res := make(helpers.ResMap)
	res[ccs.GetMapKey()] = ccs.Render(appContext.CurrentUser.GetRoleString())
	res["Count"] = count

	resp.SendResponse(res)
}
