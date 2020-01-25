package user

import (
	"encoding/json"
	"net/http"

	"../../database"
	"../../helpers"
	"../../models/userPack"
	"../../structs/appCxt"
)

func Create(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var u userPack.User
	b := appContext.RequestBody
	er := json.Unmarshal(b, &u)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	rol := u.Role

	u.SetDefault()
	u.CryptPassword()
	u.SetUUID()
	u.Role = rol

	database.BeginTx()

	if err := u.Create(); err != nil {
		if database.IsTx() {
			database.GetConnection().Tx.Rollback()
		}
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	database.CommitTx()

	m := make(helpers.ResMap)
	m[u.GetMapKey()] = u.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponseMap(m)
}
