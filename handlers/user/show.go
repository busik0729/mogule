package user

import (
	"net/http"

	"github.com/gorilla/mux"

	"../../helpers"
	"../../models/userPack"

	"github.com/satori/go.uuid"

	"../../structs/appCxt"
	"../../ws"
)

func Show(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	ws.SendMessage(appContext.CurrentDevice.WsId, "message", appContext.CurrentUser)

	vars := mux.Vars(r)
	resp := helpers.Resp{RespObj: w}

	userId, _ := uuid.FromString(vars["id"])
	u, dbErr := userPack.GetById(&userId)

	if dbErr != nil {
		helpers.LogToFile(helpers.Join(dbErr.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	m := make(helpers.ResMap)
	m[u.GetMapKey()] = u.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponse(m)
}
