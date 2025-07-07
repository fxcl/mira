package xerrors

import "errors"

// Validation Errors
var (
	// Upload
	ErrUnsupportedFileType        = errors.New("unsupported file type")
	ErrUploadDomainNotFound       = errors.New("domain not found, cannot generate access address")
	ErrUploadFileIncomplete       = errors.New("upload file data is incomplete and cannot be saved")
	ErrUploadFileMissingSuffix    = errors.New("file missing suffix")
	ErrUploadFileSizeExceedsLimit = errors.New("file size exceeds the limit")
	ErrUploadInvalidFileFormat    = errors.New("invalid file format")

	// Captcha
	ErrMismatchedPassword = errors.New("mismatched password")
	ErrCaptcha            = errors.New("captcha error")

	// Common
	ErrParam = errors.New("parameter error")

	// Auth
	ErrUsernameEmpty     = errors.New("username cannot be empty")
	ErrPasswordEmpty     = errors.New("password cannot be empty")
	ErrPasswordsNotMatch = errors.New("passwords do not match")
	ErrUsernameLength    = errors.New("username length must be between 2 and 20 characters")
	ErrPasswordLength    = errors.New("password length must be between 5 and 20 characters")

	// Config
	ErrConfigNameEmpty  = errors.New("please enter the parameter name")
	ErrConfigKeyEmpty   = errors.New("please enter the parameter key")
	ErrConfigValueEmpty = errors.New("please enter the parameter value")

	// Dept
	ErrParentDeptEmpty = errors.New("please select the parent department")
	ErrDeptNameEmpty   = errors.New("please enter the department name")
	ErrDeptParentSelf  = errors.New("the parent department cannot be itself")

	// Dict
	ErrDictNameEmpty  = errors.New("please enter the dictionary name")
	ErrDictTypeEmpty  = errors.New("please enter the dictionary type")
	ErrDictLabelEmpty = errors.New("please enter the data label")
	ErrDictValueEmpty = errors.New("please enter the data key value")

	// Menu
	ErrMenuNameEmpty      = errors.New("please enter the menu name")
	ErrMenuPathEmpty      = errors.New("please enter the route address")
	ErrMenuPathHttpPrefix = errors.New("the address must start with http(s)://")
	ErrMenuParentSelf     = errors.New("the parent menu cannot be itself")

	// Post
	ErrPostCodeEmpty = errors.New("please enter the post code")
	ErrPostNameEmpty = errors.New("please enter the post name")

	// Role
	ErrRoleNameEmpty        = errors.New("please enter the role name")
	ErrRoleKeyEmpty         = errors.New("please enter the permission string")
	ErrRoleSuperAdminDelete = errors.New("the super administrator cannot be deleted")
	ErrRoleInUseDelete      = errors.New("the role is in use and cannot be deleted")
	ErrRoleStatusEmpty      = errors.New("please select a status")

	// User
	ErrUserNicknameEmpty     = errors.New("please enter the user nickname")
	ErrUserEmailFormat       = errors.New("incorrect email format")
	ErrUserPhoneFormat       = errors.New("incorrect mobile number format")
	ErrUserOldPasswordEmpty  = errors.New("please enter the old password")
	ErrUserNewPasswordEmpty  = errors.New("please enter the new password")
	ErrUserNameEmpty         = errors.New("please enter the username")
	ErrUserPasswordEmpty     = errors.New("please enter the user password")
	ErrUserSuperAdminDelete  = errors.New("the super administrator cannot be deleted")
	ErrUserCurrentUserDelete = errors.New("the current user cannot be deleted")
	ErrUserStatusEmpty       = errors.New("please select a status")

	// General
	ErrNotImplemented = errors.New("not implemented")
	ErrInternal       = errors.New("internal server error")
)
