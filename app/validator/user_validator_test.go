package validator

import (
	"testing"

	"mira/app/dto"
	"mira/common/xerrors"
)

func TestUpdateProfileValidator(t *testing.T) {
	type args struct {
		param dto.UpdateProfileRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_nickname",
			args: args{
				param: dto.UpdateProfileRequest{
					NickName:    "",
					Email:       "test@example.com",
					Phonenumber: "13800138000",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserNicknameEmpty,
		},
		{
			name: "invalid_email",
			args: args{
				param: dto.UpdateProfileRequest{
					NickName:    "test",
					Email:       "invalid-email",
					Phonenumber: "13800138000",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserEmailFormat,
		},
		{
			name: "invalid_phone",
			args: args{
				param: dto.UpdateProfileRequest{
					NickName:    "test",
					Email:       "test@example.com",
					Phonenumber: "123",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserPhoneFormat,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateProfileRequest{
					NickName:    "test",
					Email:       "test@example.com",
					Phonenumber: "13800138000",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateProfileValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateProfileValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateProfileValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUserProfileUpdatePwdValidator(t *testing.T) {
	type args struct {
		param dto.UserProfileUpdatePwdRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_old_password",
			args: args{
				param: dto.UserProfileUpdatePwdRequest{
					OldPassword: "",
					NewPassword: "new_password",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserOldPasswordEmpty,
		},
		{
			name: "empty_new_password",
			args: args{
				param: dto.UserProfileUpdatePwdRequest{
					OldPassword: "old_password",
					NewPassword: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserNewPasswordEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UserProfileUpdatePwdRequest{
					OldPassword: "old_password",
					NewPassword: "new_password",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UserProfileUpdatePwdValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UserProfileUpdatePwdValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UserProfileUpdatePwdValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestCreateUserValidator(t *testing.T) {
	type args struct {
		param dto.CreateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_nickname",
			args: args{
				param: dto.CreateUserRequest{
					NickName: "",
					UserName: "test",
					Password: "password",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserNicknameEmpty,
		},
		{
			name: "empty_username",
			args: args{
				param: dto.CreateUserRequest{
					NickName: "test",
					UserName: "",
					Password: "password",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserNameEmpty,
		},
		{
			name: "empty_password",
			args: args{
				param: dto.CreateUserRequest{
					NickName: "test",
					UserName: "test",
					Password: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserPasswordEmpty,
		},
		{
			name: "invalid_phone",
			args: args{
				param: dto.CreateUserRequest{
					NickName:    "test",
					UserName:    "test",
					Password:    "password",
					Phonenumber: "123",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserPhoneFormat,
		},
		{
			name: "invalid_email",
			args: args{
				param: dto.CreateUserRequest{
					NickName: "test",
					UserName: "test",
					Password: "password",
					Email:    "invalid-email",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserEmailFormat,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateUserRequest{
					NickName:    "test",
					UserName:    "test",
					Password:    "password",
					Phonenumber: "13800138000",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateUserValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreateUserValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreateUserValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdateUserValidator(t *testing.T) {
	type args struct {
		param dto.UpdateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_user_id",
			args: args{
				param: dto.UpdateUserRequest{
					UserId: 0,
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_nickname",
			args: args{
				param: dto.UpdateUserRequest{
					UserId:   1,
					NickName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserNicknameEmpty,
		},
		{
			name: "invalid_phone",
			args: args{
				param: dto.UpdateUserRequest{
					UserId:      1,
					NickName:    "test",
					Phonenumber: "123",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserPhoneFormat,
		},
		{
			name: "invalid_email",
			args: args{
				param: dto.UpdateUserRequest{
					UserId:   1,
					NickName: "test",
					Email:    "invalid-email",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserEmailFormat,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateUserRequest{
					UserId:      1,
					NickName:    "test",
					Phonenumber: "13800138000",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUserValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateUserValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestRemoveUserValidator(t *testing.T) {
	type args struct {
		userIds    []int
		authUserId int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "delete_super_admin",
			args: args{
				userIds: []int{1, 2, 3},
			},
			wantErr: true,
			err:     xerrors.ErrUserSuperAdminDelete,
		},
		{
			name: "delete_current_user",
			args: args{
				userIds:    []int{2, 3},
				authUserId: 2,
			},
			wantErr: true,
			err:     xerrors.ErrUserCurrentUserDelete,
		},
		{
			name: "success",
			args: args{
				userIds:    []int{2, 3},
				authUserId: 4,
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveUserValidator(tt.args.userIds, tt.args.authUserId); (err != nil) != tt.wantErr {
				t.Errorf("RemoveUserValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("RemoveUserValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestChangeUserStatusValidator(t *testing.T) {
	type args struct {
		param dto.UpdateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_user_id",
			args: args{
				param: dto.UpdateUserRequest{
					UserId: 0,
					Status: "0",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_status",
			args: args{
				param: dto.UpdateUserRequest{
					UserId: 1,
					Status: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserStatusEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateUserRequest{
					UserId: 1,
					Status: "0",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeUserStatusValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("ChangeUserStatusValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("ChangeUserStatusValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestResetUserPwdValidator(t *testing.T) {
	type args struct {
		param dto.UpdateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_user_id",
			args: args{
				param: dto.UpdateUserRequest{
					UserId:   0,
					Password: "password",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_password",
			args: args{
				param: dto.UpdateUserRequest{
					UserId:   1,
					Password: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserPasswordEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateUserRequest{
					UserId:   1,
					Password: "password",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ResetUserPwdValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("ResetUserPwdValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("ResetUserPwdValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestImportUserValidator(t *testing.T) {
	type args struct {
		param dto.CreateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_nickname",
			args: args{
				param: dto.CreateUserRequest{
					NickName: "",
					UserName: "test",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserNicknameEmpty,
		},
		{
			name: "empty_username",
			args: args{
				param: dto.CreateUserRequest{
					NickName: "test",
					UserName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserNameEmpty,
		},
		{
			name: "invalid_phone",
			args: args{
				param: dto.CreateUserRequest{
					NickName:    "test",
					UserName:    "test",
					Phonenumber: "123",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserPhoneFormat,
		},
		{
			name: "invalid_email",
			args: args{
				param: dto.CreateUserRequest{
					NickName: "test",
					UserName: "test",
					Email:    "invalid-email",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUserEmailFormat,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateUserRequest{
					NickName:    "test",
					UserName:    "test",
					Phonenumber: "13800138000",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ImportUserValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("ImportUserValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("ImportUserValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
