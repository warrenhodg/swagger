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

func addRouteHandler(engine *gin.Engine, httpMethod string, path string, methodDefinition *Route, f gin.HandlerFunc, pathRegex *regexp.Regexp) {
	var validateFunc = func(context *gin.Context) {
		parameters := getParameters(context, methodDefinition, pathRegex)
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
			pathRegex, err := regexp.Compile(swaggerParamRegexp.ReplaceAllString(path, "([^/]*)"))

			if err != nil {
				fmt.Println("Error : ", err)
			} else {
				httpMethod := strings.ToUpper(method)
				f := routeFunctions[controller+"."+operationId]
				if f != nil {
					fmt.Println("Adding route handler for ", ginPath)

					addRouteHandler(engine, httpMethod, ginPath, &methodDefinition, f, pathRegex)
				}
			}
		}
	}
}

func getParameters(context *gin.Context, route *Route, pathRegex *regexp.Regexp) map[string]interface{} {
	parameters := route.Parameters

	result := make(map[string]interface{})

	uri := context.Request.RequestURI

	matches := pathRegex.FindStringSubmatch(uri)

	fmt.Println("Matches ", matches)

	for i := 0; i < len(parameters); i++ {
		parameter := parameters[i]

		if parameter.In == "query" {
			value, found := context.GetQuery(parameter.Name)

			if found {
				result[parameter.Name] = value
			}
		} else if parameter.In == "path" {
			//TODO: Dont assume parameter order
			if i+1 < len(matches) {
				result[parameter.Name] = matches[i+1]
			}
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
