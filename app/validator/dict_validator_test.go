package validator

import (
	"testing"

	"mira/app/dto"
	"mira/common/xerrors"
)

func TestCreateDictTypeValidator(t *testing.T) {
	type args struct {
		param dto.CreateDictTypeRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_dict_name",
			args: args{
				param: dto.CreateDictTypeRequest{
					DictName: "",
					DictType: "type",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictNameEmpty,
		},
		{
			name: "empty_dict_type",
			args: args{
				param: dto.CreateDictTypeRequest{
					DictName: "name",
					DictType: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictTypeEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateDictTypeRequest{
					DictName: "name",
					DictType: "type",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDictTypeValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreateDictTypeValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreateDictTypeValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdateDictTypeValidator(t *testing.T) {
	type args struct {
		param dto.UpdateDictTypeRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_dict_id",
			args: args{
				param: dto.UpdateDictTypeRequest{
					DictId:   0,
					DictName: "name",
					DictType: "type",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_dict_name",
			args: args{
				param: dto.UpdateDictTypeRequest{
					DictId:   1,
					DictName: "",
					DictType: "type",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictNameEmpty,
		},
		{
			name: "empty_dict_type",
			args: args{
				param: dto.UpdateDictTypeRequest{
					DictId:   1,
					DictName: "name",
					DictType: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictTypeEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateDictTypeRequest{
					DictId:   1,
					DictName: "name",
					DictType: "type",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateDictTypeValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDictTypeValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateDictTypeValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestCreateDictDataValidator(t *testing.T) {
	type args struct {
		param dto.CreateDictDataRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_dict_label",
			args: args{
				param: dto.CreateDictDataRequest{
					DictLabel: "",
					DictValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictLabelEmpty,
		},
		{
			name: "empty_dict_value",
			args: args{
				param: dto.CreateDictDataRequest{
					DictLabel: "label",
					DictValue: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictValueEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateDictDataRequest{
					DictLabel: "label",
					DictValue: "value",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDictDataValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreateDictDataValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreateDictDataValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdateDictDataValidator(t *testing.T) {
	type args struct {
		param dto.UpdateDictDataRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_dict_code",
			args: args{
				param: dto.UpdateDictDataRequest{
					DictCode:  0,
					DictLabel: "label",
					DictValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_dict_label",
			args: args{
				param: dto.UpdateDictDataRequest{
					DictCode:  1,
					DictLabel: "",
					DictValue: "value",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictLabelEmpty,
		},
		{
			name: "empty_dict_value",
			args: args{
				param: dto.UpdateDictDataRequest{
					DictCode:  1,
					DictLabel: "label",
					DictValue: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDictValueEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateDictDataRequest{
					DictCode:  1,
					DictLabel: "label",
					DictValue: "value",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateDictDataValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDictDataValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateDictDataValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
