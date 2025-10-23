package service

import (
	"context"
	"fmt"
	"time"

	"mira/anima/dal"
	"mira/common/types/redis-key"

	"gorm.io/gorm"
)

// OptimizedDataScopeService provides cached and optimized data scope operations
type OptimizedDataScopeService struct {
	*DataScopeService
	cacheService *CacheService
}

// NewOptimizedDataScopeService creates a new optimized data scope service
func NewOptimizedDataScopeService(userService DataScopeUserServiceInterface, roleService DataScopeRoleServiceInterface) *OptimizedDataScopeService {
	return &OptimizedDataScopeService{
		DataScopeService: NewDataScopeService(userService, roleService),
		cacheService:     NewCacheService(),
	}
}

// DataScopeInfo represents cached data scope information
type DataScopeInfo struct {
	UserID      int                    `json:"user_id"`
	RoleIDs     []int                  `json:"role_ids"`
	DeptID      int                    `json:"dept_id"`
	DataScope   string                 `json:"data_scope"`
	DeptIDs     []int                  `json:"dept_ids"`
	UserIDs     []int                  `json:"user_ids"`
	CacheTime   time.Time              `json:"cache_time"`
	TTL         time.Duration          `json:"ttl"`
}

// GetDataScopeOptimized returns optimized data scope filtering with caching
func (ods *OptimizedDataScopeService) GetDataScopeOptimized(ctx context.Context, deptAlias string, userId int, userAlias string) func(*gorm.DB) *gorm.DB {
	// Check cache first
	cacheKey := rediskey.UserDataScopeKey(userId)
	var scopeInfo DataScopeInfo

	err := ods.cacheService.Get(ctx, cacheKey, &scopeInfo)
	if err == nil && time.Since(scopeInfo.CacheTime) < scopeInfo.TTL {
		// Cache hit
		return ods.buildDataScopeQuery(deptAlias, userId, userAlias, &scopeInfo)
	}

	// Cache miss, calculate data scope
	scopeInfo = ods.calculateDataScope(ctx, userId)

	// Cache the result for 30 minutes
	scopeInfo.CacheTime = time.Now()
	scopeInfo.TTL = 30 * time.Minute

	if setErr := ods.cacheService.Set(ctx, cacheKey, scopeInfo, scopeInfo.TTL); setErr != nil {
		fmt.Printf("Warning: failed to cache data scope for user %d: %v\n", userId, setErr)
	}

	return ods.buildDataScopeQuery(deptAlias, userId, userAlias, &scopeInfo)
}

// calculateDataScope calculates data scope information for a user
func (ods *OptimizedDataScopeService) calculateDataScope(ctx context.Context, userId int) DataScopeInfo {
	// Get user information
	userInfo := ods.userService.GetUserByUserId(userId)
	roleList := ods.roleService.GetRoleListByUserIdCompat(userId)

	// Initialize scope info
	scopeInfo := DataScopeInfo{
		UserID:    userId,
		DeptID:    userInfo.DeptId,
		RoleIDs:   make([]int, len(roleList)),
		DataScope: DATA_SCOPE_ALL, // Default to all permissions
	}

	// Extract role IDs
	for i, role := range roleList {
		scopeInfo.RoleIDs[i] = role.RoleId
		if role.DataScope != "" {
			scopeInfo.DataScope = role.DataScope
		}
	}

	// Calculate department and user access based on data scope
	switch scopeInfo.DataScope {
	case DATA_SCOPE_ALL:
		// All data - no restrictions needed
		break

	case DATA_SCOPE_CUSTOM:
		scopeInfo.DeptIDs = ods.getCustomDeptIds(ctx, scopeInfo.RoleIDs)
		scopeInfo.UserIDs = ods.getCustomUserIds(ctx, scopeInfo.DeptIDs)

	case DATA_SCOPE_DEPT:
		scopeInfo.DeptIDs = []int{userInfo.DeptId}
		scopeInfo.UserIDs = ods.getDeptUserIds(userInfo.DeptId)

	case DATA_SCOPE_DEPT_SUB:
		scopeInfo.DeptIDs = ods.getSubDeptIds(userInfo.DeptId)
		scopeInfo.UserIDs = ods.getDeptUserIds(scopeInfo.DeptIDs...)

	case DATA_SCOPE_PERSONAL:
		scopeInfo.UserIDs = []int{userId}
	}

	return scopeInfo
}

