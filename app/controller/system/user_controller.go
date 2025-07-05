package systemcontroller

import (
	"io"
	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/password"
	"mira/common/upload"
	"mira/common/utils"
	"mira/config"
	"os"
	"strconv"
	"strings"
	"time"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
	excelize "github.com/xuri/excelize/v2"
)

type UserController struct{}

// Get department tree
func (*UserController) DeptTree(ctx *gin.Context) {
	depts := (&service.DeptService{}).GetUserDeptTree(security.GetAuthUserId(ctx))

	tree := (&service.UserService{}).DeptListToTree(depts, 0)

	response.NewSuccess().SetData("data", tree).Json(ctx)
}

// Get user list
func (*UserController) List(ctx *gin.Context) {
	var param dto.UserListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	users, total := (&service.UserService{}).GetUserList(param, security.GetAuthUserId(ctx), true)

	for key, user := range users {
		users[key].Dept.DeptName = user.DeptName
		users[key].Dept.Leader = user.Leader
	}

	response.NewSuccess().SetPageData(users, total).Json(ctx)
}

// User details
func (*UserController) Detail(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))

	response := response.NewSuccess()

	if userId > 0 {
		user := (&service.UserService{}).GetUserByUserId(userId)

		user.Admin = user.UserId == 1

		dept := (&service.DeptService{}).GetDeptByDeptId(user.DeptId)

		roles := (&service.RoleService{}).GetRoleListByUserId(user.UserId)

		response.SetData("data", dto.AuthUserInfoResponse{
			UserDetailResponse: user,
			Dept:               dept,
			Roles:              roles,
		})

		roleIds := make([]int, 0)
		for _, role := range roles {
			roleIds = append(roleIds, role.RoleId)
		}
		response.SetData("roleIds", roleIds)

		postIds := (&service.PostService{}).GetPostIdsByUserId(user.UserId)
		response.SetData("postIds", postIds)
	}

	roles, _ := (&service.RoleService{}).GetRoleList(dto.RoleListRequest{}, false)
	if userId != 1 {
		roles = utils.Filter(roles, func(role dto.RoleListResponse) bool {
			return role.RoleId != 1
		})
	}
	response.SetData("roles", roles)

	posts, _ := (&service.PostService{}).GetPostList(dto.PostListRequest{}, false)
	response.SetData("posts", posts)

	response.Json(ctx)
}

