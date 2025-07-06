package router

import (
	"mira/app"
	"mira/app/controller"
	"mira/app/middleware"
	"mira/common/types/constant"

	"github.com/gin-gonic/gin"
)

// Admin router group
func RegisterAdminGroupApi(api *gin.RouterGroup, container *app.AppContainer) {
	api.Use(middleware.Cors()) // CORS middleware

	registerAuthRoutes(api, container)
	registerSystemRoutes(api, container)
	registerMonitorRoutes(api, container)
}

func registerAuthRoutes(api *gin.RouterGroup, container *app.AppContainer) {
	authController := &controller.AuthController{}
	api.GET("/captchaImage", authController.CaptchaImage)
	api.POST("/register", authController.Register)
	api.POST("/login", container.LogininforMiddleware(), authController.Login)
	api.POST("/logout", authController.Logout)

	// Enable authentication middleware. The following routes require authentication.
	api.Use(middleware.AuthMiddleware())

	api.GET("/getInfo", authController.GetInfo)
	api.GET("/getRouters", authController.GetRouters)
}

func registerSystemRoutes(api *gin.RouterGroup, container *app.AppContainer) {
	// User Routes
	userGroup := api.Group("/system/user")
	{
		userGroup.GET("/profile", container.UserController.GetProfile)
		userGroup.PUT("/profile", container.UserController.UpdateProfile)
		userGroup.PUT("/profile/updatePwd", container.UserController.UserProfileUpdatePwd)
		userGroup.POST("/profile/avatar", container.UserController.UserProfileUpdateAvatar)
		userGroup.GET("/deptTree", container.HasPerm("system:user:list"), container.UserController.DeptTree)
		userGroup.GET("/list", container.HasPerm("system:user:list"), container.UserController.List)
		userGroup.GET("/", container.HasPerm("system:user:query"), container.UserController.Detail)
		userGroup.GET("/:userId", container.HasPerm("system:user:query"), container.UserController.Detail)
		userGroup.GET("/authRole/:userId", container.HasPerm("system:user:query"), container.UserController.AuthRole)
		userGroup.POST("", container.HasPerm("system:user:add"), container.OperLogMiddleware("Add User", constant.REQUEST_BUSINESS_TYPE_INSERT), container.UserController.Create)
		userGroup.PUT("", container.HasPerm("system:user:edit"), container.OperLogMiddleware("Update User", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.UserController.Update)
		userGroup.DELETE("/:userIds", container.HasPerm("system:user:remove"), container.OperLogMiddleware("Delete User", constant.REQUEST_BUSINESS_TYPE_DELETE), container.UserController.Remove)
		userGroup.PUT("/changeStatus", container.HasPerm("system:user:edit"), container.OperLogMiddleware("Modify User Status", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.UserController.ChangeStatus)
		userGroup.PUT("/resetPwd", container.HasPerm("system:user:edit"), container.OperLogMiddleware("Modify User Password", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.UserController.ResetPwd)
		userGroup.PUT("/authRole", container.HasPerm("system:user:edit"), container.OperLogMiddleware("User Authorized Role", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.UserController.AddAuthRole)
		userGroup.POST("/export", container.HasPerm("system:user:export"), container.OperLogMiddleware("Export User", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.UserController.Export)
		userGroup.POST("/importData", container.HasPerm("system:user:import"), container.OperLogMiddleware("Import User", constant.REQUEST_BUSINESS_TYPE_IMPORT), container.UserController.ImportData)
		userGroup.POST("/importTemplate", container.OperLogMiddleware("Import User Template", constant.REQUEST_BUSINESS_TYPE_IMPORT), container.UserController.ImportTemplate)
	}

	// Role Routes
	roleGroup := api.Group("/system/role")
	{
		roleGroup.GET("/list", container.HasPerm("system:role:list"), container.RoleController.List)
		roleGroup.GET("/:roleId", container.HasPerm("system:role:query"), container.RoleController.Detail)
		roleGroup.GET("/deptTree/:roleId", container.HasPerm("system:role:query"), container.RoleController.DeptTree)
		roleGroup.GET("/authUser/allocatedList", container.HasPerm("system:role:list"), container.RoleController.RoleAuthUserAllocatedList)
		roleGroup.GET("/authUser/unallocatedList", container.HasPerm("system:role:list"), container.RoleController.RoleAuthUserUnallocatedList)
		roleGroup.POST("", container.HasPerm("system:role:add"), container.OperLogMiddleware("Add Role", constant.REQUEST_BUSINESS_TYPE_INSERT), container.RoleController.Create)
		roleGroup.PUT("", container.HasPerm("system:role:edit"), container.OperLogMiddleware("Update Role", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.RoleController.Update)
		roleGroup.DELETE("/:roleIds", container.HasPerm("system:role:remove"), container.OperLogMiddleware("Delete Role", constant.REQUEST_BUSINESS_TYPE_DELETE), container.RoleController.Remove)
		roleGroup.PUT("/changeStatus", container.HasPerm("system:role:edit"), container.OperLogMiddleware("Modify Role Status", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.RoleController.ChangeStatus)
		roleGroup.PUT("/dataScope", container.HasPerm("system:role:edit"), container.OperLogMiddleware("Assign Data Permissions", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.RoleController.DataScope)
		roleGroup.PUT("/authUser/selectAll", container.HasPerm("system:role:edit"), container.OperLogMiddleware("Batch Select User Authorization", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.RoleController.RoleAuthUserSelectAll)
		roleGroup.PUT("/authUser/cancel", container.HasPerm("system:role:edit"), container.OperLogMiddleware("Cancel Authorized User", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.RoleController.RoleAuthUserCancel)
		roleGroup.PUT("/authUser/cancelAll", container.HasPerm("system:role:edit"), container.OperLogMiddleware("Batch Cancel Authorized User", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.RoleController.RoleAuthUserCancelAll)
		roleGroup.POST("/export", container.HasPerm("system:role:export"), container.OperLogMiddleware("Export Role", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.RoleController.Export)
	}

	// Menu Routes
	menuGroup := api.Group("/system/menu")
	{
		menuGroup.GET("/list", container.HasPerm("system:menu:list"), container.MenuController.List)
		menuGroup.GET("/treeselect", container.MenuController.Treeselect)
		menuGroup.GET("/roleMenuTreeselect/:roleId", container.MenuController.RoleMenuTreeselect)
		menuGroup.GET("/:menuId", container.HasPerm("system:menu:query"), container.MenuController.Detail)
		menuGroup.POST("", container.HasPerm("system:menu:add"), container.OperLogMiddleware("Add Menu", constant.REQUEST_BUSINESS_TYPE_INSERT), container.MenuController.Create)
		menuGroup.PUT("", container.HasPerm("system:menu:edit"), container.OperLogMiddleware("Update Menu", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.MenuController.Update)
		menuGroup.DELETE("/:menuId", container.HasPerm("system:menu:remove"), container.OperLogMiddleware("Delete Menu", constant.REQUEST_BUSINESS_TYPE_DELETE), container.MenuController.Remove)
	}

	// Dept Routes
	deptGroup := api.Group("/system/dept")
	{
		deptGroup.GET("/list", container.HasPerm("system:dept:list"), container.DeptController.List)
		deptGroup.GET("/list/exclude/:deptId", container.HasPerm("system:dept:list"), container.DeptController.ListExclude)
		deptGroup.GET("/:deptId", container.HasPerm("system:dept:query"), container.DeptController.Detail)
		deptGroup.POST("", container.HasPerm("system:dept:add"), container.OperLogMiddleware("Add Department", constant.REQUEST_BUSINESS_TYPE_INSERT), container.DeptController.Create)
		deptGroup.PUT("", container.HasPerm("system:dept:edit"), container.OperLogMiddleware("Update Department", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.DeptController.Update)
		deptGroup.DELETE("/:deptId", container.HasPerm("system:dept:remove"), container.OperLogMiddleware("Delete Department", constant.REQUEST_BUSINESS_TYPE_DELETE), container.DeptController.Remove)
	}

	// Post Routes
	postGroup := api.Group("/system/post")
	{
		postGroup.GET("/list", container.HasPerm("system:post:list"), container.PostController.List)
		postGroup.GET("/:postId", container.HasPerm("system:post:query"), container.PostController.Detail)
		postGroup.POST("", container.HasPerm("system:post:add"), container.OperLogMiddleware("Add Post", constant.REQUEST_BUSINESS_TYPE_INSERT), container.PostController.Create)
		postGroup.PUT("", container.HasPerm("system:post:edit"), container.OperLogMiddleware("Update Post", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.PostController.Update)
		postGroup.DELETE("/:postIds", container.HasPerm("system:post:remove"), container.OperLogMiddleware("Delete Post", constant.REQUEST_BUSINESS_TYPE_DELETE), container.PostController.Remove)
		postGroup.POST("/export", container.HasPerm("system:post:export"), container.OperLogMiddleware("Export Post", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.PostController.Export)
	}

	// Dict Routes
	dictGroup := api.Group("/system/dict")
	{
		dictGroup.GET("/type/list", container.HasPerm("system:dict:list"), container.DictTypeController.List)
		dictGroup.GET("/type/:dictId", container.HasPerm("system:dict:query"), container.DictTypeController.Detail)
		dictGroup.GET("/type/optionselect", container.DictTypeController.Optionselect)
		dictGroup.POST("/type", container.HasPerm("system:dict:add"), container.OperLogMiddleware("Add Dictionary Type", constant.REQUEST_BUSINESS_TYPE_INSERT), container.DictTypeController.Create)
		dictGroup.PUT("/type", container.HasPerm("system:dict:edit"), container.OperLogMiddleware("Update Dictionary Type", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.DictTypeController.Update)
		dictGroup.DELETE("/type/:dictIds", container.HasPerm("system:dict:remove"), container.OperLogMiddleware("Delete Dictionary Type", constant.REQUEST_BUSINESS_TYPE_DELETE), container.DictTypeController.Remove)
		dictGroup.POST("/type/export", container.HasPerm("system:dict:export"), container.OperLogMiddleware("Export Dictionary Type", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.DictTypeController.Export)
		dictGroup.DELETE("/type/refreshCache", container.HasPerm("system:dict:remove"), container.OperLogMiddleware("Refresh Dictionary Type Cache", constant.REQUEST_BUSINESS_TYPE_DELETE), container.DictTypeController.RefreshCache)

		dictGroup.GET("/data/list", container.HasPerm("system:dict:list"), container.DictDataController.List)
		dictGroup.GET("/data/:dictCode", container.HasPerm("system:dict:query"), container.DictDataController.Detail)
		dictGroup.GET("/data/type/:dictType", container.DictDataController.Type)
		dictGroup.POST("/data", container.HasPerm("system:dict:add"), container.OperLogMiddleware("Add Dictionary Data", constant.REQUEST_BUSINESS_TYPE_INSERT), container.DictDataController.Create)
		dictGroup.PUT("/data", container.HasPerm("system:dict:edit"), container.OperLogMiddleware("Update Dictionary Data", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.DictDataController.Update)
		dictGroup.DELETE("/data/:dictCodes", container.HasPerm("system:dict:remove"), container.OperLogMiddleware("Delete Dictionary Data", constant.REQUEST_BUSINESS_TYPE_DELETE), container.DictDataController.Remove)
		dictGroup.POST("/data/export", container.HasPerm("system:dict:export"), container.OperLogMiddleware("Export Dictionary Data", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.DictDataController.Export)
	}

	// Config Routes
	configGroup := api.Group("/system/config")
	{
		configGroup.GET("/list", container.HasPerm("system:config:list"), container.ConfigController.List)
		configGroup.GET("/:configId", container.HasPerm("system:config:query"), container.ConfigController.Detail)
		configGroup.GET("/configKey/:configKey", container.ConfigController.ConfigKey)
		configGroup.POST("", container.HasPerm("system:config:add"), container.OperLogMiddleware("Add Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_INSERT), container.ConfigController.Create)
		configGroup.PUT("", container.HasPerm("system:config:edit"), container.OperLogMiddleware("Update Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_UPDATE), container.ConfigController.Update)
		configGroup.DELETE("/:configIds", container.HasPerm("system:config:remove"), container.OperLogMiddleware("Delete Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_DELETE), container.ConfigController.Remove)
		configGroup.POST("/export", container.HasPerm("system:config:export"), container.OperLogMiddleware("Export Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.ConfigController.Export)
		configGroup.DELETE("/refreshCache", container.HasPerm("system:config:remove"), container.OperLogMiddleware("Refresh Parameter Configuration Cache", constant.REQUEST_BUSINESS_TYPE_DELETE), container.ConfigController.RefreshCache)
	}
}

func registerMonitorRoutes(api *gin.RouterGroup, container *app.AppContainer) {
	monitorGroup := api.Group("/monitor")
	{
		logininforGroup := monitorGroup.Group("/logininfor")
		{
			logininforGroup.GET("/list", container.HasPerm("monitor:logininfor:list"), container.LogininforController.List)
			logininforGroup.DELETE("/:infoIds", container.HasPerm("monitor:logininfor:remove"), container.OperLogMiddleware("Delete Login Log", constant.REQUEST_BUSINESS_TYPE_DELETE), container.LogininforController.Remove)
			logininforGroup.DELETE("/clean", container.HasPerm("monitor:logininfor:remove"), container.OperLogMiddleware("Clear Login Log", constant.REQUEST_BUSINESS_TYPE_DELETE), container.LogininforController.Clean)
			logininforGroup.GET("/unlock/:userName", container.HasPerm("monitor:logininfor:unlock"), container.OperLogMiddleware("Account Unlock", constant.REQUEST_BUSINESS_TYPE_DELETE), container.LogininforController.Unlock)
			logininforGroup.POST("/export", container.HasPerm("monitor:logininfor:export"), container.OperLogMiddleware("Export Login Log", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.LogininforController.Export)
		}
		operlogGroup := monitorGroup.Group("/operlog")
		{
			operlogGroup.GET("/list", container.HasPerm("monitor:operlog:list"), container.OperlogController.List)
			operlogGroup.DELETE("/:operIds", container.HasPerm("monitor:operlog:remove"), container.OperLogMiddleware("Delete Operation Log", constant.REQUEST_BUSINESS_TYPE_DELETE), container.OperlogController.Remove)
			operlogGroup.DELETE("/clean", container.HasPerm("monitor:operlog:remove"), container.OperLogMiddleware("Clear Operation Log", constant.REQUEST_BUSINESS_TYPE_DELETE), container.OperlogController.Clean)
			operlogGroup.POST("/export", container.HasPerm("monitor:operlog:export"), container.OperLogMiddleware("Export Operation Log", constant.REQUEST_BUSINESS_TYPE_EXPORT), container.OperlogController.Export)
		}
	}
}
