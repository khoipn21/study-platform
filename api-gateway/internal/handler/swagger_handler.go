package handler

import (
	"net/http"
	httpSwagger "github.com/swaggo/http-swagger"
)

type SwaggerHandler struct{}

func NewSwaggerHandler() *SwaggerHandler {
	return &SwaggerHandler{}
}

// ServeSwaggerUI serves the Swagger UI
func (h *SwaggerHandler) ServeSwaggerUI() http.Handler {
	return httpSwagger.Handler(
		httpSwagger.URL("/api/v1/docs/swagger.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	)
}

// SwaggerJSON serves the swagger.json file
func (h *SwaggerHandler) SwaggerJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// This will be served by the generated docs/swagger.json file
	http.ServeFile(w, r, "docs/swagger.json")
}
