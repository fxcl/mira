package validator

import (
	"errors"
	"testing"

	"mira/app/dto"
	"mira/common/xerrors"
)

func TestCreateRoleValidator(t *testing.T) {
	type args struct {
		param dto.CreateRoleRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_role_name",
			args: args{
				param: dto.CreateRoleRequest{
					RoleName: "",
					RoleKey:  "key",
				},
			},
			wantErr: true,
			err:     xerrors.ErrRoleNameEmpty,
		},
		{
			name: "empty_role_key",
			args: args{
				param: dto.CreateRoleRequest{
					RoleName: "name",
					RoleKey:  "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrRoleKeyEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateRoleRequest{
					RoleName: "name",
					RoleKey:  "key",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateRoleValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreateRoleValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreateRoleValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdateRoleValidator(t *testing.T) {
	type args struct {
		param dto.UpdateRoleRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_role_id",
			args: args{
				param: dto.UpdateRoleRequest{
					RoleId:   0,
					RoleName: "name",
					RoleKey:  "key",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_role_name",
			args: args{
				param: dto.UpdateRoleRequest{
					RoleId:   1,
					RoleName: "",
					RoleKey:  "key",
				},
			},
			wantErr: true,
			err:     xerrors.ErrRoleNameEmpty,
		},
		{
			name: "empty_role_key",
			args: args{
				param: dto.UpdateRoleRequest{
					RoleId:   1,
					RoleName: "name",
					RoleKey:  "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrRoleKeyEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateRoleRequest{
					RoleId:   1,
					RoleName: "name",
					RoleKey:  "key",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateRoleValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateRoleValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateRoleValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestRemoveRoleValidator(t *testing.T) {
	type args struct {
		roleIds  []int
		roleId   int
		roleName string
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
				roleIds: []int{1, 2, 3},
			},
			wantErr: true,
			err:     xerrors.ErrRoleSuperAdminDelete,
		},
		{
			name: "delete_current_role",
			args: args{
				roleIds:  []int{2, 3},
				roleId:   2,
				roleName: "test",
			},
			wantErr: true,
			err:     errors.New("the test role cannot be deleted"),
		},
		{
			name: "success",
			args: args{
				roleIds:  []int{2, 3},
				roleId:   4,
				roleName: "test",
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RemoveRoleValidator(tt.args.roleIds, tt.args.roleId, tt.args.roleName); (err != nil) != tt.wantErr {
				t.Errorf("RemoveRoleValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != nil && err.Error() != tt.err.Error() {
				t.Errorf("RemoveRoleValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestChangeRoleStatusValidator(t *testing.T) {
	type args struct {
		param dto.UpdateRoleRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_role_id",
			args: args{
				param: dto.UpdateRoleRequest{
					RoleId: 0,
					Status: "0",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_status",
			args: args{
				param: dto.UpdateRoleRequest{
					RoleId: 1,
					Status: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrRoleStatusEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateRoleRequest{
					RoleId: 1,
					Status: "0",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ChangeRoleStatusValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("ChangeRoleStatusValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("ChangeRoleStatusValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
