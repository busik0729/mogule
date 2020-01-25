package helpers

import (
	"fmt"
	"github.com/fatih/camelcase"
	"os"
	"strings"
)

func GetBaseDir() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return pwd
}

func getSchemaFileName(name string, httpMethod string) string {
	splittedName := camelcase.Split(name)
	lowHttpMethod := strings.ToLower(httpMethod)
	slicedName := strings.Join(splittedName, "_")
	fileName := lowHttpMethod + "_" + slicedName + ".json"

	return fileName
}
