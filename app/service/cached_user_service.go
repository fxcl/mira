package service

import (
	"context"
	"fmt"
	"time"

	"mira/app/dto"
	"mira/common/types/redis-key"
)

// CachedUserService extends UserService with caching capabilities
type CachedUserService struct {
	*UserService
	cacheService *CacheService
}

// NewCachedUserService creates a new cached user service
func NewCachedUserService() *CachedUserService {
	return &CachedUserService{
		UserService:   &UserService{},
		cacheService:  NewCacheService(),
	}
}

// GetUserByUserId retrieves user by ID with caching
func (s *CachedUserService) GetUserByUserIdCached(ctx context.Context, userId int) (dto.UserDetailResponse, error) {
	var result dto.UserDetailResponse
	cacheKey := rediskey.UserProfileKey(userId)

	// Try cache first
	err := s.cacheService.Get(ctx, cacheKey, &result)
	if err == nil {
		return result, nil // Cache hit
	}

	// Cache miss, get from database
	dbResult := s.UserService.GetUserByUserId(userId)

	// Cache the result for 30 minutes
	if setErr := s.cacheService.Set(ctx, cacheKey, dbResult, 30*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user profile for user %d: %v\n", userId, setErr)
	}

	return dbResult, nil
}

// GetUserByUsername retrieves user by username with caching
func (s *CachedUserService) GetUserByUsernameCached(ctx context.Context, userName string) (dto.UserTokenResponse, error) {
	var result dto.UserTokenResponse
	cacheKey := rediskey.UserTokenKey() + userName

	// Try cache first
	err := s.cacheService.Get(ctx, cacheKey, &result)
	if err == nil {
		return result, nil // Cache hit
	}

	// Cache miss, get from database
	dbResult := s.UserService.GetUserByUsername(userName)

	// Cache the result for 15 minutes (shorter for auth tokens)
	if setErr := s.cacheService.Set(ctx, cacheKey, dbResult, 15*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user token for username %s: %v\n", userName, setErr)
	}

	return dbResult, nil
}

// GetUserByEmail retrieves user by email with caching
func (s *CachedUserService) GetUserByEmailCached(ctx context.Context, email string) (dto.UserTokenResponse, error) {
	var result dto.UserTokenResponse
	cacheKey := rediskey.UserTokenKey() + "email:" + email

	// Try cache first
	err := s.cacheService.Get(ctx, cacheKey, &result)
	if err == nil {
		return result, nil // Cache hit
	}

	// Cache miss, get from database
	dbResult := s.UserService.GetUserByEmail(email)

	// Cache the result for 15 minutes
	if setErr := s.cacheService.Set(ctx, cacheKey, dbResult, 15*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user token for email %s: %v\n", email, setErr)
	}

	return dbResult, nil
}

// GetUserByPhonenumber retrieves user by phone number with caching
func (s *CachedUserService) GetUserByPhonenumberCached(ctx context.Context, phonenumber string) (dto.UserTokenResponse, error) {
	var result dto.UserTokenResponse
	cacheKey := rediskey.UserTokenKey() + "phone:" + phonenumber

	// Try cache first
	err := s.cacheService.Get(ctx, cacheKey, &result)
	if err == nil {
		return result, nil // Cache hit
	}

	// Cache miss, get from database
	dbResult := s.UserService.GetUserByPhonenumber(phonenumber)

	// Cache the result for 15 minutes
	if setErr := s.cacheService.Set(ctx, cacheKey, dbResult, 15*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user token for phone %s: %v\n", phonenumber, setErr)
	}

	return dbResult, nil
}

// GetUserList retrieves user list with caching for pagination
func (s *CachedUserService) GetUserListCached(ctx context.Context, param dto.UserListRequest, userId int, isPaging bool) ([]dto.UserListResponse, int, error) {
	// Create cache key based on request parameters
	cacheKey := fmt.Sprintf("%s:userlist:%d:%v:%d:%d:%s",
		rediskey.UserTokenKey(), userId, isPaging, param.PageNum, param.PageSize, param.UserName)

	var result []dto.UserListResponse
	var total int

	// Try cache first
	err := s.cacheService.Get(ctx, cacheKey, &result)
	if err == nil {
		// Get total count separately
		totalKey := cacheKey + ":total"
		s.cacheService.Get(ctx, totalKey, &total)
		return result, total, nil // Cache hit
	}

	// Cache miss, get from database
	dbResult, dbTotal := s.UserService.GetUserList(param, userId, isPaging)

	// Cache the result for 5 minutes (shorter for dynamic lists)
	if setErr := s.cacheService.Set(ctx, cacheKey, dbResult, 5*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user list: %v\n", setErr)
	}

	// Cache total count
	totalKey := cacheKey + ":total"
	if setErr := s.cacheService.Set(ctx, totalKey, dbTotal, 5*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user list total: %v\n", setErr)
	}

	return dbResult, dbTotal, nil
}

