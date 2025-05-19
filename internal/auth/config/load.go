package config

import "github.com/go-playground/validator/v10"

func (a *AuthConfig) Validate() error {
	validator := validator.New()
	if err := validator.Struct(a); err != nil {
		return err
	}
	return nil
}
