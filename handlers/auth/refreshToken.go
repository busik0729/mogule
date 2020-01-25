package auth

import (
	"../../helpers"
	"../../models/devicePack"
	"../../services"
	"../../structs/appCxt"
	"net/http"
)

func RefreshToken(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	dev, err := helpers.GetDeviceValues(r)
	if err != nil {
		resp.SendBadRequest(helpers.Message{err.Error()})
		return
	}

	device, errDevice := devicePack.GetByDeviceAndRefreshToken(dev, helpers.GetRefreshTokenHeader(r))
	if errDevice != nil {
		resp.SendNotFound(helpers.Message{errDevice.Error()})
		return
	}

	device.SetTokens()
	devicePack.Update(&device)

	services.RedisSet(device.GetCacheKey(), device)
	resp.SendResponse(devicePack.GetJWTTokens(device))
}
