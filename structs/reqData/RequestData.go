package reqData

import (
	"../../models/devicePack"
	"../../models/userPack"
	"../routeInfo"
)

type RequestData struct {
	RouteData     routeInfo.RouterInfo
	CurrentUser   userPack.User
	CurrentDevice devicePack.Device
}
