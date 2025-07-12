package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"mira/anima/dal"
	"mira/anima/datetime"
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/token"
	"mira/app/validator"
	"mira/common/captcha"
	"mira/common/password"
	"mira/common/types/constant"
	"mira/common/xerrors"
	"mira/config"

	rediskey "mira/common/types/redis-key"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthController struct{}

// Get verification code
func (*AuthController) CaptchaImage(ctx *gin.Context) {
	captcha := captcha.NewCaptcha()

	id, b64s := captcha.Generate()

	b64s = strings.Replace(b64s, "data:image/png;base64,", "", 1)

	config := (&service.ConfigService{}).GetConfigCacheByConfigKey("sys.account.captchaEnabled")

	response.NewSuccess().SetData("uuid", id).SetData("img", b64s).SetData("captchaEnabled", config.ConfigValue == "true").Json(ctx)
}

// Register
func (*AuthController) Register(ctx *gin.Context) {
	if config := (&service.ConfigService{}).GetConfigCacheByConfigKey("sys.account.registerUser"); config.ConfigValue != "true" {
		response.NewError().SetMsg("The current system does not have the registration function enabled").Json(ctx)
		return
	}

	var param dto.RegisterRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.RegisterValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if config := (&service.ConfigService{}).GetConfigCacheByConfigKey("sys.account.captchaEnabled"); config.ConfigValue == "true" {
		if err := captcha.NewCaptcha().Verify(param.Uuid, param.Code); err != nil {
			response.NewError().SetMsg(err.Error()).Json(ctx)
			return
		}
	}

	if user := (&service.UserService{}).GetUserByUsername(param.Username); user.UserId > 0 {
		response.NewError().SetMsg("Failed to save user " + param.Username + ", registration account already exists").Json(ctx)
		return
	}

	hashedPassword, err := password.Generate(param.Password)
	if err != nil {
		response.NewError().SetCode(500).SetMsg("Failed to process password").Json(ctx)
		return
	}
	if err := (&service.UserService{}).CreateUser(dto.SaveUser{
		UserName: param.Username,
		NickName: param.Username,
		Password: hashedPassword,
		Status:   "0",
		Remark:   "Registered user",
		CreateBy: "Registered user",
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Login
func (*AuthController) Login(ctx *gin.Context) {
	var param dto.LoginRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetCode(400).SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.LoginValidator(param); err != nil {
		response.NewError().SetCode(400).SetMsg(err.Error()).Json(ctx)
		return
	}

	if config := (&service.ConfigService{}).GetConfigCacheByConfigKey("sys.account.captchaEnabled"); config.ConfigValue == "true" {
		if err := captcha.NewCaptcha().Verify(param.Uuid, param.Code); err != nil {
			response.NewError().SetMsg(err.Error()).Json(ctx)
			return
		}
	}

	user := (&service.UserService{}).GetUserByUsername(param.Username)
	if user.UserId <= 0 || user.Status != constant.NORMAL_STATUS {
		response.NewError().SetMsg("User does not exist or is disabled").Json(ctx)
		return
	}

	// If the number of login password errors exceeds the limit, the account will be locked for 10 minutes
	count, _ := dal.Redis.Get(ctx.Request.Context(), rediskey.LoginPasswordErrorKey()+param.Username).Int()
	if count >= config.Data.User.Password.MaxRetryCount {
		response.NewError().SetMsg("The number of password errors has exceeded the limit, please try again in " + strconv.Itoa(config.Data.User.Password.LockTime) + " minutes").Json(ctx)
		return
	}

	if err := password.Verify(user.Password, param.Password); err != nil {
		if err == xerrors.ErrMismatchedPassword {
			// The number of password errors is increased by 1, and the cache expiration time is set to the lock time
			dal.Redis.Set(ctx.Request.Context(), rediskey.LoginPasswordErrorKey()+param.Username, count+1, time.Minute*time.Duration(config.Data.User.Password.LockTime))
			response.NewError().SetMsg("Password error").Json(ctx)
			return
		}
		response.NewError().SetCode(500).SetMsg("Failed to verify password").Json(ctx)
		return
	}

	// Login successful, delete the number of errors
	dal.Redis.Del(ctx.Request.Context(), rediskey.LoginPasswordErrorKey()+param.Username)

	claims := token.GetClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(config.Data.Token.ExpireTime)))
	token, err := token.GenerateToken(claims, user)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	// Update login ip and time
	(&service.UserService{}).UpdateUser(dto.SaveUser{
		UserId:    user.UserId,
		LoginIP:   ctx.ClientIP(),
		LoginDate: datetime.Datetime{Time: time.Now()},
	}, nil, nil)

	response.NewSuccess().SetData("token", token).Json(ctx)
}

// Get authorization information
func (*AuthController) GetInfo(ctx *gin.Context) {
	user := (&service.UserService{}).GetUserByUserId(security.GetAuthUserId(ctx))

	user.Admin = user.UserId == 1

	dept := (&service.DeptService{}).GetDeptByDeptId(user.DeptId)

	roles, err := (&service.RoleService{}).GetRoleListByUserId(user.UserId)
	if err != nil {
		response.NewError().SetMsg(fmt.Sprintf("Failed to get user roles: %v", err)).Json(ctx)
		return
	}

	data := dto.AuthUserInfoResponse{
		UserDetailResponse: user,
		Dept:               dept,
		Roles:              roles,
	}

	roleKeys, err := (&service.RoleService{}).GetRoleKeysByUserId(user.UserId)
	if err != nil {
		response.NewError().SetMsg(fmt.Sprintf("Failed to get user permission identifiers: %v", err)).Json(ctx)
		return
	}

	perms := (&service.MenuService{}).GetPermsByUserId(user.UserId)

	response.NewSuccess().SetData("user", data).SetData("roles", roleKeys).SetData("permissions", perms).Json(ctx)
}

// Get authorized routes
func (*AuthController) GetRouters(ctx *gin.Context) {
	menus := (&service.MenuService{}).GetMenuMCListByUserId(security.GetAuthUserId(ctx))

	tree := (&service.MenuService{}).MenusToTree(menus, 0)

	routers := (&service.MenuService{}).BuildRouterMenus(tree)

	response.NewSuccess().SetData("data", routers).Json(ctx)
}

// Logout
func (*AuthController) Logout(ctx *gin.Context) {
	tokenKey, err := token.GetUserTokenKey(ctx)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}
	token.DeleteToken(ctx.Request.Context(), tokenKey)

	response.NewSuccess().Json(ctx)
}
