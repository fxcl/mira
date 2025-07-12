package rediskey

import (
	"testing"

	"mira/config"
)

func TestRedisKeys(t *testing.T) {
	// Mock config for testing
	config.Data = &config.Config{}
	config.Data.Ruoyi.Name = "test-project"

	tests := []struct {
		name     string
		actual   string
		expected string
	}{
		{"CaptchaCodeKey", CaptchaCodeKey(), "test-project:captcha:code:"},
		{"LoginPasswordErrorKey", LoginPasswordErrorKey(), "test-project:login:password:error:"},
		{"UserTokenKey", UserTokenKey(), "test-project:user:token:"},
		{"RepeatSubmitKey", RepeatSubmitKey(), "test-project:repeat:submit:"},
		{"SysConfigKey", SysConfigKey(), "test-project:system:config"},
		{"SysDictKey", SysDictKey(), "test-project:system:dict:data"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.actual != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, tt.actual)
			}
		})
	}
}
