package updateemployee

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/sujathaiyer124/ems-gcp12/models"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
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

func UpdateEmployees(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("Update the employee data")
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
	var updatedEmployees []models.Employee
	err := json.NewDecoder(r.Body).Decode(&updatedEmployees)
	if err != nil {
		fmt.Println("error is ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		log.Println(w, "Write which fields you want to update in body section of JSON")
		return
	}
	iter := firestoreClient.Collection("employees").Where("id", "==", empID).Documents(ctx)
	defer iter.Stop()
	for _, updatedEmployee := range updatedEmployees {
		if updatedEmployee.Salary <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Salary must be greater than 0")
			return
		}
		if !models.CustomPasswordValidation(updatedEmployee.Password) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid password. Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one digit, and one special character.")
			return
		}
		if err := models.ValidateEmployee(updatedEmployee); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if updatedEmployee.Role != "admin" && updatedEmployee.Role != "user" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Role must be admin or user")
			return
		}

		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating through documents: %v", err)
		}

		// Update with the new values
		_, updateErr := doc.Ref.Update(ctx, []firestore.Update{
			{Path: "firstname", Value: updatedEmployee.FirstName},
			{Path: "lastname", Value: updatedEmployee.LastName},
			{Path: "email", Value: updatedEmployee.Email},
			{Path: "password", Value: updatedEmployee.Password},
			{Path: "phoneno", Value: updatedEmployee.PhoneNo},
			{Path: "role", Value: updatedEmployee.Role},
			{Path: "salary", Value: updatedEmployee.Salary},
		})

		if updateErr != nil {

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to update employee data"))
			return
		}
	}
	log.Println("Employee after append:", updatedEmployees)

	jsonData, err := json.Marshal(updatedEmployees)
	if err != nil {
		log.Fatalf("Error encoding JSON: %s", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
