package index

import (
	//"../../helpers"
	"../../structs/appCxt"
	"net/http"
	"path/filepath"
)

func Index(appContext *appCxt.AppContext, w http.ResponseWriter, r *http.Request) {
	//resp := helpers.Resp{RespObj:w}
	//
	//m := helpers.Message{"Hello world!"}
	//resp.SendMessage(m)

	path, _ := filepath.Abs("./public/index.html")
	http.ServeFile(w, r, path)
}
