package middleware

import (
	"../helpers"
	"../structs/appCxt"
	"../structs/fs"
	"net/http"
)

func FS(h http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs := fs.FS{}
		resp := helpers.Resp{RespObj: w}

		filterBy := reqData.RequestQuery["filterBy"]
		val := reqData.RequestQuery["v"]
		sort := reqData.RequestQuery["sort"]

		if filterBy != nil {
			if val == nil {
				helpers.LogToFile(helpers.Join("Invalid filter params", reqData.RouteInfo.Method, reqData.RouteInfo.Name))
				resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
				return
			}

			fs.FilterBy = filterBy
			fs.FilterValue = val
		}

		if sort != nil {
			s := sort.(string)
			orientir := s[:1]
			if orientir == "-" {
				fs.SortBy = s[1:]
				fs.OrderBy = "DESC"
			} else {
				fs.SortBy = s
				fs.OrderBy = "ASC"
			}
		}

		reqData.FS = fs

		h.ServeHTTP(w, r)
	})
}
