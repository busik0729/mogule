package routeInfo

type RouterInfo struct {
	Name              string
	Method            string
	Pattern           string
	SchemaPath        string
	RequestSchemaPath string
	CanRole           []string
	NotCanRole        []string
}
