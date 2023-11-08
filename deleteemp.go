package deleteemployee

import (
	"context"
	"ems-gcp/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gorilla/mux"
)

var firestoreClient *firestore.Client

func init() {
	functions.HTTP("UpdateEmployee", UpdateEmployees)
	ctx := context.Background()
	projectID := "excellent-math-403109"
	Client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	firestoreClient = Client
}

func DeleteEmployee(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Delete the employee data")
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	if firestoreClient == nil {
		log.Println("Firestore client is not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update employeerr 5345."))
		return
	}

	params := mux.Vars(r)
	empID := params["id"]
	var deleteemp []models.Employee
	err := json.NewDecoder(r.Body).Decode(&deleteemp)
	if err != nil {
		fmt.Println("error is ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "Error decoding JSON")
	}

	_, err1 := firestoreClient.Collection("employees").Doc(empID).Delete(ctx)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete employee."))
		return
	}
	log.Println("Employee deleted successfully")

	docs, err := firestoreClient.Collection("employees").Documents(ctx).GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to retrieve remaining employees."))
		return
	}

	var remainingEmployees []models.Employee
	for _, doc := range docs {
		var employee models.Employee
		err := doc.DataTo(&employee)
		if err != nil {
			log.Print("Cannot get the employee data", err.Error())
		}
		remainingEmployees = append(remainingEmployees, employee)
	}
	jsonData, err := json.Marshal(remainingEmployees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