// getCustomDeptIds gets department IDs for custom data scope
func (ods *OptimizedDataScopeService) getCustomDeptIds(ctx context.Context, roleIDs []int) []int {
	cacheKey := fmt.Sprintf("custom_depts:%v", roleIDs)
	var deptIDs []int

	// Try cache first
	err := ods.cacheService.Get(ctx, cacheKey, &deptIDs)
	if err == nil {
		return deptIDs
	}

	// Query database for custom department permissions
	var deptList []struct {
		DeptId int
	}

	query := `
		SELECT DISTINCT d.dept_id
		FROM sys_dept d
		INNER JOIN sys_role_dept rd ON d.dept_id = rd.dept_id
		WHERE rd.role_id IN ?
		AND d.status = '0'
	`

	if err := dal.Gorm.Raw(query, roleIDs).Scan(&deptList).Error; err != nil {
		fmt.Printf("Error getting custom department IDs: %v\n", err)
		return []int{}
	}

	deptIDs = make([]int, len(deptList))
	for i, dept := range deptList {
		deptIDs[i] = dept.DeptId
	}

	// Cache for 15 minutes
	ods.cacheService.Set(ctx, cacheKey, deptIDs, 15*time.Minute)

	return deptIDs
}

// getCustomUserIds gets user IDs for custom data scope
func (ods *OptimizedDataScopeService) getCustomUserIds(ctx context.Context, deptIDs []int) []int {
	if len(deptIDs) == 0 {
		return []int{}
	}

	cacheKey := fmt.Sprintf("custom_users:%v", deptIDs)
	var userIDs []int

	// Try cache first
	err := ods.cacheService.Get(ctx, cacheKey, &userIDs)
	if err == nil {
		return userIDs
	}

	// Query database for users in specified departments
	var userList []struct {
		UserId int
	}

	query := `
		SELECT DISTINCT user_id
		FROM sys_user
		WHERE dept_id IN ?
		AND status = '0'
	`

	if err := dal.Gorm.Raw(query, deptIDs).Scan(&userList).Error; err != nil {
		fmt.Printf("Error getting custom user IDs: %v\n", err)
		return []int{}
	}

	userIDs = make([]int, len(userList))
	for i, user := range userList {
		userIDs[i] = user.UserId
	}

	// Cache for 10 minutes
	ods.cacheService.Set(ctx, cacheKey, userIDs, 10*time.Minute)

	return userIDs
}

// getSubDeptIds gets all sub-department IDs for a given department
func (ods *OptimizedDataScopeService) getSubDeptIds(deptId int) []int {
	cacheKey := fmt.Sprintf("sub_depts:%d", deptId)
	var deptIDs []int

	// Try cache first
	err := ods.cacheService.Get(context.Background(), cacheKey, &deptIDs)
	if err == nil {
		return deptIDs
	}

	// Query database for sub-departments recursively
	var deptList []struct {
		DeptId int
	}

	query := `
		WITH RECURSIVE dept_tree AS (
			SELECT dept_id FROM sys_dept WHERE dept_id = ?
			UNION ALL
			SELECT d.dept_id FROM sys_dept d
			INNER JOIN dept_tree dt ON d.parent_id = dt.dept_id
			WHERE d.status = '0'
		)
		SELECT dept_id FROM dept_tree
	`

	if err := dal.Gorm.Raw(query, deptId).Scan(&deptList).Error; err != nil {
		fmt.Printf("Error getting sub-department IDs: %v\n", err)
		return []int{deptId}
	}

	deptIDs = make([]int, len(deptList))
	for i, dept := range deptList {
		deptIDs[i] = dept.DeptId
	}

	// Cache for 1 hour
	ods.cacheService.Set(context.Background(), cacheKey, deptIDs, 1*time.Hour)

	return deptIDs
}

// getDeptUserIds gets user IDs for specified departments
func (ods *OptimizedDataScopeService) getDeptUserIds(deptIDs ...int) []int {
	if len(deptIDs) == 0 {
		return []int{}
	}

	cacheKey := fmt.Sprintf("dept_users:%v", deptIDs)
	var userIDs []int

	// Try cache first
	err := ods.cacheService.Get(context.Background(), cacheKey, &userIDs)
	if err == nil {
		return userIDs
	}

	// Query database for users in specified departments
	var userList []struct {
		UserId int
	}

	query := `
		SELECT user_id
		FROM sys_user
		WHERE dept_id IN ?
		AND status = '0'
	`

	if err := dal.Gorm.Raw(query, deptIDs).Scan(&userList).Error; err != nil {
		fmt.Printf("Error getting department user IDs: %v\n", err)
		return []int{}
	}

	userIDs = make([]int, len(userList))
	for i, user := range userList {
		userIDs[i] = user.UserId
	}

	// Cache for 15 minutes
	ods.cacheService.Set(context.Background(), cacheKey, userIDs, 15*time.Minute)

	return userIDs
}

