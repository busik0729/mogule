package context

import (
	"../appCxt"
	"net/http"
)

type ContextedHandler struct {
	*appCxt.AppContext
	//ContextedHandlerFunc is the interface which our Handlers will implement
	ContextedHandlerFunc func(*appCxt.AppContext, http.ResponseWriter, *http.Request)
}

func (handler ContextedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.ContextedHandlerFunc(handler.AppContext, w, r)
}
