package handlers

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

type User struct {
	Email    string `json:"Email" validate:"required,email"`
	Name     string `json:"Name" validate:"required,min=3,max=100"`
	Password string `json:"Password" validate:"required,min=6"`
}

type AuthUser struct {
	Email    string `json:"Email" validate:"required,email"`
	Password string `json:"Password" validate:"required,min=6"`
}

type task struct {
	Name        string `json:"Name" validate:"required"`
	Description string `json:"Description" validate:"required"`
	Status      string `json:"Status" validate:"required"`
}

func Validation(d interface{}) error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(d)
	if err != nil {
		return err
	}
	return nil
}
