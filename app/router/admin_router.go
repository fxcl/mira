package router

import (
	"mira/app/controller"
	monitorcontroller "mira/app/controller/monitor"
	systemcontroller "mira/app/controller/system"
	"mira/app/middleware"
	"mira/app/service"
	"mira/common/types/constant"

	"github.com/gin-gonic/gin"
)

// Admin router group
func RegisterAdminGroupApi(api *gin.RouterGroup) {
	api.Use(middleware.Cors()) // CORS middleware

	registerAuthRoutes(api)
	registerSystemRoutes(api)
	registerMonitorRoutes(api)
}

func registerAuthRoutes(api *gin.RouterGroup) {
	authController := &controller.AuthController{}
	api.GET("/captchaImage", authController.CaptchaImage)
	api.POST("/register", authController.Register)
	api.POST("/login", middleware.LogininforMiddleware(), authController.Login)
	api.POST("/logout", authController.Logout)

	// Enable authentication middleware. The following routes require authentication.
	api.Use(middleware.AuthMiddleware())

	api.GET("/getInfo", authController.GetInfo)
	api.GET("/getRouters", authController.GetRouters)
}

func registerSystemRoutes(api *gin.RouterGroup) {
	// Instantiate services
	userService := &service.UserService{}
	deptService := &service.DeptService{}
	roleService := &service.RoleService{}
	postService := &service.PostService{}
	menuService := &service.MenuService{}
	configService := &service.ConfigService{}
	dictTypeService := &service.DictTypeService{}
	dictDataService := &service.DictDataService{}

	// Instantiate controllers with dependencies
	userController := systemcontroller.NewUserController(userService, deptService, roleService, postService, configService)
	roleController := systemcontroller.NewRoleController(roleService, deptService, userService)
	menuController := systemcontroller.NewMenuController(menuService)
	deptController := systemcontroller.NewDeptController(deptService, userService)
	postController := systemcontroller.NewPostController(postService)
	dictTypeController := systemcontroller.NewDictTypeController(dictTypeService)
	dictDataController := systemcontroller.NewDictDataController(dictDataService)
	configController := systemcontroller.NewConfigController(configService)

	// User Routes
	userGroup := api.Group("/system/user")
	{
		userGroup.GET("/profile", userController.GetProfile)
		userGroup.PUT("/profile", userController.UpdateProfile)
		userGroup.PUT("/profile/updatePwd", userController.UserProfileUpdatePwd)
		userGroup.POST("/profile/avatar", userController.UserProfileUpdateAvatar)
		userGroup.GET("/deptTree", middleware.HasPerm("system:user:list"), userController.DeptTree)
		userGroup.GET("/list", middleware.HasPerm("system:user:list"), userController.List)
		userGroup.GET("/", middleware.HasPerm("system:user:query"), userController.Detail)
		userGroup.GET("/:userId", middleware.HasPerm("system:user:query"), userController.Detail)
		userGroup.GET("/authRole/:userId", middleware.HasPerm("system:user:query"), userController.AuthRole)
		userGroup.POST("", middleware.HasPerm("system:user:add"), middleware.OperLogMiddleware("Add User", constant.REQUEST_BUSINESS_TYPE_INSERT), userController.Create)
		userGroup.PUT("", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("Update User", constant.REQUEST_BUSINESS_TYPE_UPDATE), userController.Update)
		userGroup.DELETE("/:userIds", middleware.HasPerm("system:user:remove"), middleware.OperLogMiddleware("Delete User", constant.REQUEST_BUSINESS_TYPE_DELETE), userController.Remove)
		userGroup.PUT("/changeStatus", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("Modify User Status", constant.REQUEST_BUSINESS_TYPE_UPDATE), userController.ChangeStatus)
		userGroup.PUT("/resetPwd", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("Modify User Password", constant.REQUEST_BUSINESS_TYPE_UPDATE), userController.ResetPwd)
		userGroup.PUT("/authRole", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("User Authorized Role", constant.REQUEST_BUSINESS_TYPE_UPDATE), userController.AddAuthRole)
		userGroup.POST("/export", middleware.HasPerm("system:user:export"), middleware.OperLogMiddleware("Export User", constant.REQUEST_BUSINESS_TYPE_EXPORT), userController.Export)
		userGroup.POST("/importData", middleware.HasPerm("system:user:import"), middleware.OperLogMiddleware("Import User", constant.REQUEST_BUSINESS_TYPE_IMPORT), userController.ImportData)
		userGroup.POST("/importTemplate", middleware.OperLogMiddleware("Import User Template", constant.REQUEST_BUSINESS_TYPE_IMPORT), userController.ImportTemplate)
	}

	// Role Routes
	roleGroup := api.Group("/system/role")
	{
		roleGroup.GET("/list", middleware.HasPerm("system:role:list"), roleController.List)
		roleGroup.GET("/:roleId", middleware.HasPerm("system:role:query"), roleController.Detail)
		roleGroup.GET("/deptTree/:roleId", middleware.HasPerm("system:role:query"), roleController.DeptTree)
		roleGroup.GET("/authUser/allocatedList", middleware.HasPerm("system:role:list"), roleController.RoleAuthUserAllocatedList)
		roleGroup.GET("/authUser/unallocatedList", middleware.HasPerm("system:role:list"), roleController.RoleAuthUserUnallocatedList)
		roleGroup.POST("", middleware.HasPerm("system:role:add"), middleware.OperLogMiddleware("Add Role", constant.REQUEST_BUSINESS_TYPE_INSERT), roleController.Create)
		roleGroup.PUT("", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Update Role", constant.REQUEST_BUSINESS_TYPE_UPDATE), roleController.Update)
		roleGroup.DELETE("/:roleIds", middleware.HasPerm("system:role:remove"), middleware.OperLogMiddleware("Delete Role", constant.REQUEST_BUSINESS_TYPE_DELETE), roleController.Remove)
		roleGroup.PUT("/changeStatus", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Modify Role Status", constant.REQUEST_BUSINESS_TYPE_UPDATE), roleController.ChangeStatus)
		roleGroup.PUT("/dataScope", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Assign Data Permissions", constant.REQUEST_BUSINESS_TYPE_UPDATE), roleController.DataScope)
		roleGroup.PUT("/authUser/selectAll", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Batch Select User Authorization", constant.REQUEST_BUSINESS_TYPE_UPDATE), roleController.RoleAuthUserSelectAll)
		roleGroup.PUT("/authUser/cancel", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Cancel Authorized User", constant.REQUEST_BUSINESS_TYPE_UPDATE), roleController.RoleAuthUserCancel)
		roleGroup.PUT("/authUser/cancelAll", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Batch Cancel Authorized User", constant.REQUEST_BUSINESS_TYPE_UPDATE), roleController.RoleAuthUserCancelAll)
		roleGroup.POST("/export", middleware.HasPerm("system:role:export"), middleware.OperLogMiddleware("Export Role", constant.REQUEST_BUSINESS_TYPE_EXPORT), roleController.Export)
	}

	// Menu Routes
	menuGroup := api.Group("/system/menu")
	{
		menuGroup.GET("/list", middleware.HasPerm("system:menu:list"), menuController.List)
		menuGroup.GET("/treeselect", menuController.Treeselect)
		menuGroup.GET("/roleMenuTreeselect/:roleId", menuController.RoleMenuTreeselect)
		menuGroup.GET("/:menuId", middleware.HasPerm("system:menu:query"), menuController.Detail)
		menuGroup.POST("", middleware.HasPerm("system:menu:add"), middleware.OperLogMiddleware("Add Menu", constant.REQUEST_BUSINESS_TYPE_INSERT), menuController.Create)
		menuGroup.PUT("", middleware.HasPerm("system:menu:edit"), middleware.OperLogMiddleware("Update Menu", constant.REQUEST_BUSINESS_TYPE_UPDATE), menuController.Update)
		menuGroup.DELETE("/:menuId", middleware.HasPerm("system:menu:remove"), middleware.OperLogMiddleware("Delete Menu", constant.REQUEST_BUSINESS_TYPE_DELETE), menuController.Remove)
	}

	// Dept Routes
	deptGroup := api.Group("/system/dept")
	{
		deptGroup.GET("/list", middleware.HasPerm("system:dept:list"), deptController.List)
		deptGroup.GET("/list/exclude/:deptId", middleware.HasPerm("system:dept:list"), deptController.ListExclude)
		deptGroup.GET("/:deptId", middleware.HasPerm("system:dept:query"), deptController.Detail)
		deptGroup.POST("", middleware.HasPerm("system:dept:add"), middleware.OperLogMiddleware("Add Department", constant.REQUEST_BUSINESS_TYPE_INSERT), deptController.Create)
		deptGroup.PUT("", middleware.HasPerm("system:dept:edit"), middleware.OperLogMiddleware("Update Department", constant.REQUEST_BUSINESS_TYPE_UPDATE), deptController.Update)
		deptGroup.DELETE("/:deptId", middleware.HasPerm("system:dept:remove"), middleware.OperLogMiddleware("Delete Department", constant.REQUEST_BUSINESS_TYPE_DELETE), deptController.Remove)
	}

	// Post Routes
	postGroup := api.Group("/system/post")
	{
		postGroup.GET("/list", middleware.HasPerm("system:post:list"), postController.List)
		postGroup.GET("/:postId", middleware.HasPerm("system:post:query"), postController.Detail)
		postGroup.POST("", middleware.HasPerm("system:post:add"), middleware.OperLogMiddleware("Add Post", constant.REQUEST_BUSINESS_TYPE_INSERT), postController.Create)
		postGroup.PUT("", middleware.HasPerm("system:post:edit"), middleware.OperLogMiddleware("Update Post", constant.REQUEST_BUSINESS_TYPE_UPDATE), postController.Update)
		postGroup.DELETE("/:postIds", middleware.HasPerm("system:post:remove"), middleware.OperLogMiddleware("Delete Post", constant.REQUEST_BUSINESS_TYPE_DELETE), postController.Remove)
		postGroup.POST("/export", middleware.HasPerm("system:post:export"), middleware.OperLogMiddleware("Export Post", constant.REQUEST_BUSINESS_TYPE_EXPORT), postController.Export)
	}

	// Dict Routes
	dictGroup := api.Group("/system/dict")
	{
		dictGroup.GET("/type/list", middleware.HasPerm("system:dict:list"), dictTypeController.List)
		dictGroup.GET("/type/:dictId", middleware.HasPerm("system:dict:query"), dictTypeController.Detail)
		dictGroup.GET("/type/optionselect", dictTypeController.Optionselect)
		dictGroup.POST("/type", middleware.HasPerm("system:dict:add"), middleware.OperLogMiddleware("Add Dictionary Type", constant.REQUEST_BUSINESS_TYPE_INSERT), dictTypeController.Create)
		dictGroup.PUT("/type", middleware.HasPerm("system:dict:edit"), middleware.OperLogMiddleware("Update Dictionary Type", constant.REQUEST_BUSINESS_TYPE_UPDATE), dictTypeController.Update)
		dictGroup.DELETE("/type/:dictIds", middleware.HasPerm("system:dict:remove"), middleware.OperLogMiddleware("Delete Dictionary Type", constant.REQUEST_BUSINESS_TYPE_DELETE), dictTypeController.Remove)
		dictGroup.POST("/type/export", middleware.HasPerm("system:dict:export"), middleware.OperLogMiddleware("Export Dictionary Type", constant.REQUEST_BUSINESS_TYPE_EXPORT), dictTypeController.Export)
		dictGroup.DELETE("/type/refreshCache", middleware.HasPerm("system:dict:remove"), middleware.OperLogMiddleware("Refresh Dictionary Type Cache", constant.REQUEST_BUSINESS_TYPE_DELETE), dictTypeController.RefreshCache)

		dictGroup.GET("/data/list", middleware.HasPerm("system:dict:list"), dictDataController.List)
		dictGroup.GET("/data/:dictCode", middleware.HasPerm("system:dict:query"), dictDataController.Detail)
		dictGroup.GET("/data/type/:dictType", dictDataController.Type)
		dictGroup.POST("/data", middleware.HasPerm("system:dict:add"), middleware.OperLogMiddleware("Add Dictionary Data", constant.REQUEST_BUSINESS_TYPE_INSERT), dictDataController.Create)
		dictGroup.PUT("/data", middleware.HasPerm("system:dict:edit"), middleware.OperLogMiddleware("Update Dictionary Data", constant.REQUEST_BUSINESS_TYPE_UPDATE), dictDataController.Update)
		dictGroup.DELETE("/data/:dictCodes", middleware.HasPerm("system:dict:remove"), middleware.OperLogMiddleware("Delete Dictionary Data", constant.REQUEST_BUSINESS_TYPE_DELETE), dictDataController.Remove)
		dictGroup.POST("/data/export", middleware.HasPerm("system:dict:export"), middleware.OperLogMiddleware("Export Dictionary Data", constant.REQUEST_BUSINESS_TYPE_EXPORT), dictDataController.Export)
	}

	// Config Routes
	configGroup := api.Group("/system/config")
	{
		configGroup.GET("/list", middleware.HasPerm("system:config:list"), configController.List)
		configGroup.GET("/:configId", middleware.HasPerm("system:config:query"), configController.Detail)
		configGroup.GET("/configKey/:configKey", configController.ConfigKey)
		configGroup.POST("", middleware.HasPerm("system:config:add"), middleware.OperLogMiddleware("Add Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_INSERT), configController.Create)
		configGroup.PUT("", middleware.HasPerm("system:config:edit"), middleware.OperLogMiddleware("Update Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_UPDATE), configController.Update)
		configGroup.DELETE("/:configIds", middleware.HasPerm("system:config:remove"), middleware.OperLogMiddleware("Delete Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_DELETE), configController.Remove)
		configGroup.POST("/export", middleware.HasPerm("system:config:export"), middleware.OperLogMiddleware("Export Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_EXPORT), configController.Export)
		configGroup.DELETE("/refreshCache", middleware.HasPerm("system:config:remove"), middleware.OperLogMiddleware("Refresh Parameter Configuration Cache", constant.REQUEST_BUSINESS_TYPE_DELETE), configController.RefreshCache)
	}
}

func registerMonitorRoutes(api *gin.RouterGroup) {
	logininforService := &service.LogininforService{}
	logininforController := monitorcontroller.NewLogininforController(logininforService)

	operLogService := &service.OperLogService{}
	operlogController := monitorcontroller.NewOperlogController(operLogService)

	monitorGroup := api.Group("/monitor")
	{
		logininforGroup := monitorGroup.Group("/logininfor")
		{
			logininforGroup.GET("/list", middleware.HasPerm("monitor:logininfor:list"), logininforController.List)
			logininforGroup.DELETE("/:infoIds", middleware.HasPerm("monitor:logininfor:remove"), middleware.OperLogMiddleware("Delete Login Log", constant.REQUEST_BUSINESS_TYPE_DELETE), logininforController.Remove)
			logininforGroup.DELETE("/clean", middleware.HasPerm("monitor:logininfor:remove"), middleware.OperLogMiddleware("Clear Login Log", constant.REQUEST_BUSINESS_TYPE_DELETE), logininforController.Clean)
			logininforGroup.GET("/unlock/:userName", middleware.HasPerm("monitor:logininfor:unlock"), middleware.OperLogMiddleware("Account Unlock", constant.REQUEST_BUSINESS_TYPE_DELETE), logininforController.Unlock)
			logininforGroup.POST("/export", middleware.HasPerm("monitor:logininfor:export"), middleware.OperLogMiddleware("Export Login Log", constant.REQUEST_BUSINESS_TYPE_EXPORT), logininforController.Export)
		}
		operlogGroup := monitorGroup.Group("/operlog")
		{
			operlogGroup.GET("/list", middleware.HasPerm("monitor:operlog:list"), operlogController.List)
			operlogGroup.DELETE("/:operIds", middleware.HasPerm("monitor:operlog:remove"), middleware.OperLogMiddleware("Delete Operation Log", constant.REQUEST_BUSINESS_TYPE_DELETE), operlogController.Remove)
			operlogGroup.DELETE("/clean", middleware.HasPerm("monitor:operlog:remove"), middleware.OperLogMiddleware("Clear Operation Log", constant.REQUEST_BUSINESS_TYPE_DELETE), operlogController.Clean)
			operlogGroup.POST("/export", middleware.HasPerm("monitor:operlog:export"), middleware.OperLogMiddleware("Export Operation Log", constant.REQUEST_BUSINESS_TYPE_EXPORT), operlogController.Export)
		}
	}
}
