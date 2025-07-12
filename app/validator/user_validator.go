package validator

import (
	"mira/app/dto"
	"mira/common/types/regexp"
	"mira/common/utils"
	"mira/common/xerrors"
)

// UpdateProfileValidator validates the request to update a user's profile.
func UpdateProfileValidator(param dto.UpdateProfileRequest) error {
	switch {
	case param.NickName == "":
		return xerrors.ErrUserNicknameEmpty
	case param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email):
		return xerrors.ErrUserEmailFormat
	case param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber):
		return xerrors.ErrUserPhoneFormat
	default:
		return nil
	}
}

// UserProfileUpdatePwdValidator validates the request to update a user's password.
func UserProfileUpdatePwdValidator(param dto.UserProfileUpdatePwdRequest) error {
	switch {
	case param.OldPassword == "":
		return xerrors.ErrUserOldPasswordEmpty
	case param.NewPassword == "":
		return xerrors.ErrUserNewPasswordEmpty
	default:
		return nil
	}
}

// CreateUserValidator validates the request to create a user.
func CreateUserValidator(param dto.CreateUserRequest) error {
	switch {
	case param.NickName == "":
		return xerrors.ErrUserNicknameEmpty
	case param.UserName == "":
		return xerrors.ErrUserNameEmpty
	case param.Password == "":
		return xerrors.ErrUserPasswordEmpty
	case param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber):
		return xerrors.ErrUserPhoneFormat
	case param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email):
		return xerrors.ErrUserEmailFormat
	default:
		return nil
	}
}

// UpdateUserValidator validates the request to update a user.
func UpdateUserValidator(param dto.UpdateUserRequest) error {
	switch {
	case param.UserId <= 0:
		return xerrors.ErrParam
	case param.NickName == "":
		return xerrors.ErrUserNicknameEmpty
	case param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber):
		return xerrors.ErrUserPhoneFormat
	case param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email):
		return xerrors.ErrUserEmailFormat
	default:
		return nil
	}
}

// RemoveUserValidator validates the request to remove a user.
func RemoveUserValidator(userIds []int, authUserId int) error {
	switch {
	case utils.Contains(userIds, 1):
		return xerrors.ErrUserSuperAdminDelete
	case utils.Contains(userIds, authUserId):
		return xerrors.ErrUserCurrentUserDelete
	default:
		return nil
	}
}

// ChangeUserStatusValidator validates the request to change the user status.
func ChangeUserStatusValidator(param dto.UpdateUserRequest) error {
	switch {
	case param.UserId <= 0:
		return xerrors.ErrParam
	case param.Status == "":
		return xerrors.ErrUserStatusEmpty
	default:
		return nil
	}
}

// ResetUserPwdValidator validates the request to reset a user's password.
func ResetUserPwdValidator(param dto.UpdateUserRequest) error {
	switch {
	case param.UserId <= 0:
		return xerrors.ErrParam
	case param.Password == "":
		return xerrors.ErrUserPasswordEmpty
	default:
		return nil
	}
}

// ImportUserValidator validates the request to import a user.
func ImportUserValidator(param dto.CreateUserRequest) error {
	switch {
	case param.NickName == "":
		return xerrors.ErrUserNicknameEmpty
	case param.UserName == "":
		return xerrors.ErrUserNameEmpty
	case param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber):
		return xerrors.ErrUserPhoneFormat
	case param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email):
		return xerrors.ErrUserEmailFormat
	default:
		return nil
	}
}
