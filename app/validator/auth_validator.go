package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// RegisterValidator validates the registration request.
func RegisterValidator(param dto.RegisterRequest) error {
	switch {
	case param.Username == "":
		return xerrors.ErrUsernameEmpty
	case param.Password == "":
		return xerrors.ErrPasswordEmpty
	case param.ConfirmPassword != param.Password:
		return xerrors.ErrPasswordsNotMatch
	case len(param.Username) < 2 || len(param.Username) > 20:
		return xerrors.ErrUsernameLength
	case len(param.Password) < 5 || len(param.Password) > 20:
		return xerrors.ErrPasswordLength
	default:
		return nil
	}
}

// LoginValidator validates the login request.
func LoginValidator(param dto.LoginRequest) error {
	switch {
	case param.Username == "":
		return xerrors.ErrUsernameEmpty
	case param.Password == "":
		return xerrors.ErrPasswordEmpty
	default:
		return nil
	}
}
