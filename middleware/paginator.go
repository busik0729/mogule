package middleware

import (
	"net/http"
	"strconv"

	"../structs/appCxt"
	"../structs/paginator"
)

func Paginator(h http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pg := paginator.Paginator{}

		limit := reqData.RequestQuery["limit"]
		page := reqData.RequestQuery["page"]

		if limit != nil {
			pg.Limit, _ = strconv.Atoi(limit.(string))
		} else {
			pg.Limit = 8
		}

		if page != nil {
			pg.Page, _ = strconv.Atoi(page.(string))
		} else {
			pg.Page = 1
		}

		pg.Offset = (pg.Page - 1) * pg.Limit.(int)

		reqData.Paginator = pg

		h.ServeHTTP(w, r)
	})
}
