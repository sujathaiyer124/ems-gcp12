package reademployee

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/sujathaiyer124/ems-gcp12/models"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"google.golang.org/api/iterator"
)

var firestoreClient *firestore.Client

func init() {
	functions.HTTP("Reademployee", ReadEmployees)
	ctx := context.Background()
	projectID := "excellent-math-403109"
	Client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	firestoreClient = Client

}

func ReadEmployees(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	if firestoreClient == nil {
		log.Println("Firestore client is not initialized")
		return
	}
	//firestoreClient.Close()
	iter := firestoreClient.Collection("employees").Documents(ctx)
	defer iter.Stop()

	var emp []models.Employee
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating through documents: %v", err)
			continue
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
	if err != nil {
		log.Fatalf("Error encoding JSON: %s", err.Error())
	}
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
	firestoreClient.Close()
}
