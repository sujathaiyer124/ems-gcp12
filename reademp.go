package reademployee

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sujathaiyer124/ems-gcp12/models"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/iterator"
)

func init() {
	functions.HTTP("Reademployee", ReadEmployees)
}

func ReadEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	projectID := "excellent-math-403109"
	firestoreClient, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	if firestoreClient == nil {
		log.Println("Firestore client is not initialized")
		return
	}
	var iter *firestore.DocumentIterator
	params := mux.Vars(r)
	fname := params["firstname"]
	empID := params["id"]
	lname := params["lastname"]
	emprole := params["role"]
	emp_email := params["email"]

	if fname != "" && empID == "" && lname == "" && emprole == "" && emp_email == "" {
		iter = firestoreClient.Collection("employees").Where("firstname", "==", fname).Documents(ctx)
		defer iter.Stop()
	} else if fname == "" && empID == "" && lname != "" && emprole == "" && emp_email == "" {
		iter = firestoreClient.Collection("employees").Where("lastname", "==", lname).Documents(ctx)
		defer iter.Stop()
	} else if fname == "" && empID == "" && lname == "" && emprole == "" && emp_email != "" {
		iter = firestoreClient.Collection("employees").Where("email", "==", emp_email).Documents(ctx)
		defer iter.Stop()
	} else if fname == "" && empID == "" && lname == "" && emprole != "" && emp_email == "" {
		iter = firestoreClient.Collection("employees").Where("role", "==", emprole).Documents(ctx)
		defer iter.Stop()
	} else if fname == "" && empID != "" && lname == "" && emprole == "" && emp_email == "" {
		iter = firestoreClient.Collection("employees").Where("id", "==", empID).Documents(ctx)
		defer iter.Stop()
	} else {
		iter = firestoreClient.Collection("employees").Documents(ctx)
		defer iter.Stop()
	}
	var found bool
	var emp []models.Employee
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating through documents: %v", err.Error())
		}
		var employee models.Employee
		if err != nil {
			log.Printf("Error converting document to Employee: %v", err)

		}
		if id, ok := doc.Data()["id"].(string); ok {
			employee.ID = id
		}

		employee.FirstName = doc.Data()["firstname"].(string)
		employee.LastName = doc.Data()["lastname"].(string)
		employee.Email = doc.Data()["email"].(string)
		employee.Password = doc.Data()["password"].(string)
		employee.PhoneNo = doc.Data()["phoneno"].(string)
		employee.Role = doc.Data()["role"].(string)
		employee.Salary = doc.Data()["salary"].(float64)
		emp = append(emp, employee)
	}
	jsonData, err := json.Marshal(emp)
	// if !found {
	// 	w.WriteHeader(http.StatusNotFound)
	// 	json.NewEncoder(w).Encode("Not found")
	// }
	if err != nil {
		log.Fatalf("Error encoding JSON: %s", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

}
