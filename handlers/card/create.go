package card

import (
	"encoding/json"
	"log"
	"net/http"

	"../../helpers"
	"../../models/cardPack"
	"../../structs/appCxt"
)

func Create(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var u cardPack.Card
	b := appContext.RequestBody
	er := json.Unmarshal(b, &u)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	u.SetUUID()
	u.SetDefault()
	_, err := u.Create()
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	m := make(helpers.ResMap)
	m[u.GetMapKey()] = u.Render(appContext.CurrentUser.GetRoleString())
	log.Println(m)

	resp.SendResponseMap(m)
}
