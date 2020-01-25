package user

import (
	"../../helpers"
	"../../models/userPack"
	"../../structs/appCxt"
	"net/http"
)

func Roles(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}

	ccs := userPack.GetRUSRoles()

	res := make(helpers.ResMap)
	res["Roles"] = ccs

	resp.SendResponse(res)
}
