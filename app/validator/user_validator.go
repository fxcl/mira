package validator

import (
	"mira/app/dto"
	"mira/common/types/regexp"
	"mira/common/utils"
	"mira/common/xerrors"
)

// UpdateProfileValidator validates the request to update a user's profile.
func UpdateProfileValidator(param dto.UpdateProfileRequest) error {
	if param.NickName == "" {
		return xerrors.ErrUserNicknameEmpty
	}

	if !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return xerrors.ErrUserEmailFormat
	}

	if !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return xerrors.ErrUserPhoneFormat
	}

	return nil
}

// UserProfileUpdatePwdValidator validates the request to update a user's password.
func UserProfileUpdatePwdValidator(param dto.UserProfileUpdatePwdRequest) error {
	if param.OldPassword == "" {
		return xerrors.ErrUserOldPasswordEmpty
	}

	if param.NewPassword == "" {
		return xerrors.ErrUserNewPasswordEmpty
	}

	return nil
}

// CreateUserValidator validates the request to create a user.
func CreateUserValidator(param dto.CreateUserRequest) error {
	if param.NickName == "" {
		return xerrors.ErrUserNicknameEmpty
	}

	if param.UserName == "" {
		return xerrors.ErrUserNameEmpty
	}

	if param.Password == "" {
		return xerrors.ErrUserPasswordEmpty
	}

	if param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return xerrors.ErrUserPhoneFormat
	}

	if param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return xerrors.ErrUserEmailFormat
	}

	return nil
}

// UpdateUserValidator validates the request to update a user.
func UpdateUserValidator(param dto.UpdateUserRequest) error {
	if param.UserId <= 0 {
		return xerrors.ErrParam
	}

	if param.NickName == "" {
		return xerrors.ErrUserNicknameEmpty
	}

	if param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return xerrors.ErrUserPhoneFormat
	}

	if param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return xerrors.ErrUserEmailFormat
	}

	return nil
}

// RemoveUserValidator validates the request to remove a user.
func RemoveUserValidator(userIds []int, authUserId int) error {
	if utils.Contains(userIds, 1) {
		return xerrors.ErrUserSuperAdminDelete
	}

	if utils.Contains(userIds, authUserId) {
		return xerrors.ErrUserCurrentUserDelete
	}

	return nil
}

// ChangeUserStatusValidator validates the request to change the user status.
func ChangeUserStatusValidator(param dto.UpdateUserRequest) error {
	if param.UserId <= 0 {
		return xerrors.ErrParam
	}

	if param.Status == "" {
		return xerrors.ErrUserStatusEmpty
	}

	return nil
}

// ResetUserPwdValidator validates the request to reset a user's password.
func ResetUserPwdValidator(param dto.UpdateUserRequest) error {
	if param.UserId <= 0 {
		return xerrors.ErrParam
	}

	if param.Password == "" {
		return xerrors.ErrUserPasswordEmpty
	}

	return nil
}

// ImportUserValidator validates the request to import a user.
func ImportUserValidator(param dto.CreateUserRequest) error {
	if param.NickName == "" {
		return xerrors.ErrUserNicknameEmpty
	}

	if param.UserName == "" {
		return xerrors.ErrUserNameEmpty
	}

	if param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return xerrors.ErrUserPhoneFormat
	}

	if param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return xerrors.ErrUserEmailFormat
	}

	return nil
}
