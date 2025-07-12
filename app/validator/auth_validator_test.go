package validator

import (
	"testing"

	"mira/app/dto"
	"mira/common/xerrors"
)

func TestRegisterValidator(t *testing.T) {
	type args struct {
		param dto.RegisterRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_username",
			args: args{
				param: dto.RegisterRequest{
					Username:        "",
					Password:        "123456",
					ConfirmPassword: "123456",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUsernameEmpty,
		},
		{
			name: "empty_password",
			args: args{
				param: dto.RegisterRequest{
					Username:        "test",
					Password:        "",
					ConfirmPassword: "123456",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPasswordEmpty,
		},
		{
			name: "passwords_not_match",
			args: args{
				param: dto.RegisterRequest{
					Username:        "test",
					Password:        "123456",
					ConfirmPassword: "654321",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPasswordsNotMatch,
		},
		{
			name: "username_too_short",
			args: args{
				param: dto.RegisterRequest{
					Username:        "a",
					Password:        "123456",
					ConfirmPassword: "123456",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUsernameLength,
		},
		{
			name: "username_too_long",
			args: args{
				param: dto.RegisterRequest{
					Username:        "123456789012345678901",
					Password:        "123456",
					ConfirmPassword: "123456",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUsernameLength,
		},
		{
			name: "password_too_short",
			args: args{
				param: dto.RegisterRequest{
					Username:        "test",
					Password:        "1234",
					ConfirmPassword: "1234",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPasswordLength,
		},
		{
			name: "password_too_long",
			args: args{
				param: dto.RegisterRequest{
					Username:        "test",
					Password:        "123456789012345678901",
					ConfirmPassword: "123456789012345678901",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPasswordLength,
		},
		{
			name: "success",
			args: args{
				param: dto.RegisterRequest{
					Username:        "test",
					Password:        "123456",
					ConfirmPassword: "123456",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegisterValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("RegisterValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("RegisterValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestLoginValidator(t *testing.T) {
	type args struct {
		param dto.LoginRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_username",
			args: args{
				param: dto.LoginRequest{
					Username: "",
					Password: "123456",
				},
			},
			wantErr: true,
			err:     xerrors.ErrUsernameEmpty,
		},
		{
			name: "empty_password",
			args: args{
				param: dto.LoginRequest{
					Username: "test",
					Password: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPasswordEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.LoginRequest{
					Username: "test",
					Password: "123456",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoginValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("LoginValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("LoginValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
