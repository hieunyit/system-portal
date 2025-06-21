package docs

import "github.com/swaggo/swag"

const swaggerDoc = `{
  "swagger": "2.0",
  "info": {
    "title": "System Portal API",
    "version": "1.0"
  },
  "paths": {}
}`

func init() {
	swag.Register("swagger", &swag.Spec{
		Version:          "1.0",
		Host:             "localhost",
		BasePath:         "/",
		Schemes:          []string{},
		Title:            "System Portal API",
		Description:      "Auto-generated docs",
		InfoInstanceName: "swagger",
		SwaggerTemplate:  swaggerDoc,
	})
}
