package middleware

import (
	"log"
	"net/http"
	"time"

	"../structs/appCxt"
)

func Logger(inner http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			reqData.RouteInfo.Name,
			time.Since(start),
		)
	})
}
