package router

import (
	"mira/app/controller"
	monitorcontroller "mira/app/controller/monitor"
	systemcontroller "mira/app/controller/system"
	"mira/app/middleware"
	"mira/common/types/constant"

	"github.com/gin-gonic/gin"
)

// Admin router group
func RegisterAdminGroupApi(api *gin.RouterGroup) {
	api.Use(middleware.Cors()) // CORS middleware

	registerAuthRoutes(api)
	registerUserRoutes(api)
	registerRoleRoutes(api)
	registerMenuRoutes(api)
	registerDeptRoutes(api)
	registerPostRoutes(api)
	registerDictRoutes(api)
	registerConfigRoutes(api)
	registerMonitorRoutes(api)
}

func registerAuthRoutes(api *gin.RouterGroup) {
	api.GET("/captchaImage", (&controller.AuthController{}).CaptchaImage)                       // Get captcha image
	api.POST("/register", (&controller.AuthController{}).Register)                              // Register
	api.POST("/login", middleware.LogininforMiddleware(), (&controller.AuthController{}).Login) // Login
	api.POST("/logout", (&controller.AuthController{}).Logout)                                  // Logout

	// Enable authentication middleware. The following routes require authentication.
	api.Use(middleware.AuthMiddleware())

	api.GET("/getInfo", (&controller.AuthController{}).GetInfo)       // Get user info
	api.GET("/getRouters", (&controller.AuthController{}).GetRouters) // Get routing information
}

func registerUserRoutes(api *gin.RouterGroup) {
	api.GET("/system/user/profile", (&systemcontroller.UserController{}).GetProfile)                      // Personal information
	api.PUT("/system/user/profile", (&systemcontroller.UserController{}).UpdateProfile)                   // Update user
	api.PUT("/system/user/profile/updatePwd", (&systemcontroller.UserController{}).UserProfileUpdatePwd)  // Reset password
	api.POST("/system/user/profile/avatar", (&systemcontroller.UserController{}).UserProfileUpdateAvatar) // Update avatar

	api.GET("/system/user/deptTree", middleware.HasPerm("system:user:list"), (&systemcontroller.UserController{}).DeptTree)          // Get department tree list
	api.GET("/system/user/list", middleware.HasPerm("system:user:list"), (&systemcontroller.UserController{}).List)                  // Get user list
	api.GET("/system/user/", middleware.HasPerm("system:user:query"), (&systemcontroller.UserController{}).Detail)                   // Get user details by user ID
	api.GET("/system/user/:userId", middleware.HasPerm("system:user:query"), (&systemcontroller.UserController{}).Detail)            // Get user details by user ID
	api.GET("/system/user/authRole/:userId", middleware.HasPerm("system:user:query"), (&systemcontroller.UserController{}).AuthRole) // Get user details by user ID

	api.POST("/system/user", middleware.HasPerm("system:user:add"), middleware.OperLogMiddleware("Add User", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.UserController{}).Create)
	api.PUT("/system/user", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("Update User", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.UserController{}).Update)
	api.DELETE("/system/user/:userIds", middleware.HasPerm("system:user:remove"), middleware.OperLogMiddleware("Delete User", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.UserController{}).Remove)
	api.PUT("/system/user/changeStatus", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("Modify User Status", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.UserController{}).ChangeStatus)
	api.PUT("/system/user/resetPwd", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("Modify User Password", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.UserController{}).ResetPwd)
	api.PUT("/system/user/authRole", middleware.HasPerm("system:user:edit"), middleware.OperLogMiddleware("User Authorized Role", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.UserController{}).AddAuthRole)
	api.POST("/system/user/export", middleware.HasPerm("system:user:export"), middleware.OperLogMiddleware("Export User", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&systemcontroller.UserController{}).Export)
	api.POST("/system/user/importData", middleware.HasPerm("system:user:import"), middleware.OperLogMiddleware("Import User", constant.REQUEST_BUSINESS_TYPE_IMPORT), (&systemcontroller.UserController{}).ImportData)
	api.POST("/system/user/importTemplate", middleware.OperLogMiddleware("Import User Template", constant.REQUEST_BUSINESS_TYPE_IMPORT), (&systemcontroller.UserController{}).ImportTemplate)
}

