package user

import (
	"../../helpers"
	"../../models/devicePack"
	"../../models/userPack"
	"../../services"
	"../../structs/appCxt"
	"net/http"

	"encoding/json"
)

func Login(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	var m userPack.User

	b := appContext.RequestBody
	er := json.Unmarshal(b, &m)
	if er != nil {
		helpers.LogToFile(helpers.Join(er.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
		return
	}

	user, err := userPack.GetByUsername(m.Username)
	if err != nil {
		helpers.LogToFile(helpers.Join(err.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendNotFound(helpers.Message{"ERROR: Некорректный логин/пароль!"})
		return
	}

	if !user.CheckHashPassword(m.Password) {
		helpers.LogToFile(helpers.Join("Incorrect password", appContext.RouteInfo.Method, appContext.RouteInfo.Name))
		resp.SendBadRequest(helpers.Message{"ERROR: Некорректный логин.пароль!"})
		return
	}

	deviceArr, err := helpers.GetDeviceValues(r)
	if err != nil {
		resp.SendBadRequest(helpers.Message{err.Error()})
		return
	}
	device, errDevice := devicePack.GetByHeader(deviceArr, user.Id)
	if errDevice != nil {
		device = devicePack.Device{
			Platform:   deviceArr[0],
			DeviceUuid: deviceArr[1],
			Model:      deviceArr[2],
			Serial:     deviceArr[3],
			VersionOS:  deviceArr[4],
			VersionApp: deviceArr[5],
			UserId:     user.Id}
		device.SetUUID()
		device.SetTokens()

		if _, errDevice := device.Create(); errDevice != nil {
			helpers.LogToFile(helpers.Join(errDevice.Error(), appContext.RouteInfo.Method, appContext.RouteInfo.Name))
			resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
			return
		}
	} else {
		device.SetTokens()
		devicePack.Update(&device)
	}

	services.RedisSet(device.GetCacheKey(), device)
	services.RedisSet(user.GetCacheKey(), user)

	res := make(helpers.ResMap)
	res[device.GetMapKey()] = devicePack.GetJWTTokens(device)
	res[user.GetMapKey()] = user.Render(appContext.CurrentUser.GetRoleString())

	resp.SendResponse(res)
}
