package systemcontroller

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/password"
	"mira/common/upload"
	"mira/common/utils"
	"mira/common/xerrors"
	"mira/config"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
	excelize "github.com/xuri/excelize/v2"
)

// UserController handles user-related operations.
type UserController struct {
	UserService   *service.UserService
	DeptService   *service.DeptService
	RoleService   *service.RoleService
	PostService   *service.PostService
	ConfigService *service.ConfigService
}

// NewUserController creates a new UserController.
func NewUserController(userService *service.UserService, deptService *service.DeptService, roleService *service.RoleService, postService *service.PostService, configService *service.ConfigService) *UserController {
	return &UserController{
		UserService:   userService,
		DeptService:   deptService,
		RoleService:   roleService,
		PostService:   postService,
		ConfigService: configService,
	}
}

// DeptTree retrieves the department tree for the current user.
// @Summary Get department tree
// @Description Retrieves the department tree accessible to the current user.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.DeptTreeResponse} "Success"
// @Router /system/user/deptTree [get]
func (c *UserController) DeptTree(ctx *gin.Context) {
	depts := c.DeptService.GetUserDeptTree(security.GetAuthUserId(ctx))

	tree := c.UserService.DeptListToTree(depts, 0)

	response.NewSuccess().SetData("data", tree).Json(ctx)
}

// List retrieves a paginated list of users.
// @Summary Get user list
// @Description Retrieves a paginated list of users based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.UserListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.UserListResponse}} "Success"
// @Router /system/user/list [get]
func (c *UserController) List(ctx *gin.Context) {
	var param dto.UserListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	users, total := c.UserService.GetUserList(param, security.GetAuthUserId(ctx), true)

	for key, user := range users {
		users[key].Dept.DeptName = user.DeptName
		users[key].Dept.Leader = user.Leader
	}

	response.NewSuccess().SetPageData(users, total).Json(ctx)
}

// Detail retrieves the details of a specific user.
// @Summary User details
// @Description Retrieves comprehensive details for a specific user, including roles, posts, and department info.
// @Tags System
// @Accept json
// @Produce json
// @Param userId path int false "User ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "Success"
// @Router /system/user/{userId} [get]
func (c *UserController) Detail(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))

	resp := response.NewSuccess()

	if userId > 0 {
		user := c.UserService.GetUserByUserId(userId)
		user.Admin = user.UserId == 1
		dept := c.DeptService.GetDeptByDeptId(user.DeptId)
		roles, err := c.RoleService.GetRoleListByUserId(user.UserId)
		if err != nil {
			response.NewError().SetMsg(fmt.Sprintf("获取用户角色列表失败: %v", err)).Json(ctx)
			return
		}

		resp.SetData("data", dto.AuthUserInfoResponse{
			UserDetailResponse: user,
			Dept:               dept,
			Roles:              roles,
		})

		roleIds := make([]int, 0)
		for _, role := range roles {
			roleIds = append(roleIds, role.RoleId)
		}
		resp.SetData("roleIds", roleIds)

		postIds := c.PostService.GetPostIdsByUserId(user.UserId)
		resp.SetData("postIds", postIds)
	}

	roles, _ := c.RoleService.GetRoleList(dto.RoleListRequest{}, false)
	if userId != 1 {
		roles = utils.Filter(roles, func(role dto.RoleListResponse) bool {
			return role.RoleId != 1
		})
	}
	resp.SetData("roles", roles)

	posts, _ := c.PostService.GetPostList(dto.PostListRequest{}, false)
	resp.SetData("posts", posts)

	resp.Json(ctx)
}

