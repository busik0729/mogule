package router

import (
	"../appCxt"
	"../context"
	"net/http"
)

type Route struct {
	Name              string
	Method            string
	Pattern           string
	ContextedHandler  *context.ContextedHandler
	Mds               Middleware
	SchemaPath        string
	RequestSchemaPath string
	CanRole           []string
	NotCanRole        []string
}

type Routes []Route

type Md func(http.Handler, *appCxt.AppContext) http.Handler
type Middleware []Md