func registerRoleRoutes(api *gin.RouterGroup) {
	api.GET("/system/role/list", middleware.HasPerm("system:role:list"), (&systemcontroller.RoleController{}).List)                                            // Get role list
	api.GET("/system/role/:roleId", middleware.HasPerm("system:role:query"), (&systemcontroller.RoleController{}).Detail)                                      // Get role details
	api.GET("/system/role/deptTree/:roleId", middleware.HasPerm("system:role:query"), (&systemcontroller.RoleController{}).DeptTree)                           // Get department tree
	api.GET("/system/role/authUser/allocatedList", middleware.HasPerm("system:role:list"), (&systemcontroller.RoleController{}).RoleAuthUserAllocatedList)     // Query allocated user role list
	api.GET("/system/role/authUser/unallocatedList", middleware.HasPerm("system:role:list"), (&systemcontroller.RoleController{}).RoleAuthUserUnallocatedList) // Query unallocated user role list

	api.POST("/system/role", middleware.HasPerm("system:role:add"), middleware.OperLogMiddleware("Add Role", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.RoleController{}).Create)
	api.PUT("/system/role", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Update Role", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.RoleController{}).Update)
	api.DELETE("/system/role/:roleIds", middleware.HasPerm("system:role:remove"), middleware.OperLogMiddleware("Delete Role", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.RoleController{}).Remove)
	api.PUT("/system/role/changeStatus", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Modify Role Status", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.RoleController{}).ChangeStatus)
	api.PUT("/system/role/dataScope", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Assign Data Permissions", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.RoleController{}).DataScope)
	api.PUT("/system/role/authUser/selectAll", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Batch Select User Authorization", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.RoleController{}).RoleAuthUserSelectAll)
	api.PUT("/system/role/authUser/cancel", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Cancel Authorized User", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.RoleController{}).RoleAuthUserCancel)
	api.PUT("/system/role/authUser/cancelAll", middleware.HasPerm("system:role:edit"), middleware.OperLogMiddleware("Batch Cancel Authorized User", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.RoleController{}).RoleAuthUserCancelAll)
	api.POST("/system/role/export", middleware.HasPerm("system:role:export"), middleware.OperLogMiddleware("Export Role", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&systemcontroller.RoleController{}).Export)
}

func registerMenuRoutes(api *gin.RouterGroup) {
	api.GET("/system/menu/list", middleware.HasPerm("system:menu:list"), (&systemcontroller.MenuController{}).List)       // Get menu list
	api.GET("/system/menu/treeselect", (&systemcontroller.MenuController{}).Treeselect)                                   // Get menu dropdown tree list
	api.GET("/system/menu/roleMenuTreeselect/:roleId", (&systemcontroller.MenuController{}).RoleMenuTreeselect)           // Load corresponding role menu list tree
	api.GET("/system/menu/:menuId", middleware.HasPerm("system:menu:query"), (&systemcontroller.MenuController{}).Detail) // Get menu details

	api.POST("/system/menu", middleware.HasPerm("system:menu:add"), middleware.OperLogMiddleware("Add Menu", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.MenuController{}).Create)
	api.PUT("/system/menu", middleware.HasPerm("system:menu:edit"), middleware.OperLogMiddleware("Update Menu", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.MenuController{}).Update)
	api.DELETE("/system/menu/:menuId", middleware.HasPerm("system:menu:remove"), middleware.OperLogMiddleware("Delete Menu", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.MenuController{}).Remove)
}

func registerDeptRoutes(api *gin.RouterGroup) {
	api.GET("/system/dept/list", middleware.HasPerm("system:dept:list"), (&systemcontroller.DeptController{}).List)                        // Get department list
	api.GET("/system/dept/list/exclude/:deptId", middleware.HasPerm("system:dept:list"), (&systemcontroller.DeptController{}).ListExclude) // Query department list (exclude nodes)
	api.GET("/system/dept/:deptId", middleware.HasPerm("system:dept:query"), (&systemcontroller.DeptController{}).Detail)                  // Get department details

	api.POST("/system/dept", middleware.HasPerm("system:dept:add"), middleware.OperLogMiddleware("Add Department", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.DeptController{}).Create)
	api.PUT("/system/dept", middleware.HasPerm("system:dept:edit"), middleware.OperLogMiddleware("Update Department", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.DeptController{}).Update)
	api.DELETE("/system/dept/:deptId", middleware.HasPerm("system:dept:remove"), middleware.OperLogMiddleware("Delete Department", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.DeptController{}).Remove)
}

func registerPostRoutes(api *gin.RouterGroup) {
	api.GET("/system/post/list", middleware.HasPerm("system:post:list"), (&systemcontroller.PostController{}).List)       // Get post list
	api.GET("/system/post/:postId", middleware.HasPerm("system:post:query"), (&systemcontroller.PostController{}).Detail) // Get post details

	api.POST("/system/post", middleware.HasPerm("system:post:add"), middleware.OperLogMiddleware("Add Post", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.PostController{}).Create)
	api.PUT("/system/post", middleware.HasPerm("system:post:edit"), middleware.OperLogMiddleware("Update Post", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.PostController{}).Update)
	api.DELETE("/system/post/:postIds", middleware.HasPerm("system:post:remove"), middleware.OperLogMiddleware("Delete Post", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.PostController{}).Remove)
	api.POST("/system/post/export", middleware.HasPerm("system:post:export"), middleware.OperLogMiddleware("Export Post", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&systemcontroller.PostController{}).Export)
}

func registerDictRoutes(api *gin.RouterGroup) {
	api.GET("/system/dict/list", middleware.HasPerm("system:dict:list"), (&systemcontroller.DictTypeController{}).List)            // Get dictionary type list
	api.GET("/system/dict/type/:dictId", middleware.HasPerm("system:dict:query"), (&systemcontroller.DictTypeController{}).Detail) // Get dictionary type details
	api.GET("/system/dict/type/optionselect", (&systemcontroller.DictTypeController{}).Optionselect)                               // Get dictionary select box list

	api.POST("/system/dict/type", middleware.HasPerm("system:dict:add"), middleware.OperLogMiddleware("Add Dictionary Type", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.DictTypeController{}).Create)
	api.PUT("/system/dict/type", middleware.HasPerm("system:dict:edit"), middleware.OperLogMiddleware("Update Dictionary Type", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.DictTypeController{}).Update)
	api.DELETE("/system/dict/type/:dictIds", middleware.HasPerm("system:dict:remove"), middleware.OperLogMiddleware("Delete Dictionary Type", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.DictTypeController{}).Remove)
	api.POST("/system/dict/type/export", middleware.HasPerm("system:dict:export"), middleware.OperLogMiddleware("Export Dictionary Type", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&systemcontroller.DictTypeController{}).Export)
	api.DELETE("/system/dict/type/refreshCache", middleware.HasPerm("system:dict:remove"), middleware.OperLogMiddleware("Refresh Dictionary Type Cache", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.DictTypeController{}).RefreshCache)

	api.GET("/system/dict/data/list", middleware.HasPerm("system:dict:list"), (&systemcontroller.DictDataController{}).List)         // Get dictionary data list
	api.GET("/system/dict/data/:dictCode", middleware.HasPerm("system:dict:query"), (&systemcontroller.DictDataController{}).Detail) // Get dictionary data details
	api.GET("/system/dict/data/type/:dictType", (&systemcontroller.DictDataController{}).Type)                                       // Query dictionary data by dictionary type

	api.POST("/system/dict/data", middleware.HasPerm("system:dict:add"), middleware.OperLogMiddleware("Add Dictionary Data", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.DictDataController{}).Create)
	api.PUT("/system/dict/data", middleware.HasPerm("system:dict:edit"), middleware.OperLogMiddleware("Update Dictionary Data", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.DictDataController{}).Update)
	api.DELETE("/system/dict/data/:dictCodes", middleware.HasPerm("system:dict:remove"), middleware.OperLogMiddleware("Delete Dictionary Data", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.DictDataController{}).Remove)
	api.POST("/system/dict/data/export", middleware.HasPerm("system:dict:export"), middleware.OperLogMiddleware("Export Dictionary Data", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&systemcontroller.DictDataController{}).Export)
}

func registerConfigRoutes(api *gin.RouterGroup) {
	api.GET("/system/config/list", middleware.HasPerm("system:config:list"), (&systemcontroller.ConfigController{}).List)         // Get parameter configuration list
	api.GET("/system/config/:configId", middleware.HasPerm("system:config:query"), (&systemcontroller.ConfigController{}).Detail) // Get parameter configuration details
	api.GET("/system/config/configKey/:configKey", (&systemcontroller.ConfigController{}).ConfigKey)                              // Query parameter value by parameter key name

	api.POST("/system/config", middleware.HasPerm("system:config:add"), middleware.OperLogMiddleware("Add Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_INSERT), (&systemcontroller.ConfigController{}).Create)
	api.PUT("/system/config", middleware.HasPerm("system:config:edit"), middleware.OperLogMiddleware("Update Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_UPDATE), (&systemcontroller.ConfigController{}).Update)
	api.DELETE("/system/config/:configIds", middleware.HasPerm("system:config:remove"), middleware.OperLogMiddleware("Delete Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.ConfigController{}).Remove)
	api.POST("/system/config/export", middleware.HasPerm("system:config:export"), middleware.OperLogMiddleware("Export Parameter Configuration", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&systemcontroller.ConfigController{}).Export)
	api.DELETE("/system/config/refreshCache", middleware.HasPerm("system:config:remove"), middleware.OperLogMiddleware("Refresh Parameter Configuration Cache", constant.REQUEST_BUSINESS_TYPE_DELETE), (&systemcontroller.ConfigController{}).RefreshCache)
}

func registerMonitorRoutes(api *gin.RouterGroup) {
	api.GET("/monitor/logininfor/list", middleware.HasPerm("monitor:operlog:list"), (&monitorcontroller.LogininforController{}).List) // Get login log list

	api.DELETE("/monitor/logininfor/:infoIds", middleware.HasPerm("monitor:logininfor:remove"), middleware.OperLogMiddleware("Delete Login Log", constant.REQUEST_BUSINESS_TYPE_DELETE), (&monitorcontroller.LogininforController{}).Remove)
	api.DELETE("/monitor/logininfor/clean", middleware.HasPerm("monitor:logininfor:remove"), middleware.OperLogMiddleware("Clear Login Log", constant.REQUEST_BUSINESS_TYPE_DELETE), (&monitorcontroller.LogininforController{}).Clean)
	api.GET("/monitor/logininfor/unlock/:userName", middleware.HasPerm("monitor:logininfor:unlock"), middleware.OperLogMiddleware("Account Unlock", constant.REQUEST_BUSINESS_TYPE_DELETE), (&monitorcontroller.LogininforController{}).Unlock)
	api.POST("/monitor/logininfor/export", middleware.HasPerm("monitor:logininfor:export"), middleware.OperLogMiddleware("Export Login Log", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&monitorcontroller.LogininforController{}).Export)

	api.GET("/monitor/operlog/list", middleware.HasPerm("monitor:logininfor:list"), (&monitorcontroller.OperlogController{}).List) // Get operation log list

	api.DELETE("/monitor/operlog/:operIds", middleware.HasPerm("monitor:operlog:remove"), middleware.OperLogMiddleware("Delete Operation Log", constant.REQUEST_BUSINESS_TYPE_DELETE), (&monitorcontroller.OperlogController{}).Remove)
	api.DELETE("/monitor/operlog/clean", middleware.HasPerm("monitor:operlog:remove"), middleware.OperLogMiddleware("Clear Operation Log", constant.REQUEST_BUSINESS_TYPE_DELETE), (&monitorcontroller.OperlogController{}).Clean)
	api.POST("/monitor/operlog/export", middleware.HasPerm("monitor:operlog:export"), middleware.OperLogMiddleware("Export Operation Log", constant.REQUEST_BUSINESS_TYPE_EXPORT), (&monitorcontroller.OperlogController{}).Export)
}
