package validator

import (
	"errors"
	"mira/app/dto"
	"mira/common/types/regexp"
	"mira/common/utils"
)

// UpdateProfileValidator validates the request to update a user's profile.
func UpdateProfileValidator(param dto.UpdateProfileRequest) error {
	if param.NickName == "" {
		return errors.New("please enter the user nickname")
	}

	if !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return errors.New("incorrect email format")
	}

	if !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return errors.New("incorrect mobile number format")
	}

	return nil
}

// UserProfileUpdatePwdValidator validates the request to update a user's password.
func UserProfileUpdatePwdValidator(param dto.UserProfileUpdatePwdRequest) error {
	if param.OldPassword == "" {
		return errors.New("please enter the old password")
	}

	if param.NewPassword == "" {
		return errors.New("please enter the new password")
	}

	return nil
}

// CreateUserValidator validates the request to create a user.
func CreateUserValidator(param dto.CreateUserRequest) error {
	if param.NickName == "" {
		return errors.New("please enter the user nickname")
	}

	if param.UserName == "" {
		return errors.New("please enter the username")
	}

	if param.Password == "" {
		return errors.New("please enter the user password")
	}

	if param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return errors.New("incorrect mobile number format")
	}

	if param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return errors.New("incorrect email account format")
	}

	return nil
}

// UpdateUserValidator validates the request to update a user.
func UpdateUserValidator(param dto.UpdateUserRequest) error {
	if param.UserId <= 0 {
		return errors.New("parameter error")
	}

	if param.NickName == "" {
		return errors.New("please enter the user nickname")
	}

	if param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return errors.New("incorrect mobile number format")
	}

	if param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return errors.New("incorrect email account format")
	}

	return nil
}

// RemoveUserValidator validates the request to remove a user.
func RemoveUserValidator(userIds []int, authUserId int) error {
	if utils.Contains(userIds, 1) {
		return errors.New("the super administrator cannot be deleted")
	}

	if utils.Contains(userIds, authUserId) {
		return errors.New("the current user cannot be deleted")
	}

	return nil
}

// ChangeUserStatusValidator validates the request to change the user status.
func ChangeUserStatusValidator(param dto.UpdateUserRequest) error {
	if param.UserId <= 0 {
		return errors.New("parameter error")
	}

	if param.Status == "" {
		return errors.New("please select a status")
	}

	return nil
}

// ResetUserPwdValidator validates the request to reset a user's password.
func ResetUserPwdValidator(param dto.UpdateUserRequest) error {
	if param.UserId <= 0 {
		return errors.New("parameter error")
	}

	if param.Password == "" {
		return errors.New("please enter the user password")
	}

	return nil
}

// ImportUserValidator validates the request to import a user.
func ImportUserValidator(param dto.CreateUserRequest) error {
	if param.NickName == "" {
		return errors.New("please enter the user nickname")
	}

	if param.UserName == "" {
		return errors.New("please enter the username")
	}

	if param.Phonenumber != "" && !utils.CheckRegex(regexp.PHONE, param.Phonenumber) {
		return errors.New("incorrect mobile number format")
	}

	if param.Email != "" && !utils.CheckRegex(regexp.EMAIL, param.Email) {
		return errors.New("incorrect email account format")
	}

	return nil
}
