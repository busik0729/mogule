package errors

import (
	"../../helpers"
	"net/http"
)

func Handle404(w http.ResponseWriter, r *http.Request) {
	resp := helpers.Resp{RespObj: w}
	resp.Send404()
}
