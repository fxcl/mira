package rediskey

import "mira/config"

var (
	// Captcha code redis key
	CaptchaCodeKey = config.Data.Ruoyi.Name + ":captcha:code:"

	// Login account password error count redis key
	LoginPasswordErrorKey = config.Data.Ruoyi.Name + ":login:password:error:"

	// Login user redis key
	UserTokenKey = config.Data.Ruoyi.Name + ":user:token:"

	// Anti-resubmission redis key
	RepeatSubmitKey = config.Data.Ruoyi.Name + ":repeat:submit:"

	// System config data redis key
	SysConfigKey = config.Data.Ruoyi.Name + ":system:config"

	// System dictionary data redis key
	SysDictKey = config.Data.Ruoyi.Name + ":system:dict:data"
)
