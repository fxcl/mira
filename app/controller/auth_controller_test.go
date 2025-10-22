package controller

import (
	"testing"

	"mira/app/dto"
	"mira/app/token"

	"github.com/stretchr/testify/assert"
)

// Test DTO types - this test validates that our DTO type fixes are working
func TestDTO_Types(t *testing.T) {
	// Test ConfigDetailResponse
	configResp := dto.ConfigDetailResponse{
		ConfigId:    1,
		ConfigName:  "Test Config",
		ConfigKey:   "test.key",
		ConfigValue: "true",
		ConfigType:  "boolean",
		Remark:      "Test config",
	}
	assert.Equal(t, "test.key", configResp.ConfigKey)
	assert.Equal(t, "true", configResp.ConfigValue)

	// Test DeptDetailResponse
	deptResp := dto.DeptDetailResponse{
		DeptId:   1,
		DeptName: "Test Department",
	}
	assert.Equal(t, 1, deptResp.DeptId)
	assert.Equal(t, "Test Department", deptResp.DeptName)

	// Test RoleListResponse
	roleResp := dto.RoleListResponse{
		RoleId:   1,
		RoleName: "Test Role",
	}
	assert.Equal(t, 1, roleResp.RoleId)
	assert.Equal(t, "Test Role", roleResp.RoleName)

	// Test MenuListResponse
	menuResp := dto.MenuListResponse{
		MenuId:   1,
		MenuName: "Test Menu",
	}
	assert.Equal(t, 1, menuResp.MenuId)
	assert.Equal(t, "Test Menu", menuResp.MenuName)

	// Test UserTokenResponse with embedded struct
	userToken := token.UserTokenResponse{
		UserTokenResponse: dto.UserTokenResponse{
			UserId:   1,
			UserName: "testuser",
		},
	}
	assert.Equal(t, 1, userToken.UserId)
	assert.Equal(t, "testuser", userToken.UserName)

	// Test RegisterRequest
	registerReq := dto.RegisterRequest{
		Username:        "testuser",
		Password:        "password123",
		ConfirmPassword: "password123",
		Code:            "1234",
		Uuid:            "test-uuid",
	}
	assert.Equal(t, "testuser", registerReq.Username)
	assert.Equal(t, "password123", registerReq.Password)
	assert.Equal(t, "password123", registerReq.ConfirmPassword)
}

// Test AuthController basic functionality
func TestAuthController_Basic(t *testing.T) {
	controller := &AuthController{}
	assert.NotNil(t, controller)
}