// Add user
func (*UserController) Create(ctx *gin.Context) {
	var param dto.CreateUserRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateUserValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if user := (&service.UserService{}).GetUserByUsername(param.UserName); user.UserId > 0 {
		response.NewError().SetMsg("Failed to add user " + param.UserName + ", username already exists").Json(ctx)
		return
	}

	if param.Email != "" {
		if user := (&service.UserService{}).GetUserByEmail(param.Email); user.UserId > 0 {
			response.NewError().SetMsg("Failed to add user " + param.UserName + ", email already exists").Json(ctx)
			return
		}
	}

	if param.Phonenumber != "" {
		if user := (&service.UserService{}).GetUserByPhonenumber(param.Phonenumber); user.UserId > 0 {
			response.NewError().SetMsg("Failed to add user " + param.UserName + ", phone number already exists").Json(ctx)
			return
		}
	}

	if err := (&service.UserService{}).CreateUser(dto.SaveUser{
		DeptId:      param.DeptId,
		UserName:    param.UserName,
		NickName:    param.NickName,
		Email:       param.Email,
		Phonenumber: param.Phonenumber,
		Sex:         param.Sex,
		Password:    password.Generate(param.Password),
		Status:      param.Status,
		Remark:      param.Remark,
		CreateBy:    security.GetAuthUserName(ctx),
	}, param.RoleIds, param.PostIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update user
func (*UserController) Update(ctx *gin.Context) {
	var param dto.UpdateUserRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateUserValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if param.Email != "" {
		if user := (&service.UserService{}).GetUserByEmail(param.Email); user.UserId > 0 && user.UserId != param.UserId {
			response.NewError().SetMsg("Failed to modify user " + param.UserName + ", email already exists").Json(ctx)
			return
		}
	}

	if param.Phonenumber != "" {
		if user := (&service.UserService{}).GetUserByPhonenumber(param.Phonenumber); user.UserId > 0 && user.UserId != param.UserId {
			response.NewError().SetMsg("Failed to modify user " + param.UserName + ", phone number already exists").Json(ctx)
			return
		}
	}

	if err := (&service.UserService{}).UpdateUser(dto.SaveUser{
		UserId:      param.UserId,
		DeptId:      param.DeptId,
		NickName:    param.NickName,
		Email:       param.Email,
		Phonenumber: param.Phonenumber,
		Sex:         param.Sex,
		Status:      param.Status,
		Remark:      param.Remark,
		UpdateBy:    security.GetAuthUserName(ctx),
	}, param.RoleIds, param.PostIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Delete user
func (*UserController) Remove(ctx *gin.Context) {
	userIds, err := utils.StringToIntSlice(ctx.Param("userIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = validator.RemoveUserValidator(userIds, security.GetAuthUserId(ctx)); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = (&service.UserService{}).DeleteUser(userIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Change user status
func (*UserController) ChangeStatus(ctx *gin.Context) {
	var param dto.UpdateUserRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.ChangeUserStatusValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.UserService{}).UpdateUser(dto.SaveUser{
		UserId:   param.UserId,
		Status:   param.Status,
		UpdateBy: security.GetAuthUserName(ctx),
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Reset user password
func (*UserController) ResetPwd(ctx *gin.Context) {
	var param dto.UpdateUserRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.ResetUserPwdValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.UserService{}).UpdateUser(dto.SaveUser{
		UserId:   param.UserId,
		Password: password.Generate(param.Password),
		UpdateBy: security.GetAuthUserName(ctx),
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Get authorized roles by user ID
func (*UserController) AuthRole(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))

	response := response.NewSuccess()

	var userHasRoleIds []int

	if userId > 0 {
		user := (&service.UserService{}).GetUserByUserId(userId)

		user.Admin = user.UserId == 1

		dept := (&service.DeptService{}).GetDeptByDeptId(user.DeptId)

		roles := (&service.RoleService{}).GetRoleListByUserId(user.UserId)
		for _, role := range roles {
			userHasRoleIds = append(userHasRoleIds, role.RoleId)
		}

		response.SetData("user", dto.AuthUserInfoResponse{
			UserDetailResponse: user,
			Dept:               dept,
			Roles:              roles,
		})
	}

	roles, _ := (&service.RoleService{}).GetRoleList(dto.RoleListRequest{}, false)
	if userId != 1 {
		roles = utils.Filter(roles, func(role dto.RoleListResponse) bool {
			return role.RoleId != 1
		})
		// Set the role selection flag, if the role is in the user's role list, set the flag to true
		for key, role := range roles {
			if utils.Contains(userHasRoleIds, role.RoleId) {
				roles[key].Flag = true
			}
		}
	}
	response.SetData("roles", roles)

	response.Json(ctx)
}

// User authorized role
func (*UserController) AddAuthRole(ctx *gin.Context) {
	var param dto.AddUserAuthRoleRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	roleIds, err := utils.StringToIntSlice(param.RoleIds, ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.UserService{}).AddAuthRole(param.UserId, roleIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Import user template
func (*UserController) ImportTemplate(ctx *gin.Context) {
	list := make([]dto.UserImportRequest, 0)

	list = append(list, dto.UserImportRequest{
		DeptId:      1,
		UserName:    "example",
		NickName:    "template",
		Email:       "example@example.com",
		Phonenumber: "12345678901",
		Sex:         "1",
		Status:      "0",
	})

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("user_template_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}

// Import user data
func (*UserController) ImportData(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	fileName := config.Data.Ruoyi.UploadPath + file.Filename

	// Temporarily save the file
	ctx.SaveUploadedFile(file, fileName)
	defer os.Remove(fileName)

	// Whether to update existing user data
	updateSupport, _ := strconv.ParseBool(ctx.Query("updateSupport"))

	excelFile, err := excelize.OpenFile(fileName)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.UserImportRequest, 0)

	if err = excel.ImportExcel(excelFile, &list, 0, 1); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if len(list) <= 0 {
		response.NewError().SetMsg("Import user data cannot be empty").Json(ctx)
		return
	}

	var successNum, failNum int
	var failMsg []string

	authUserName := security.GetAuthUserName(ctx)

	for _, item := range list {
		user := (&service.UserService{}).GetUserByUsername(item.UserName)

		// Insert new user
		if user.UserId <= 0 {
			if err = validator.ImportUserValidator(dto.CreateUserRequest{
				DeptId:      item.DeptId,
				UserName:    item.UserName,
				NickName:    item.NickName,
				Email:       item.Email,
				Phonenumber: item.Phonenumber,
				Sex:         item.Sex,
				Status:      item.Status,
			}); err != nil {
				failNum = failNum + 1
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be added: "+err.Error())
				continue
			}
			if err = (&service.UserService{}).CreateUser(dto.SaveUser{
				DeptId:      item.DeptId,
				UserName:    item.UserName,
				NickName:    item.NickName,
				Email:       item.Email,
				Phonenumber: item.Phonenumber,
				Sex:         item.Sex,
				Password:    password.Generate((&service.ConfigService{}).GetConfigCacheByConfigKey("sys.user.initPassword").ConfigValue),
				Status:      item.Status,
				CreateBy:    authUserName,
			}, nil, nil); err != nil {
				failNum = failNum + 1
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be added: "+err.Error())
				continue
			}
			successNum = successNum + 1
			continue
		} else if updateSupport {
			if err = validator.UpdateUserValidator(dto.UpdateUserRequest{
				UserId:      user.UserId,
				DeptId:      item.DeptId,
				NickName:    item.NickName,
				Email:       item.Email,
				Phonenumber: item.Phonenumber,
				Sex:         item.Sex,
				Status:      item.Status,
			}); err != nil {
				failNum = failNum + 1
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be updated: "+err.Error())
				continue
			}
			// Update existing user
			if err = (&service.UserService{}).UpdateUser(dto.SaveUser{
				UserId:      user.UserId,
				DeptId:      item.DeptId,
				NickName:    item.NickName,
				Email:       item.Email,
				Phonenumber: item.Phonenumber,
				Sex:         item.Sex,
				Status:      item.Status,
				UpdateBy:    authUserName,
			}, nil, nil); err != nil {
				failNum = failNum + 1
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be updated: "+err.Error())
				continue
			}
			successNum = successNum + 1
			// successMsg = append(successMsg, strconv.Itoa(successNum)+", Account "+item.UserName+" updated successfully")
			continue
		} else {
			failNum = failNum + 1
			failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" already exists")
		}
	}

	if failNum > 0 {
		response.NewError().SetMsg("Import failed, " + strconv.Itoa(failNum) + " pieces of data are wrong, the errors are as follows:" + strings.Join(failMsg, "<br/>")).Json(ctx)
		return
	}

	response.NewSuccess().SetMsg("Import successful, " + strconv.Itoa(successNum) + " pieces of data in total").Json(ctx)
}

// Export user data
func (*UserController) Export(ctx *gin.Context) {
	var param dto.UserListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.UserExportResponse, 0)

	users, _ := (&service.UserService{}).GetUserList(param, security.GetAuthUserId(ctx), false)
	for _, user := range users {

		loginDate := user.LoginDate.Format("2006-01-02 15:04:05")
		if user.LoginDate.IsZero() {
			loginDate = ""
		}

		list = append(list, dto.UserExportResponse{
			UserId:      user.UserId,
			UserName:    user.UserName,
			NickName:    user.NickName,
			Email:       user.Email,
			Phonenumber: user.Phonenumber,
			Sex:         user.Sex,
			Status:      user.Status,
			LoginIp:     user.LoginIp,
			LoginDate:   loginDate,
			DeptName:    user.DeptName,
			DeptLeader:  user.Leader,
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("user_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}

// Get personal information
func (*UserController) GetProfile(ctx *gin.Context) {
	user := (&service.UserService{}).GetUserByUserId(security.GetAuthUserId(ctx))

	user.Admin = user.UserId == 1

	dept := (&service.DeptService{}).GetDeptByDeptId(user.DeptId)

	roles := (&service.RoleService{}).GetRoleListByUserId(user.UserId)

	data := dto.AuthUserInfoResponse{
		UserDetailResponse: user,
		Dept:               dept,
		Roles:              roles,
	}

	// Get role group
	roleGroup := (&service.RoleService{}).GetRoleNamesByUserId(user.UserId)

	// Get post group
	postGroup := (&service.PostService{}).GetPostNamesByUserId(user.UserId)

	response.NewSuccess().
		SetData("data", data).
		SetData("roleGroup", strings.Join(roleGroup, ",")).
		SetData("postGroup", strings.Join(postGroup, ",")).
		Json(ctx)
}

// Update personal information
func (*UserController) UpdateProfile(ctx *gin.Context) {
	var param dto.UpdateProfileRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateProfileValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := (&service.UserService{}).UpdateUser(dto.SaveUser{
		UserId:      security.GetAuthUserId(ctx),
		NickName:    param.NickName,
		Email:       param.Email,
		Phonenumber: param.Phonenumber,
		Sex:         param.Sex,
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update personal password
func (*UserController) UserProfileUpdatePwd(ctx *gin.Context) {
	var param dto.UserProfileUpdatePwdRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UserProfileUpdatePwdValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	user := (&service.UserService{}).GetUserByUserId(security.GetAuthUserId(ctx))
	if !password.Verify(user.Password, param.OldPassword) {
		response.NewError().SetMsg("Incorrect old password").Json(ctx)
		return
	}

	if err := (&service.UserService{}).UpdateUser(dto.SaveUser{
		UserId:   user.UserId,
		Password: password.Generate(param.NewPassword),
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Upload avatar
func (*UserController) UserProfileUpdateAvatar(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("avatarfile")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	fileResult, err := upload.New(
		upload.SetLimitType([]string{
			"image/jpeg",
			"image/png",
			"image/svg+xml",
		}),
	).SetFile(&upload.File{
		FileName:    fileHeader.Filename,
		FileType:    fileHeader.Header.Get("Content-Type"),
		FileHeader:  fileHeader.Header,
		FileContent: fileContent,
	}).Save()
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	imgUrl := "/" + fileResult.UrlPath + fileResult.FileName

	if err = (&service.UserService{}).UpdateUser(dto.SaveUser{
		UserId: security.GetAuthUserId(ctx),
		Avatar: imgUrl,
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().SetData("imgUrl", imgUrl).Json(ctx)
}
