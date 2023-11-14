package createemployee

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sujathaiyer124/ems-gcp12/models"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/google/uuid"
	//"google.golang.org/genproto/googleapis/cloud/functions/v2"
)

var firestoreClient *firestore.Client

func init() {
	functions.HTTP("CreateEmployee", CreateEmployees)
	ctx := context.Background()
	projectID := "excellent-math-403109"
	//credentialsFile := "excellent-math-403109-3a5f8960ab6f.json"
	Client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	firestoreClient = Client

}
func generateShortID() string {
	id := uuid.New()
	return id.String()[:6]
}

var err error

func CreateEmployees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")
	w.Header().Set("Content-Type", "application/json")
	var employee []models.Employee
	ctx := r.Context()

	var createemp []models.Employee
	err = json.NewDecoder(r.Body).Decode(&createemp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode("No data inside Json")
		return
	}

	for _, emp := range createemp {
		if emp.Salary <= 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Salary must be greater than 0")
			return
		}
		if !models.CustomPasswordValidation(emp.Password) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid password. Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one digit, and one special character.")
			return
		}
		if err := models.ValidateEmployee(emp); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if emp.Role != "admin" && emp.Role != "user" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Role must be admin or user")
			return
		}
		id := generateShortID()
		_, _, err := firestoreClient.Collection("employees").Add(ctx, map[string]interface{}{
			"id":        id,
			"firstname": emp.FirstName,
			"lastname":  emp.LastName,
			"email":     emp.Email,
			"password":  emp.Password,
			"phoneno":   emp.PhoneNo,
			"role":      emp.Role,
			"salary":    emp.Salary,
		})
		if err != nil {
			log.Fatalf("Failed to add employee to Firestore: %v", err)
		}

	}
	createemp = append(employee, createemp...)

	//fmt.Fprintln(w, "Created employees")

	log.Println("Employee Created")
	jsonData, err := json.Marshal(createemp)
	if err != nil {
		log.Fatalf("Error encoding JSON: %s", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	defer firestoreClient.Close()
}