// Create adds a new user.
// @Summary Add user
// @Description Adds a new user to the system with specified roles and posts.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreateUserRequest true "User data"
// @Success 200 {object} response.Response "Success"
// @Router /system/user [post]
func (c *UserController) Create(ctx *gin.Context) {
	var param dto.CreateUserRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreateUserValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if user := c.UserService.GetUserByUsername(param.UserName); user.UserId > 0 {
		response.NewError().SetMsg("Failed to add user " + param.UserName + ", username already exists").Json(ctx)
		return
	}

	if param.Email != "" {
		if user := c.UserService.GetUserByEmail(param.Email); user.UserId > 0 {
			response.NewError().SetMsg("Failed to add user " + param.UserName + ", email already exists").Json(ctx)
			return
		}
	}

	if param.Phonenumber != "" {
		if user := c.UserService.GetUserByPhonenumber(param.Phonenumber); user.UserId > 0 {
			response.NewError().SetMsg("Failed to add user " + param.UserName + ", phone number already exists").Json(ctx)
			return
		}
	}

	hashedPassword, err := password.Generate(param.Password)
	if err != nil {
		response.NewError().SetCode(500).SetMsg("Failed to process password").Json(ctx)
		return
	}
	if err := c.UserService.CreateUser(dto.SaveUser{
		DeptId:      param.DeptId,
		UserName:    param.UserName,
		NickName:    param.NickName,
		Email:       param.Email,
		Phonenumber: param.Phonenumber,
		Sex:         param.Sex,
		Password:    hashedPassword,
		Status:      param.Status,
		Remark:      param.Remark,
		CreateBy:    security.GetAuthUserName(ctx),
	}, param.RoleIds, param.PostIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update modifies an existing user.
// @Summary Update user
// @Description Modifies an existing user's information, roles, and posts.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateUserRequest true "User data"
// @Success 200 {object} response.Response "Success"
// @Router /system/user [put]
func (c *UserController) Update(ctx *gin.Context) {
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
		if user := c.UserService.GetUserByEmail(param.Email); user.UserId > 0 && user.UserId != param.UserId {
			response.NewError().SetMsg("Failed to modify user " + param.UserName + ", email already exists").Json(ctx)
			return
		}
	}

	if param.Phonenumber != "" {
		if user := c.UserService.GetUserByPhonenumber(param.Phonenumber); user.UserId > 0 && user.UserId != param.UserId {
			response.NewError().SetMsg("Failed to modify user " + param.UserName + ", phone number already exists").Json(ctx)
			return
		}
	}

	if err := c.UserService.UpdateUser(dto.SaveUser{
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

// Remove deletes one or more users.
// @Summary Delete user
// @Description Deletes users by their IDs.
// @Tags System
// @Accept json
// @Produce json
// @Param userIds path string true "User IDs, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /system/user/{userIds} [delete]
func (c *UserController) Remove(ctx *gin.Context) {
	userIds, err := utils.StringToIntSlice(ctx.Param("userIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = validator.RemoveUserValidator(userIds, security.GetAuthUserId(ctx)); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = c.UserService.DeleteUser(userIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// ChangeStatus changes the status of a user.
// @Summary Change user status
// @Description Changes the status (e.g., active/inactive) of a user.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateUserRequest true "User status data"
// @Success 200 {object} response.Response "Success"
// @Router /system/user/changeStatus [put]
func (c *UserController) ChangeStatus(ctx *gin.Context) {
	var param dto.UpdateUserRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.ChangeUserStatusValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := c.UserService.UpdateUser(dto.SaveUser{
		UserId:   param.UserId,
		Status:   param.Status,
		UpdateBy: security.GetAuthUserName(ctx),
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// ResetPwd resets a user's password.
// @Summary Reset user password
// @Description Resets the password for a specific user.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateUserRequest true "User password data"
// @Success 200 {object} response.Response "Success"
// @Router /system/user/resetPwd [put]
func (c *UserController) ResetPwd(ctx *gin.Context) {
	var param dto.UpdateUserRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.ResetUserPwdValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	hashedPassword, err := password.Generate(param.Password)
	if err != nil {
		response.NewError().SetCode(500).SetMsg("Failed to process password").Json(ctx)
		return
	}
	if err := c.UserService.UpdateUser(dto.SaveUser{
		UserId:   param.UserId,
		Password: hashedPassword,
		UpdateBy: security.GetAuthUserName(ctx),
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// AuthRole retrieves authorized roles for a user.
// @Summary Get authorized roles by user ID
// @Description Retrieves a list of all roles, indicating which are assigned to a specific user.
// @Tags System
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "Success"
// @Router /system/user/authRole/{userId} [get]
func (c *UserController) AuthRole(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Param("userId"))

	resp := response.NewSuccess()

	var userHasRoleIds []int

	if userId > 0 {
		user := c.UserService.GetUserByUserId(userId)
		user.Admin = user.UserId == 1
		dept := c.DeptService.GetDeptByDeptId(user.DeptId)
		roles, err := c.RoleService.GetRoleListByUserId(user.UserId)
		if err != nil {
			response.NewError().SetMsg(fmt.Sprintf("Failed to get user role list: %v", err)).Json(ctx)
			return
		}
		for _, role := range roles {
			userHasRoleIds = append(userHasRoleIds, role.RoleId)
		}

		resp.SetData("user", dto.AuthUserInfoResponse{
			UserDetailResponse: user,
			Dept:               dept,
			Roles:              roles,
		})
	}

	roles, _ := c.RoleService.GetRoleList(dto.RoleListRequest{}, false)
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
	resp.SetData("roles", roles)

	resp.Json(ctx)
}

// AddAuthRole assigns roles to a user.
// @Summary User authorized role
// @Description Assigns a set of roles to a specific user.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.AddUserAuthRoleRequest true "User and Role IDs"
// @Success 200 {object} response.Response "Success"
// @Router /system/user/authRole [put]
func (c *UserController) AddAuthRole(ctx *gin.Context) {
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

	if err := c.UserService.AddAuthRole(param.UserId, roleIds); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// ImportTemplate provides a template for importing user data.
// @Summary Import user template
// @Description Downloads an Excel template for importing user data.
// @Tags System
// @Produce octet-stream
// @Success 200 {file} file "Excel template"
// @Router /system/user/importTemplate [post]
func (c *UserController) ImportTemplate(ctx *gin.Context) {
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

// ImportData imports user data from an Excel file.
// @Summary Import user data
// @Description Imports user data from an Excel file, with an option to update existing users.
// @Tags System
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Excel file"
// @Param updateSupport query bool false "Whether to update existing user data"
// @Success 200 {object} response.Response "Success"
// @Router /system/user/importData [post]
func (c *UserController) ImportData(ctx *gin.Context) {
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
		user := c.UserService.GetUserByUsername(item.UserName)

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
				failNum++
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be added: "+err.Error())
				continue
			}
			hashedPassword, err := password.Generate(c.ConfigService.GetConfigCacheByConfigKey("sys.user.initPassword").ConfigValue)
			if err != nil {
				failNum++
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be added: "+err.Error())
				continue
			}
			if err = c.UserService.CreateUser(dto.SaveUser{
				DeptId:      item.DeptId,
				UserName:    item.UserName,
				NickName:    item.NickName,
				Email:       item.Email,
				Phonenumber: item.Phonenumber,
				Sex:         item.Sex,
				Password:    hashedPassword,
				Status:      item.Status,
				CreateBy:    authUserName,
			}, nil, nil); err != nil {
				failNum++
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be added: "+err.Error())
				continue
			}
			successNum++
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
				failNum++
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be updated: "+err.Error())
				continue
			}
			// Update existing user
			if err = c.UserService.UpdateUser(dto.SaveUser{
				UserId:      user.UserId,
				DeptId:      item.DeptId,
				NickName:    item.NickName,
				Email:       item.Email,
				Phonenumber: item.Phonenumber,
				Sex:         item.Sex,
				Status:      item.Status,
				UpdateBy:    authUserName,
			}, nil, nil); err != nil {
				failNum++
				failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" failed to be updated: "+err.Error())
				continue
			}
			successNum++
			continue
		} else {
			failNum++
			failMsg = append(failMsg, strconv.Itoa(failNum)+", Account "+item.UserName+" already exists")
		}
	}

	if failNum > 0 {
		response.NewError().SetMsg("Import failed, " + strconv.Itoa(failNum) + " pieces of data are wrong, the errors are as follows:" + strings.Join(failMsg, "<br/>")).Json(ctx)
		return
	}

	response.NewSuccess().SetMsg("Import successful, " + strconv.Itoa(successNum) + " pieces of data in total").Json(ctx)
}

// Export exports user data to an Excel file.
// @Summary Export user data
// @Description Exports user data to an Excel file based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.UserListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /system/user/export [post]
func (c *UserController) Export(ctx *gin.Context) {
	var param dto.UserListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.UserExportResponse, 0)

	users, _ := c.UserService.GetUserList(param, security.GetAuthUserId(ctx), false)
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

// GetProfile retrieves the profile of the currently authenticated user.
// @Summary Get personal information
// @Description Retrieves the profile, roles, and posts of the currently authenticated user.
// @Tags System
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "Success"
// @Router /system/user/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
	user := c.UserService.GetUserByUserId(security.GetAuthUserId(ctx))
	user.Admin = user.UserId == 1
	dept := c.DeptService.GetDeptByDeptId(user.DeptId)
	roles, err := c.RoleService.GetRoleListByUserId(user.UserId)
	if err != nil {
		response.NewError().SetMsg(fmt.Sprintf("Failed to get user role list: %v", err)).Json(ctx)
		return
	}

	data := dto.AuthUserInfoResponse{
		UserDetailResponse: user,
		Dept:               dept,
		Roles:              roles,
	}

	// Get role group
	roleGroup, err := c.RoleService.GetRoleNamesByUserId(user.UserId)
	if err != nil {
		response.NewError().SetMsg(fmt.Sprintf("Failed to get user role names: %v", err)).Json(ctx)
		return
	}

	// Get post group
	postGroup := c.PostService.GetPostNamesByUserId(user.UserId)

	response.NewSuccess().
		SetData("data", data).
		SetData("roleGroup", strings.Join(roleGroup, ",")).
		SetData("postGroup", strings.Join(postGroup, ",")).
		Json(ctx)
}

// UpdateProfile updates the profile of the currently authenticated user.
// @Summary Update personal information
// @Description Updates the profile information of the currently authenticated user.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdateProfileRequest true "Profile data"
// @Success 200 {object} response.Response "Success"
// @Router /system/user/profile [put]
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	var param dto.UpdateProfileRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdateProfileValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := c.UserService.UpdateUser(dto.SaveUser{
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

// UserProfileUpdatePwd updates the password of the currently authenticated user.
// @Summary Update personal password
// @Description Updates the password for the currently authenticated user.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UserProfileUpdatePwdRequest true "Password data"
// @Success 200 {object} response.Response "Success"
// @Router /system/user/profile/updatePwd [put]
func (c *UserController) UserProfileUpdatePwd(ctx *gin.Context) {
	var param dto.UserProfileUpdatePwdRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UserProfileUpdatePwdValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	user := c.UserService.GetUserByUserId(security.GetAuthUserId(ctx))
	if err := password.Verify(user.Password, param.OldPassword); err != nil {
		if err == xerrors.ErrMismatchedPassword {
			response.NewError().SetMsg("Incorrect old password").Json(ctx)
			return
		}
		response.NewError().SetCode(500).SetMsg("Failed to verify password").Json(ctx)
		return
	}

	hashedPassword, err := password.Generate(param.NewPassword)
	if err != nil {
		response.NewError().SetCode(500).SetMsg("Failed to process password").Json(ctx)
		return
	}
	if err := c.UserService.UpdateUser(dto.SaveUser{
		UserId:   user.UserId,
		Password: hashedPassword,
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// UserProfileUpdateAvatar updates the avatar of the currently authenticated user.
// @Summary Upload avatar
// @Description Updates the avatar for the currently authenticated user.
// @Tags System
// @Accept multipart/form-data
// @Produce json
// @Param avatarfile formData file true "Avatar file"
// @Success 200 {object} response.Response{data=map[string]string} "Success"
// @Router /system/user/profile/avatar [post]
func (c *UserController) UserProfileUpdateAvatar(ctx *gin.Context) {
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
	defer file.Close()

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

	if err = c.UserService.UpdateUser(dto.SaveUser{
		UserId: security.GetAuthUserId(ctx),
		Avatar: imgUrl,
	}, nil, nil); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().SetData("imgUrl", imgUrl).Json(ctx)
}
