package swagger

import (
	"github.com/gin-gonic/gin"
)

type SwaggerInfo struct {
	Title       string
	Description string
	Version     string
}

type SecurityDefinition struct {
	Type        string
	Description string
}

type SecurityDefinitions map[string]SecurityDefinition

type Response map[string]interface{}

type Responses map[string]Response

type Parameter struct {
	Name     string
	In       string
	Type     string
	Required bool
}

type Route struct {
	XSwaggerRouterController string `json:"x-swagger-router-controller"`
	OperationId              string
	Tags                     []string
	Description              string
	Security                 []interface{}
	Parameters               []Parameter
	Responses                Responses
}

type Routes map[string]Route

type Paths map[string]Routes

type Definition map[string]interface{}

type Definitions map[string]Definition

type Swagger struct {
	Swagger             string
	Info                SwaggerInfo
	Produces            []string
	Host                string
	BasePath            string
	SecurityDefinitions SecurityDefinition
	Paths               Paths
	Definitions         Definitions
}

type RouteFunctions map[string]gin.HandlerFunc