// buildDataScopeQuery builds the actual GORM query based on scope information
func (ods *OptimizedDataScopeService) buildDataScopeQuery(deptAlias string, userId int, userAlias string, scopeInfo *DataScopeInfo) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		// Skip data scope filtering for super admin
		if userId == SUPER_ADMIN_USER_ID {
			return db
		}

		switch scopeInfo.DataScope {
		case DATA_SCOPE_ALL:
			// No filtering needed
			return db

		case DATA_SCOPE_CUSTOM:
			if len(scopeInfo.DeptIDs) > 0 {
				db = db.Where(fmt.Sprintf("%s.dept_id IN ?", deptAlias), scopeInfo.DeptIDs)
			}
			if len(scopeInfo.UserIDs) > 0 {
				db = db.Where(fmt.Sprintf("%s.user_id IN ?", userAlias), scopeInfo.UserIDs)
			}

		case DATA_SCOPE_DEPT:
			if scopeInfo.DeptID > 0 {
				db = db.Where(fmt.Sprintf("%s.dept_id = ?", deptAlias), scopeInfo.DeptID)
			}

		case DATA_SCOPE_DEPT_SUB:
			if len(scopeInfo.DeptIDs) > 0 {
				db = db.Where(fmt.Sprintf("%s.dept_id IN ?", deptAlias), scopeInfo.DeptIDs)
			}

		case DATA_SCOPE_PERSONAL:
			db = db.Where(fmt.Sprintf("%s.user_id = ?", userAlias), userId)

		default:
			// Default to personal data if scope is unknown
			db = db.Where(fmt.Sprintf("%s.user_id = ?", userAlias), userId)
		}

		return db
	}
}

// InvalidateDataScopeCache invalidates data scope cache for a user
func (ods *OptimizedDataScopeService) InvalidateDataScopeCache(ctx context.Context, userId int) error {
	cacheKey := rediskey.UserDataScopeKey(userId)
	return ods.cacheService.Delete(ctx, cacheKey)
}

// InvalidateDeptCache invalidates department-related caches
func (ods *OptimizedDataScopeService) InvalidateDeptCache(ctx context.Context, deptID int) error {
	patterns := []string{
		fmt.Sprintf("sub_depts:%d", deptID),
		fmt.Sprintf("dept_users:%d", deptID),
		fmt.Sprintf("custom_users:%d", deptID),
	}

	for _, pattern := range patterns {
		if err := ods.cacheService.Delete(ctx, pattern); err != nil {
			return err
		}
	}

	return nil
}

// InvalidateRoleCache invalidates role-related caches
func (ods *OptimizedDataScopeService) InvalidateRoleCache(ctx context.Context, roleID int) error {
	patterns := []string{
		fmt.Sprintf("custom_depts:%d", roleID),
	}

	for _, pattern := range patterns {
		if err := ods.cacheService.Delete(ctx, pattern); err != nil {
			return err
		}
	}

	return nil
}

// BatchInvalidateCache invalidates multiple cache entries
func (ods *OptimizedDataScopeService) BatchInvalidateCache(ctx context.Context, userIDs, deptIDs, roleIDs []int) error {
	var keys []string

	// Add user cache keys
	for _, userID := range userIDs {
		keys = append(keys, rediskey.UserDataScopeKey(userID))
	}

	// Add department cache keys
	for _, deptID := range deptIDs {
		keys = append(keys, fmt.Sprintf("sub_depts:%d", deptID))
		keys = append(keys, fmt.Sprintf("dept_users:%d", deptID))
	}

	// Add role cache keys
	for _, roleID := range roleIDs {
		keys = append(keys, fmt.Sprintf("custom_depts:%d", roleID))
	}

	return ods.cacheService.DeleteMultiple(ctx, keys)
}

// GetCachedDataScopeInfo returns cached data scope information for a user
func (ods *OptimizedDataScopeService) GetCachedDataScopeInfo(ctx context.Context, userID int) (*DataScopeInfo, error) {
	cacheKey := rediskey.UserDataScopeKey(userID)
	var scopeInfo DataScopeInfo

	err := ods.cacheService.Get(ctx, cacheKey, &scopeInfo)
	if err != nil {
		return nil, err
	}

	return &scopeInfo, nil
}

// PreloadDataScopeCache preloads data scope cache for multiple users
func (ods *OptimizedDataScopeService) PreloadDataScopeCache(ctx context.Context, userIDs []int) error {
	if len(userIDs) == 0 {
		return nil
	}

	// Create batch of cache items
	items := make([]CacheItem, 0, len(userIDs))

	for _, userID := range userIDs {
		scopeInfo := ods.calculateDataScope(ctx, userID)
		scopeInfo.CacheTime = time.Now()
		scopeInfo.TTL = 30 * time.Minute

		cacheKey := rediskey.UserDataScopeKey(userID)
		items = append(items, CacheItem{
			Key:        cacheKey,
			Value:      scopeInfo,
			Expiration: scopeInfo.TTL,
		})
	}

	return ods.cacheService.SetMultiple(ctx, items)
}