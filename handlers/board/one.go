package board

import (
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"

	"../../helpers"
	"../../models/boardPack"
	"../../structs/appCxt"
)

func One(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	params := appContext.RequestParams
	bId := params["id"]

	bUUID, _ := uuid.FromString(bId)
	ccs, err := boardPack.GetById(&bUUID)
	log.Println(ccs)
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	res := make(helpers.ResMap)
	cR := ccs.Render(appContext.CurrentUser.GetRoleString())
	res[ccs.GetMapKey()] = cR

	resp.SendResponse(res)
}
