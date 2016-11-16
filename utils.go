package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"regexp"
	"strings"
)

var swaggerParamRegexp *regexp.Regexp = regexp.MustCompile("{([^}]*)}")

func addRouteHandler(engine *gin.Engine, httpMethod string, path string, methodDefinition *Route, f gin.HandlerFunc) {
	var validateFunc = func(context *gin.Context) {
		parameters := getParameters(context, methodDefinition)
		context.Set("parameters", parameters)
		context.Next()
	}

	engine.Handle(httpMethod, path, validateFunc, f)
}

func AddRoutesFromSwaggerSpec(spec *Swagger, engine *gin.Engine, routeFunctions RouteFunctions) {
	for path, pathDefinition := range spec.Paths {
		for method, methodDefinition := range pathDefinition {
			controller := methodDefinition.XSwaggerRouterController
			operationId := methodDefinition.OperationId

			ginPath := swaggerParamRegexp.ReplaceAllString(path, ":$1")

			httpMethod := strings.ToUpper(method)
			f := routeFunctions[controller+"."+operationId]
			if f != nil {
				fmt.Println("Adding route handler for ", ginPath)

				addRouteHandler(engine, httpMethod, ginPath, &methodDefinition, f)
			}
		}
	}
}

func getParameters(context *gin.Context, route *Route) map[string]interface{} {
	parameters := route.Parameters

	result := make(map[string]interface{})

	for i := 0; i < len(parameters); i++ {
		parameter := parameters[i]

		if parameter.In == "query" {
			value, found := context.GetQuery(parameter.Name)

			if found {
				result[parameter.Name] = value
			}
		} else if parameter.In == "path" {
			result[parameter.Name] = context.Param(parameter.Name)
		} else if parameter.In == "body" {
			bytes, err := ioutil.ReadAll(context.Request.Body)

			if err != nil {
			} else {
				var value interface{}

				err = json.Unmarshal(bytes, &value)

				if err != nil {
				} else {
					result[parameter.Name] = value
				}
			}
		}
	}

	return result
}
