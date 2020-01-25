package routes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	errorsH "../handlers/errors"
	"../helpers"
	"../models/devicePack"
	"../models/userPack"
	"../structs/appCxt"
	"../structs/routeInfo"
	"../structs/router"

	"github.com/gorilla/mux"
	"github.com/xeipuuv/gojsonschema"
)

const HTTP_REQUEST_SCHEMAS_BASE_PATH = "/schemas/httpInputSchemas/generated/request/"

func getRoutes() router.Routes {
	routes := getIndexRoutes()
	routes = append(routes, getAuthRoutes()...)
	routes = append(routes, getUserRoutes()...)
	routes = append(routes, getClientRoutes()...)
	routes = append(routes, getCategoryClientRoutes()...)
	routes = append(routes, getBoardRoutes()...)
	routes = append(routes, getListRoutes()...)
	routes = append(routes, getCardRoutes()...)
	routes = append(routes, getCommentRoutes()...)
	routes = append(routes, getLabelRoutes()...)
	routes = append(routes, getMailRoutes()...)
	routes = append(routes, getWsRoutes()...)

	return routes
}
func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Content-Type", "application/json; charset=UTF-8")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, X-AT, X-RT, Content-Length, X-Device")
}

func parseRequest(h http.Handler, appContext *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		resp := helpers.Resp{RespObj: w}
		if r.Method == "OPTIONS" {
			resp.SendResponse(helpers.Message{})
		}

		u, err := url.Parse(r.URL.String())
		if err != nil {
			panic(err)
		}

		m := helpers.ParseQuery(u.RawQuery)
		if appContext.RouteInfo.RequestSchemaPath != "" {
			filePath := helpers.GetBaseDir() + HTTP_REQUEST_SCHEMAS_BASE_PATH + appContext.RouteInfo.RequestSchemaPath + ".json"

			schemaLoader := gojsonschema.NewReferenceLoader("file://" + filePath)
			documentLoader := gojsonschema.NewGoLoader(m)

			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				panic(err.Error())
			}

			if result.Valid() {
				fmt.Printf("The document is valid\n")
			} else {
				error := ""
				fmt.Printf("The document is not valid. see errors :\n")
				for _, desc := range result.Errors() {
					fmt.Printf("- %s\n", desc)
					error = error + desc.String() + "\n"
				}

				resp.SendBadRequest(helpers.Message{error})
				return
			}
		}
		appContext.RequestQuery = m

		vars := mux.Vars(r)
		appContext.RequestParams = vars

		b, _ := ioutil.ReadAll(r.Body)
		appContext.RequestBody = b

		h.ServeHTTP(w, r)
	})
}

func NewRouter() *mux.Router {

	routerMain := mux.NewRouter().StrictSlash(true)
	for _, route := range getRoutes() {
		var handler http.Handler
		routeInf := routeInfo.RouterInfo{route.Name, route.Method, route.Pattern, route.SchemaPath, route.RequestSchemaPath, route.CanRole, route.NotCanRole}

		route.ContextedHandler.AppContext.RouteInfo = routeInf
		route.ContextedHandler.AppContext.CurrentUser = userPack.User{}
		route.ContextedHandler.AppContext.CurrentUser.SetRole("guest")
		route.ContextedHandler.AppContext.CurrentDevice = devicePack.Device{}

		handler = route.ContextedHandler

		for _, middle := range route.Mds {
			handler = middle(handler, route.ContextedHandler.AppContext)
		}

		handler = parseRequest(handler, route.ContextedHandler.AppContext)

		routerMain.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}
	routerMain.NotFoundHandler = http.HandlerFunc(errorsH.Handle404)

	return routerMain
}
