package validator

import (
	"testing"

	"mira/app/dto"
	"mira/common/xerrors"
)

func TestCreateDeptValidator(t *testing.T) {
	type args struct {
		param dto.CreateDeptRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_parent_id",
			args: args{
				param: dto.CreateDeptRequest{
					ParentId: 0,
					DeptName: "dept",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParentDeptEmpty,
		},
		{
			name: "empty_dept_name",
			args: args{
				param: dto.CreateDeptRequest{
					ParentId: 1,
					DeptName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDeptNameEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.CreateDeptRequest{
					ParentId: 1,
					DeptName: "dept",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateDeptValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreateDeptValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreateDeptValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdateDeptValidator(t *testing.T) {
	type args struct {
		param dto.UpdateDeptRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_dept_id",
			args: args{
				param: dto.UpdateDeptRequest{
					DeptId:   0,
					ParentId: 1,
					DeptName: "dept",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "invalid_parent_id",
			args: args{
				param: dto.UpdateDeptRequest{
					DeptId:   2,
					ParentId: 0,
					DeptName: "dept",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParentDeptEmpty,
		},
		{
			name: "empty_dept_name",
			args: args{
				param: dto.UpdateDeptRequest{
					DeptId:   2,
					ParentId: 1,
					DeptName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDeptNameEmpty,
		},
		{
			name: "parent_is_self",
			args: args{
				param: dto.UpdateDeptRequest{
					DeptId:   2,
					ParentId: 2,
					DeptName: "dept",
				},
			},
			wantErr: true,
			err:     xerrors.ErrDeptParentSelf,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdateDeptRequest{
					DeptId:   2,
					ParentId: 1,
					DeptName: "dept",
				},
			},
			wantErr: false,
			err:     nil,
		},
		{
			name: "success_no_parent",
			args: args{
				param: dto.UpdateDeptRequest{
					DeptId:   100,
					ParentId: 0,
					DeptName: "dept",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateDeptValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdateDeptValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdateDeptValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
