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
	var emp []models.Employee
	fname := r.URL.Query().Get("firstname")
	lname := r.URL.Query().Get("lastname")
	empID := r.URL.Query().Get("id")
	emp_email := r.URL.Query().Get("email")
	emprole := r.URL.Query().Get("role")
	//log.Println(fname)
	if fname != "" && empID == "" && lname == "" && emprole == "" && emp_email == "" {
		iter := firestoreClient.Collection("employees").Where("firstname", "==", fname).Documents(ctx)
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				w.WriteHeader(http.StatusNotFound)
				log.Println("Firstname not found")
				json.NewEncoder(w).Encode("firstname not found")
				return
			}
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var employee models.Employee
			doc.DataTo(&employee)
			json.NewEncoder(w).Encode(employee)
			break
		}

	} else if fname == "" && empID == "" && lname != "" && emprole == "" && emp_email == "" {
		iter := firestoreClient.Collection("employees").Where("lastname", "==", lname).Documents(ctx)
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				w.WriteHeader(http.StatusNotFound)
				log.Println("Lastname not found")
				json.NewEncoder(w).Encode("Lastname not found")
				return
			}
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var employee models.Employee
			doc.DataTo(&employee)
			json.NewEncoder(w).Encode(employee)
			break
		}

	} else if fname == "" && empID == "" && lname == "" && emprole == "" && emp_email != "" {
		iter := firestoreClient.Collection("employees").Where("email", "==", emp_email).Documents(ctx)
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				w.WriteHeader(http.StatusNotFound)
				log.Println("Email not found")
				json.NewEncoder(w).Encode("Email not found")
				return
			}
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var employee models.Employee
			doc.DataTo(&employee)
			json.NewEncoder(w).Encode(employee)
			w.WriteHeader(http.StatusOK)
			break
		}

	} else if fname == "" && empID == "" && lname == "" && emprole != "" && emp_email == "" {
		var employees []models.Employee
		var found bool
		iter := firestoreClient.Collection("employees").Where("role", "==", emprole).Documents(ctx)
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var employee models.Employee
			doc.DataTo(&employee)
			employees = append(employees, employee)
			found = true
		}
		if !found {
			w.WriteHeader(http.StatusNotFound)
			log.Println("role not found")
			json.NewEncoder(w).Encode("Role not found")
			return
		}
		json.NewEncoder(w).Encode(employees)
	} else if fname == "" && empID != "" && lname == "" && emprole == "" && emp_email == "" {
		iter := firestoreClient.Collection("employees").Where("id", "==", empID).Documents(ctx)
		defer iter.Stop()

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				w.WriteHeader(http.StatusNotFound)
				log.Println("Id not found")
				json.NewEncoder(w).Encode("Id not found")
				return
			}
			if err != nil {
				log.Println("Error:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			var employee models.Employee
			doc.DataTo(&employee)
			json.NewEncoder(w).Encode(employee)
			w.WriteHeader(http.StatusOK)
			break
		}
	} else {
		iter := firestoreClient.Collection("employees").Documents(ctx)
		defer iter.Stop()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Printf("Error iterating through documents: %v", err.Error())
			}
			var employee models.Employee
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
	}

}
