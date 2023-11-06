package models

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type Employee struct {
	ID        string  `json:"id,omitempty"`
	FirstName string  `json:"firstname" validate:"required,max=10,min=3"`
	LastName  string  `json:"lastname" validate:"required,max=10,min=3"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required"`
	PhoneNo   string  `json:"phoneno" validate:"required,len=10"`
	Role      string  `json:"role" validate:"required"`
	Salary    float64 `json:"salary" validate:"required,gt=0.0"`
}

func ValidateEmployee(emp Employee) error {
	err := validate.Struct(emp)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if err := validate.Struct(emp); err != nil {
		return err
	}
	if emp.Salary <= 0 {
		return errors.New("salary must be greater than 0")
	}
	return nil
}
func CustomPasswordValidation(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasDigit = true
		case strings.ContainsAny("!@#$%^&*()_+{}|:<>?", string(char)):
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}
