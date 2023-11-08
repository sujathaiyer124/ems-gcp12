package deleteemployee

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sujathaiyer124/ems-gcp12/models"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gorilla/mux"
)

var firestoreClient *firestore.Client

func init() {
	functions.HTTP("DeleteEmployee", DeleteEmployees)
	ctx := context.Background()
	projectID := "excellent-math-403109"
	Client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	firestoreClient = Client
}

func DeleteEmployees(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Delete the employee data")
	w.Header().Set("Content-Type", "application/json")
	//ctx := r.Context()
	if firestoreClient == nil {
		log.Println("Firestore client is not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update employeerr 5345."))
		return
	}

	params := mux.Vars(r)
	empID := params["id"]
	log.Println(empID)
	_, err := firestoreClient.Collection("employees").Doc(empID).Delete(ctx)
	if err != nil {
		log.Printf("Error deleting employee with ID %s: %s", empID, err.Error())
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
