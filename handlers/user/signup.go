package user

import (
	"../../database"
	"../../helpers"
	"../../models/devicePack"
	"../../models/userPack"
	"../../services"
	"../../structs/appCxt"
	"net/http"

	"encoding/json"
)

func Signup(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var u userPack.User
	b := appContext.RequestBody
	er := json.Unmarshal(b, &u)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	u.SetDefault()
	u.CryptPassword()
	u.SetUUID()

	database.BeginTx()

	if err := u.Create(); err != nil {
		if database.IsTx() {
			database.GetConnection().Tx.Rollback()
		}
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	deviceArr, err := helpers.GetDeviceValues(r)
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}
	device := devicePack.Device{
		Platform:   deviceArr[0],
		DeviceUuid: deviceArr[1],
		Model:      deviceArr[2],
		Serial:     deviceArr[3],
		VersionOS:  deviceArr[4],
		VersionApp: deviceArr[5],
		UserId:     u.Id}
	device.SetUUID()
	device.SetTokens()

	if _, errDevice := device.Create(); errDevice != nil {
		if database.IsTx() {
			database.GetConnection().Tx.Rollback()
		}
		helpers.LogToFile(helpers.Join(errDevice.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	database.CommitTx()

	services.RedisSet(device.GetCacheKey(), device)

	m := make(helpers.ResMap)
	m[device.GetMapKey()] = devicePack.GetJWTTokens(device)
	m[u.GetMapKey()] = u.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponseMap(m)
}
