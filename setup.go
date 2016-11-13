package swagger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func AddFromSwagger(spec *Swagger, engine *gin.Engine, routeFunctions RouteFunctions) {
	for path, pathDefinition := range spec.Paths {
		for method, methodDefinition := range pathDefinition {
			controller := methodDefinition.XSwaggerRouterController
			operationId := methodDefinition.OperationId

			httpMethod := strings.ToUpper(method)
			f := routeFunctions[controller+"."+operationId]
			if f != nil {
				fmt.Printf("Adding %s %s => %s.%s\n", method, path, controller, operationId)
				engine.Handle(httpMethod, path, f)
			}
		}
	}
}
