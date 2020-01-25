package middleware

import (
	"../reqData"
	"net/http"
)

type Md func(http.Handler, *reqData.RequestData) http.Handler
type Middleware []Md
