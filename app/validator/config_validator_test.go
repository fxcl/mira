package validator

import (
	"testing"

	"mira/app/dto"
	"mira/common/xerrors"
)

func TestCreateConfigValidator(t *testing.T) {
	type args struct {
		param dto.CreateConfigRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_config_name",
			args: args{
				param: dto.CreateConfigRequest{
					ConfigName:  "",
					ConfigKey:   "key",
					ConfigValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrConfigNameEmpty,
		},
		{
			name: "empty_config_key",
			args: args{
				param: dto.CreateConfigRequest{
					ConfigName:  "name",
					ConfigKey:   "",
					ConfigValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrConfigKeyEmpty,
		},
		{
			name: "empty_config_value",
			args: args{
				param: dto.CreateConfigRequest{
					ConfigName:  "name",
					ConfigKey:   "key",
					ConfigValue: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrConfigValueEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateConfigRequest{
					ConfigName:  "name",
					ConfigKey:   "key",
					ConfigValue: "value",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateConfigValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreateConfigValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreateConfigValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdateConfigValidator(t *testing.T) {
	type args struct {
		param dto.UpdateConfigRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_config_id",
			args: args{
				param: dto.UpdateConfigRequest{
					ConfigId:    0,
					ConfigName:  "name",
					ConfigKey:   "key",
					ConfigValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_config_name",
			args: args{
				param: dto.UpdateConfigRequest{
					ConfigId:    1,
					ConfigName:  "",
					ConfigKey:   "key",
					ConfigValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrConfigNameEmpty,
		},
		{
			name: "empty_config_key",
			args: args{
				param: dto.UpdateConfigRequest{
					ConfigId:    1,
					ConfigName:  "name",
					ConfigKey:   "",
					ConfigValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrConfigKeyEmpty,
		},
		{
			name: "empty_config_value",
			args: args{
				param: dto.UpdateConfigRequest{
					ConfigId:    1,
					ConfigName:  "name",
					ConfigKey:   "key",
					ConfigValue: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrConfigValueEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateConfigRequest{
					ConfigId:    1,
					ConfigName:  "name",
					ConfigKey:   "key",
					ConfigValue: "value",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateConfigValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateConfigValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateConfigValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