// UserHasPermsCached checks if user has permissions with caching
func (s *CachedUserService) UserHasPermsCached(ctx context.Context, userId int, perms []string) bool {
	cacheKey := rediskey.UserPermsKey(userId)

	// Try cache first
	var cachedPerms []string
	err := s.cacheService.Get(ctx, cacheKey, &cachedPerms)
	if err == nil {
		// Check permissions against cached list
		return s.checkPermissions(cachedPerms, perms)
	}

	// Cache miss, get from database and cache the result
	hasPerms := s.UserService.UserHasPerms(userId, perms)

	// Cache the user's full permission list for 30 minutes
	if setErr := s.cacheService.Set(ctx, cacheKey, perms, 30*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user permissions for user %d: %v\n", userId, setErr)
	}

	return hasPerms
}

// UserHasRolesCached checks if user has roles with caching
func (s *CachedUserService) UserHasRolesCached(ctx context.Context, userId int, roles []string) bool {
	cacheKey := rediskey.UserRolesKey(userId)

	// Try cache first
	var cachedRoles []string
	err := s.cacheService.Get(ctx, cacheKey, &cachedRoles)
	if err == nil {
		// Check roles against cached list
		return s.checkRoles(cachedRoles, roles)
	}

	// Cache miss, get from database and cache the result
	hasRoles := s.UserService.UserHasRoles(userId, roles)

	// Cache the user's role list for 30 minutes
	if setErr := s.cacheService.Set(ctx, cacheKey, roles, 30*time.Minute); setErr != nil {
		fmt.Printf("Warning: failed to cache user roles for user %d: %v\n", userId, setErr)
	}

	return hasRoles
}

// InvalidateUserCache removes all cached data for a user
func (s *CachedUserService) InvalidateUserCache(ctx context.Context, userId int) error {
	return s.cacheService.InvalidateUserCache(ctx, userId)
}

// InvalidateUserListCache removes cached user list data
func (s *CachedUserService) InvalidateUserListCache(ctx context.Context) error {
	return s.cacheService.InvalidatePattern(ctx, rediskey.UserPattern()+":userlist:*")
}

// CreateUser creates a user and invalidates relevant caches
func (s *CachedUserService) CreateUserWithCache(ctx context.Context, param dto.SaveUser, roleIds, postIds []int) error {
	err := s.UserService.CreateUser(param, roleIds, postIds)
	if err != nil {
		return err
	}

	// Invalidate user list cache
	return s.InvalidateUserListCache(ctx)
}

// UpdateUser updates a user and invalidates relevant caches
func (s *CachedUserService) UpdateUserWithCache(ctx context.Context, param dto.SaveUser, roleIds, postIds []int) error {
	err := s.UserService.UpdateUser(param, roleIds, postIds)
	if err != nil {
		return err
	}

	// Invalidate all caches for this user
	return s.InvalidateUserCache(ctx, param.UserId)
}

// DeleteUser deletes users and invalidates relevant caches
func (s *CachedUserService) DeleteUserWithCache(ctx context.Context, userIds []int) error {
	err := s.UserService.DeleteUser(userIds)
	if err != nil {
		return err
	}

	// Invalidate caches for all deleted users
	for _, userId := range userIds {
		if invalidateErr := s.InvalidateUserCache(ctx, userId); invalidateErr != nil {
			fmt.Printf("Warning: failed to invalidate cache for user %d: %v\n", userId, invalidateErr)
		}
	}

	// Invalidate user list cache
	return s.InvalidateUserListCache(ctx)
}

// Helper methods
func (s *CachedUserService) checkPermissions(cachedPerms, requiredPerms []string) bool {
	permMap := make(map[string]bool)
	for _, perm := range cachedPerms {
		permMap[perm] = true
	}

	for _, requiredPerm := range requiredPerms {
		if !permMap[requiredPerm] {
			return false
		}
	}

	return true
}

func (s *CachedUserService) checkRoles(cachedRoles, requiredRoles []string) bool {
	roleMap := make(map[string]bool)
	for _, role := range cachedRoles {
		roleMap[role] = true
	}

	for _, requiredRole := range requiredRoles {
		if !roleMap[requiredRole] {
			return false
		}
	}

	return true
}