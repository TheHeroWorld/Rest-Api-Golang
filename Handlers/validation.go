package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

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
	Name        string `json:"Name"`
	Description string `json:"Description" `
	Status      string `json:"Status" `
}

func Validation(d interface{}) error {
	validate = validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(d)
	if err != nil {
		log.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Validation failed")
		return err
	}
	log.Info("Validation succeeded")
	return nil
}
