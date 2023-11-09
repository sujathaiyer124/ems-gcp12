package deleteemployee

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/sujathaiyer124/ems-gcp12/models"

	"cloud.google.com/go/firestore"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	//"github.com/gorilla/mux"
)

func init() {
	functions.HTTP("DeleteEmployee", DeleteEmployees)

}

func DeleteEmployees(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Delete the employee data")
	w.Header().Set("Content-Type", "application/json")
	//ctx := r.Context()
	ctx := context.Background()
	projectID := "excellent-math-403109"
	Client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		log.Println(err)
	}
	//documentId:kD0DLp9An4iwzEFjlGcz

	if Client == nil {
		log.Println(Client)
		log.Println("Firestore client is not initialized", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to update employeerr 5345."))
		return
	}
	log.Println("request is:", r)
	//params := mux.Vars(r)
	//log.Println("params is:", params)

	//empID := r.URL.Path[len("/DeleteEmployee/"):]
	//empID := params["id"]

	urlParts := strings.Split(r.URL.String(), "/")
	empID := urlParts[len(urlParts)-1]

	log.Println("empID is:", empID)
	log.Println("Url is", r.URL)
	//_, err = Client.Collection("employees").Doc(empID).Delete(ctx)
	docRef := Client.Collection("employees").Doc(empID)
	_, err = docRef.Delete(ctx)
	if err != nil {
		// Handle error
		log.Printf("Error deleting employee with ID %s: %s", empID, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to delete employee."))
		return
	}

	// if err != nil {
	// 	log.Printf("Error deleting employee with ID %s: %s", empID, err.Error())
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	w.Write([]byte("Failed to delete employee."))
	// 	return
	// }
	log.Println("Employee deleted successfully")

	docs, err := Client.Collection("employees").Documents(ctx).GetAll()
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
