package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// RegisterValidator validates the registration request.
func RegisterValidator(param dto.RegisterRequest) error {
	if param.Username == "" {
		return xerrors.ErrUsernameEmpty
	}

	if param.Password == "" {
		return xerrors.ErrPasswordEmpty
	}

	if param.ConfirmPassword != param.Password {
		return xerrors.ErrPasswordsNotMatch
	}

	if len(param.Username) < 2 || len(param.Username) > 20 {
		return xerrors.ErrUsernameLength
	}

	if len(param.Password) < 5 || len(param.Password) > 20 {
		return xerrors.ErrPasswordLength
	}

	return nil
}

// LoginValidator validates the login request.
func LoginValidator(param dto.LoginRequest) error {
	if param.Username == "" {
		return xerrors.ErrUsernameEmpty
	}

	if param.Password == "" {
		return xerrors.ErrPasswordEmpty
	}

	return nil
}
