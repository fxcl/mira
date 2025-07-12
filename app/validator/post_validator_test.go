package validator

import (
	"testing"

	"mira/app/dto"
	"mira/common/xerrors"
)

func TestCreatePostValidator(t *testing.T) {
	type args struct {
		param dto.CreatePostRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "empty_post_code",
			args: args{
				param: dto.CreatePostRequest{
					PostCode: "",
					PostName: "name",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPostCodeEmpty,
		},
		{
			name: "empty_post_name",
			args: args{
				param: dto.CreatePostRequest{
					PostCode: "code",
					PostName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPostNameEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.CreatePostRequest{
					PostCode: "code",
					PostName: "name",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreatePostValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("CreatePostValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("CreatePostValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}

func TestUpdatePostValidator(t *testing.T) {
	type args struct {
		param dto.UpdatePostRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		err     error
	}{
		{
			name: "invalid_post_id",
			args: args{
				param: dto.UpdatePostRequest{
					PostId:   0,
					PostCode: "code",
					PostName: "name",
				},
			},
			wantErr: true,
			err:     xerrors.ErrParam,
		},
		{
			name: "empty_post_code",
			args: args{
				param: dto.UpdatePostRequest{
					PostId:   1,
					PostCode: "",
					PostName: "name",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPostCodeEmpty,
		},
		{
			name: "empty_post_name",
			args: args{
				param: dto.UpdatePostRequest{
					PostId:   1,
					PostCode: "code",
					PostName: "",
				},
			},
			wantErr: true,
			err:     xerrors.ErrPostNameEmpty,
		},
		{
			name: "success",
			args: args{
				param: dto.UpdatePostRequest{
					PostId:   1,
					PostCode: "code",
					PostName: "name",
				},
			},
			wantErr: false,
			err:     nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdatePostValidator(tt.args.param); (err != nil) != tt.wantErr {
				t.Errorf("UpdatePostValidator() error = %v, wantErr %v", err, tt.wantErr)
			} else if err != tt.err {
				t.Errorf("UpdatePostValidator() error = %v, want %v", err, tt.err)
			}
		})
	}
}
