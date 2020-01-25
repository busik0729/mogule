package middleware

import (
	"fmt"
	"github.com/xeipuuv/gojsonschema"
	"log"
	"net/http"

	"../helpers"
	"../structs/appCxt"
)

const HTTP_SCHEMAS_BASE_PATH = "/home/busik/web/crm/schemas/httpInputSchemas/generated/"

func ValidateHttpInput(inner http.Handler, reqData *appCxt.AppContext) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if reqData.RouteInfo.SchemaPath != "" {
			filePath := HTTP_SCHEMAS_BASE_PATH + reqData.RouteInfo.SchemaPath + ".json"
			// filePath :=  helpers.GetBaseDir() + HTTP_SCHEMAS_BASE_PATH + reqData.RouteInfo.SchemaPath + ".json"
			log.Println(filePath)
			b := reqData.RequestBody

			schemaLoader := gojsonschema.NewReferenceLoader("file://" + filePath)
			documentLoader := gojsonschema.NewStringLoader(string(b[:]))

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

				resp := helpers.Resp{RespObj: w}
				helpers.LogToFile(helpers.Join(error, reqData.RouteInfo.Method, reqData.RouteInfo.Name))
				resp.SendBadRequest(helpers.Message{"ERROR: Ошибка сервера!"})
				return
			}
		}

		inner.ServeHTTP(w, r)
	})
}
