package list

import (
	uuid "github.com/satori/go.uuid"

	"../../helpers"
	"../../models/listPack"
	"../../structs/appCxt"
	"net/http"
)

func Delete(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	params := appContext.RequestParams
	bId := params["id"]

	bUUID, _ := uuid.FromString(bId)
	ccs, err := listPack.Delete(&bUUID)
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	res := make(helpers.ResMap)
	res[ccs.GetMapKey()] = ccs.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponse(res)
}
