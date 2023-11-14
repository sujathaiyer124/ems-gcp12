package swagger_function

import (
	"embed"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

var staticFiles embed.FS

func SwaggerFunction(w http.ResponseWriter, r *http.Request) {

	content, err := staticFiles.ReadFile("openapi.yaml")
	if err != nil {
		http.Error(w, "Failed to read YAML file", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	w.WriteHeader(http.StatusOK)
	w.Write(content)
}

func init() {
	functions.HTTP("Swagger", SwaggerFunction)
}
