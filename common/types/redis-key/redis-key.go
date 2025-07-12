package rediskey

import "mira/config"

// CaptchaCodeKey returns the redis key for the captcha code.
func CaptchaCodeKey() string {
	return config.Data.Ruoyi.Name + ":captcha:code:"
}

// LoginPasswordErrorKey returns the redis key for the login password error count.
func LoginPasswordErrorKey() string {
	return config.Data.Ruoyi.Name + ":login:password:error:"
}

// UserTokenKey returns the redis key for the login user.
func UserTokenKey() string {
	return config.Data.Ruoyi.Name + ":user:token:"
}

// RepeatSubmitKey returns the redis key for anti-resubmission.
func RepeatSubmitKey() string {
	return config.Data.Ruoyi.Name + ":repeat:submit:"
}

// SysConfigKey returns the redis key for the system config data.
func SysConfigKey() string {
	return config.Data.Ruoyi.Name + ":system:config"
}

// SysDictKey returns the redis key for the system dictionary data.
func SysDictKey() string {
	return config.Data.Ruoyi.Name + ":system:dict:data"
}
