package board

import (
	"../../database"
	"../../helpers"
	"../../models/boardPack"
	"../../models/listPack"
	"../../structs/appCxt"
	"../../ws"
	"net/http"

	"encoding/json"
)

func Create(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var u boardPack.Board
	b := appContext.RequestBody
	json.Unmarshal(b, &u)

	u.SetDefault()
	u.SetUUID()
	u.SetPM(appContext.CurrentUser)

	database.BeginTx()
	if _, err := u.Create(); err != nil {
		database.GetConnection().Tx.Rollback()
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	kpiList := listPack.List{Name: "KPI", BoardId: u.Id, Tracking: 1}
	if _, err := kpiList.Create(); err != nil {
		database.GetConnection().Tx.Rollback()
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
	}
	proccessList := listPack.List{Name: "В процессе", BoardId: u.Id, Tracking: 2}
	if _, err := proccessList.Create(); err != nil {
		database.GetConnection().Tx.Rollback()
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
	}
	doneList := listPack.List{Name: "Готово", BoardId: u.Id, Tracking: 3}
	if _, err := doneList.Create(); err != nil {
		database.GetConnection().Tx.Rollback()
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
	}

	database.CommitTx()

	m := make(helpers.ResMap)
	m[u.GetMapKey()] = u.Render(appContext.CurrentUser.GetRoleString())

	ws.SendBroadcast("create-board", m)

	resp.SendResponseMap(m)
}
