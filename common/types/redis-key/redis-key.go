package rediskey

import (
	"fmt"
	"mira/config"
)

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

// User-specific cache keys for performance optimization
func UserProfileKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:profile:" + fmt.Sprintf("%d", userID)
}

func UserPermsKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:permissions:" + fmt.Sprintf("%d", userID)
}

func UserRolesKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:roles:" + fmt.Sprintf("%d", userID)
}

func UserDataScopeKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:data_scope:" + fmt.Sprintf("%d", userID)
}

func UserSessionKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:session:" + fmt.Sprintf("%d", userID)
}

func UserAuthTokensKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:auth_tokens:" + fmt.Sprintf("%d", userID)
}

// System-wide cache keys
func RolePermsKey(roleID int) string {
	return config.Data.Ruoyi.Name + ":role:permissions:" + fmt.Sprintf("%d", roleID)
}

func MenuTreeKey() string {
	return config.Data.Ruoyi.Name + ":menu:tree"
}

func DeptTreeKey() string {
	return config.Data.Ruoyi.Name + ":dept:tree"
}

func SystemStatusKey() string {
	return config.Data.Ruoyi.Name + ":system:status"
}

func OnlineUsersKey() string {
	return config.Data.Ruoyi.Name + ":users:online"
}

// Data scope cache keys
func UserDataScopeDeptsKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:data_scope:depts:" + fmt.Sprintf("%d", userID)
}

func UserDataScopeUsersKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:data_scope:users:" + fmt.Sprintf("%d", userID)
}

// Permission cache keys
func UserAllPermsKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:permissions:all:" + fmt.Sprintf("%d", userID)
}

func UserMenuPermsKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:permissions:menu:" + fmt.Sprintf("%d", userID)
}

func UserBtnPermsKey(userID int) string {
	return config.Data.Ruoyi.Name + ":user:permissions:btn:" + fmt.Sprintf("%d", userID)
}

// Cache key patterns for invalidation
func UserPattern() string {
	return config.Data.Ruoyi.Name + ":user:*"
}

func RolePattern() string {
	return config.Data.Ruoyi.Name + ":role:*"
}

func SystemPattern() string {
	return config.Data.Ruoyi.Name + ":system:*"
}

// GetCacheKeyWithID generates a cache key with ID parameter
func GetCacheKeyWithID(baseKey string, id int) string {
	return fmt.Sprintf(baseKey, id)
}
