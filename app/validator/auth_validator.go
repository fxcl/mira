package validator

import (
	"errors"
	"mira/app/dto"
)

// RegisterValidator validates the registration request.
func RegisterValidator(param dto.RegisterRequest) error {
	if param.Username == "" {
		return errors.New("username cannot be empty")
	}

	if param.Password == "" {
		return errors.New("password cannot be empty")
	}

	if param.ConfirmPassword != param.Password {
		return errors.New("passwords do not match")
	}

	if len(param.Username) < 2 || len(param.Username) > 20 {
		return errors.New("username length must be between 2 and 20 characters")
	}

	if len(param.Password) < 5 || len(param.Password) > 20 {
		return errors.New("password length must be between 5 and 20 characters")
	}

	return nil
}

// LoginValidator validates the login request.
func LoginValidator(param dto.LoginRequest) error {
	if param.Username == "" {
		return errors.New("username cannot be empty")
	}

	if param.Password == "" {
		return errors.New("password cannot be empty")
	}

	return nil
}
