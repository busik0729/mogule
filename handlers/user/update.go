package user

import (
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"

	"../../helpers"
	"../../models/userPack"
	"../../structs/appCxt"
)

func Update(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var m map[string]interface{}

	b := appContext.RequestBody
	defer r.Body.Close()
	er := json.Unmarshal(b, &m)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Некорректные данные!"})
	}

	mid, _ := m["id"].(string)
	muuid, _ := uuid.FromString(mid)

	if !appContext.CurrentUser.IsAdmin() && !appContext.CurrentUser.IsLeadership() && muuid.String() != appContext.CurrentUser.Id.String() {
		helpers.LogToFile(helpers.Join("Forbidden", appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendForbidden()
		return
	}

	dbFromUser, err := userPack.GetById(&muuid)
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	if val, ok := m["role"]; ok {
		vals := int(val.(float64))
		bv, _ := helpers.InArray(vals, userPack.USER_ROLES_INT)
		if !bv {
			helpers.LogToFile(helpers.Join("Not Exists: role", appContext.RouteInfo.Method, appContext.RouteInfo.Name))
			resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
			return
		}

		dbFromUser.Role = vals
	}
	if val, ok := m["username"]; ok {
		username, _ := val.(string)
		dbFromUser.Username = username
	}

	if val, ok := m["name"]; ok {
		name, _ := val.(string)
		dbFromUser.Name = name
	}

	if val, ok := m["surname"]; ok {
		surname, _ := val.(string)
		dbFromUser.Surname = surname
	}

	val, ok := m["password"]
	if ok && val != "" && val != nil {
		password, _ := val.(string)
		dbFromUser.Password = password
		dbFromUser.CryptPassword()
	}

	_, errDb := userPack.Update(&dbFromUser)
	if errDb != nil {
		resp.SendBadRequest(helpers.Message{errDb.Error()})
		return
	}

	mr := make(helpers.ResMap)
	mr[dbFromUser.GetMapKey()] = dbFromUser.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponse(mr)
}
