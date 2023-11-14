package swagger_function

import (
	"embed"
	"net/http"

	"github.com/gorilla/mux"
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
	r := mux.NewRouter()
	r.HandleFunc("/Swagger", SwaggerFunction)
	http.Handle("/", r)
}
