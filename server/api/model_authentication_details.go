package api

import "github.com/go-playground/validator"

type AuthenticationDetails struct {
	// the auth key required for the call to succeed
	AuthKey string `json:"authKey" validate:"required"`
}

func (n *AuthenticationDetails) Validate() error {
	validate := validator.New()
	return validate.Struct(n)
}

