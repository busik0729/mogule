package appCxt

import (
	"../../models/devicePack"
	"../../models/userPack"
	"../../structs/fs"
	"../../structs/paginator"
	"../routeInfo"
)

type AppContext struct {
	CurrentUser   userPack.User
	CurrentDevice devicePack.Device
	RouteInfo     routeInfo.RouterInfo
	RequestBody   []byte
	RequestParams map[string]string
	RequestQuery  map[string]interface{}
	Paginator     paginator.Paginator
	FS            fs.FS
}
