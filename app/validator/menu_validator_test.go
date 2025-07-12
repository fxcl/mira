package validator

import (
	"testing"

	"mira/app/dto"
	"mira/common/types/constant"
	"mira/common/xerrors"
)

func TestCreateMenuValidator(t *testing.T) {
	type args struct {
		param dto.CreateMenuRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_menu_name",
			args: args{
				param: dto.CreateMenuRequest{
					MenuName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuNameEmpty,
		},
		{
			name: "empty_path_for_directory",
			args: args{
				param: dto.CreateMenuRequest{
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_DIRECTORY,
					Path:     "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuPathEmpty,
		},
		{
			name: "empty_path_for_menu",
			args: args{
				param: dto.CreateMenuRequest{
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_MENU,
					Path:     "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuPathEmpty,
		},
		{
			name: "invalid_frame_path",
			args: args{
				param: dto.CreateMenuRequest{
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_MENU,
					Path:     "invalid_path",
					IsFrame:  constant.MENU_YES_FRAME,
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuPathHttpPrefix,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateMenuRequest{
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_MENU,
					Path:     "/path",
					IsFrame:  constant.MENU_NO_FRAME,
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateMenuValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreateMenuValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreateMenuValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdateMenuValidator(t *testing.T) {
	type args struct {
		param dto.UpdateMenuRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_menu_id",
			args: args{
				param: dto.UpdateMenuRequest{
					MenuId: 0,
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_menu_name",
			args: args{
				param: dto.UpdateMenuRequest{
					MenuId:   1,
					MenuName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuNameEmpty,
		},
		{
			name: "empty_path_for_directory",
			args: args{
				param: dto.UpdateMenuRequest{
					MenuId:   1,
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_DIRECTORY,
					Path:     "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuPathEmpty,
		},
		{
			name: "empty_path_for_menu",
			args: args{
				param: dto.UpdateMenuRequest{
					MenuId:   1,
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_MENU,
					Path:     "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuPathEmpty,
		},
		{
			name: "invalid_frame_path",
			args: args{
				param: dto.UpdateMenuRequest{
					MenuId:   1,
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_MENU,
					Path:     "invalid_path",
					IsFrame:  constant.MENU_YES_FRAME,
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuPathHttpPrefix,
		},
		{
			name: "parent_is_self",
			args: args{
				param: dto.UpdateMenuRequest{
					MenuId:   1,
					ParentId: 1,
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_MENU,
					Path:     "/path",
					IsFrame:  constant.MENU_NO_FRAME,
				},
			},
			wantErr: true,
			err:     xerrors.ErrMenuParentSelf,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateMenuRequest{
					MenuId:   1,
					MenuName: "menu",
					MenuType: constant.MENU_TYPE_MENU,
					Path:     "/path",
					IsFrame:  constant.MENU_NO_FRAME,
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateMenuValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateMenuValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateMenuValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
