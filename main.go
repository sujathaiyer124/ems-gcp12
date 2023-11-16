package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	fmt.Println("Welcome to employee management system")
	r := mux.NewRouter()
	r.HandleFunc("/openapi.yaml", OpenAPIHandler)
	r.PathPrefix("/swagger-ui/").Handler(http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./swagger-ui/dist"))))

	fmt.Println("Server  is getting started ....")
	log.Fatal(http.ListenAndServe(":8000", r))

}
func OpenAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css")
	filename := "openapi.yaml"
	//log.Println(filename)
	http.ServeFile(w, r, filename)
}
