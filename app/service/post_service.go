package service

import (
	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"
	"mira/common/types/constant"

	"github.com/pkg/errors"
)

// PostServiceInterface defines operations for post management
type PostServiceInterface interface {
	CreatePost(param dto.SavePost) error
	DeletePost(postIds []int) error
	UpdatePost(param dto.SavePost) error
	GetPostList(param dto.PostListRequest, isPaging bool) ([]dto.PostListResponse, int)
	GetPostByPostId(postId int) dto.PostDetailResponse
	GetPostByPostName(postName string) dto.PostDetailResponse
	GetPostByPostCode(postCode string) dto.PostDetailResponse
	GetPostIdsByUserId(userId int) []int
	GetPostNamesByUserId(userId int) []string
}

// PostService implements the post management interface
type PostService struct{}

// Ensure PostService implements PostServiceInterface
var _ PostServiceInterface = (*PostService)(nil)

// NewPostService creates a new PostService
func NewPostService() *PostService {
	return &PostService{}
}

// CreatePost creates a new post
func (s *PostService) CreatePost(param dto.SavePost) error {
	return s.CreatePostWithErr(param)
}

// CreatePostWithErr creates a new post with proper error handling
func (s *PostService) CreatePostWithErr(param dto.SavePost) error {
	// Input validation
	if param.PostCode == "" {
		return errors.New("post code cannot be empty")
	}

	if param.PostName == "" {
		return errors.New("post name cannot be empty")
	}

	err := dal.Gorm.Model(model.SysPost{}).Create(&model.SysPost{
		PostCode: param.PostCode,
		PostName: param.PostName,
		PostSort: param.PostSort,
		Status:   param.Status,
		Remark:   param.Remark,
		CreateBy: param.CreateBy,
	}).Error
	if err != nil {
		return errors.Wrap(err, "failed to create post")
	}

	return nil
}

// DeletePost deletes posts by IDs
func (s *PostService) DeletePost(postIds []int) error {
	return s.DeletePostWithErr(postIds)
}

// DeletePostWithErr deletes posts by IDs with proper error handling
func (s *PostService) DeletePostWithErr(postIds []int) error {
	// Input validation
	if len(postIds) == 0 {
		return errors.New("post IDs cannot be empty")
	}

	err := dal.Gorm.Model(model.SysPost{}).Where("post_id IN ?", postIds).Delete(&model.SysPost{}).Error
	if err != nil {
		return errors.Wrap(err, "failed to delete posts")
	}

	return nil
}

// UpdatePost updates an existing post
func (s *PostService) UpdatePost(param dto.SavePost) error {
	return s.UpdatePostWithErr(param)
}

// UpdatePostWithErr updates an existing post with proper error handling
func (s *PostService) UpdatePostWithErr(param dto.SavePost) error {
	// Input validation
	if param.PostId <= 0 {
		return errors.New("invalid post ID")
	}

	if param.PostCode == "" {
		return errors.New("post code cannot be empty")
	}

	if param.PostName == "" {
		return errors.New("post name cannot be empty")
	}

	err := dal.Gorm.Model(model.SysPost{}).Where("post_id = ?", param.PostId).Updates(&model.SysPost{
		PostCode: param.PostCode,
		PostName: param.PostName,
		PostSort: param.PostSort,
		Status:   param.Status,
		Remark:   param.Remark,
		UpdateBy: param.UpdateBy,
	}).Error
	if err != nil {
		return errors.Wrap(err, "failed to update post")
	}

	return nil
}

// GetPostList retrieves a list of posts based on search parameters
func (s *PostService) GetPostList(param dto.PostListRequest, isPaging bool) ([]dto.PostListResponse, int) {
	posts, count, _ := s.GetPostListWithErr(param, isPaging)
	return posts, count
}

// GetPostListWithErr retrieves a list of posts with proper error handling
func (s *PostService) GetPostListWithErr(param dto.PostListRequest, isPaging bool) ([]dto.PostListResponse, int, error) {
	var count int64
	posts := make([]dto.PostListResponse, 0)

	query := dal.Gorm.Model(model.SysPost{}).Order("post_sort, post_id")

	if param.PostCode != "" {
		query = query.Where("post_code LIKE ?", "%"+param.PostCode+"%")
	}

	if param.PostName != "" {
		query = query.Where("post_name LIKE ?", "%"+param.PostName+"%")
	}

	if param.Status != "" {
		query = query.Where("status = ?", param.Status)
	}

	if isPaging {
		err := query.Count(&count).Error
		if err != nil {
			return nil, 0, errors.Wrap(err, "failed to count posts")
		}
		query = query.Offset((param.PageNum - 1) * param.PageSize).Limit(param.PageSize)
	}

	err := query.Find(&posts).Error
	if err != nil {
		return nil, 0, errors.Wrap(err, "failed to retrieve posts")
	}

	return posts, int(count), nil
}

// GetPostByPostId retrieves post details by ID
func (s *PostService) GetPostByPostId(postId int) dto.PostDetailResponse {
	post, _ := s.GetPostByPostIdWithErr(postId)
	return post
}

// GetPostByPostIdWithErr retrieves post details by ID with proper error handling
func (s *PostService) GetPostByPostIdWithErr(postId int) (dto.PostDetailResponse, error) {
	var post dto.PostDetailResponse

	if postId <= 0 {
		return post, errors.New("invalid post ID")
	}

	err := dal.Gorm.Model(model.SysPost{}).Where("post_id = ?", postId).Last(&post).Error
	if err != nil {
		return post, errors.Wrap(err, "failed to retrieve post by ID")
	}

	return post, nil
}

// GetPostByPostName retrieves post details by name
func (s *PostService) GetPostByPostName(postName string) dto.PostDetailResponse {
	post, _ := s.GetPostByPostNameWithErr(postName)
	return post
}

// GetPostByPostNameWithErr retrieves post details by name with proper error handling
func (s *PostService) GetPostByPostNameWithErr(postName string) (dto.PostDetailResponse, error) {
	var post dto.PostDetailResponse

	if postName == "" {
		return post, errors.New("post name cannot be empty")
	}

	err := dal.Gorm.Model(model.SysPost{}).Where("post_name = ?", postName).Last(&post).Error
	if err != nil {
		return post, errors.Wrap(err, "failed to retrieve post by name")
	}

	return post, nil
}

// GetPostByPostCode retrieves post details by code
func (s *PostService) GetPostByPostCode(postCode string) dto.PostDetailResponse {
	post, _ := s.GetPostByPostCodeWithErr(postCode)
	return post
}

// GetPostByPostCodeWithErr retrieves post details by code with proper error handling
func (s *PostService) GetPostByPostCodeWithErr(postCode string) (dto.PostDetailResponse, error) {
	var post dto.PostDetailResponse

	if postCode == "" {
		return post, errors.New("post code cannot be empty")
	}

	err := dal.Gorm.Model(model.SysPost{}).Where("post_code = ?", postCode).Last(&post).Error
	if err != nil {
		return post, errors.Wrap(err, "failed to retrieve post by code")
	}

	return post, nil
}

// GetPostIdsByUserId retrieves post IDs assigned to a user
func (s *PostService) GetPostIdsByUserId(userId int) []int {
	postIds, _ := s.GetPostIdsByUserIdWithErr(userId)
	return postIds
}

// GetPostIdsByUserIdWithErr retrieves post IDs assigned to a user with proper error handling
func (s *PostService) GetPostIdsByUserIdWithErr(userId int) ([]int, error) {
	if userId <= 0 {
		return nil, errors.New("invalid user ID")
	}

	var postIds []int

	err := dal.Gorm.Model(model.SysPost{}).
		Joins("JOIN sys_user_post ON sys_user_post.post_id = sys_post.post_id").
		Where("sys_user_post.user_id = ? AND sys_post.status = ?", userId, constant.NORMAL_STATUS).
		Pluck("sys_post.post_id", &postIds).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve post IDs for user")
	}

	return postIds, nil
}

// GetPostNamesByUserId retrieves post names assigned to a user
func (s *PostService) GetPostNamesByUserId(userId int) []string {
	postNames, _ := s.GetPostNamesByUserIdWithErr(userId)
	return postNames
}

// GetPostNamesByUserIdWithErr retrieves post names assigned to a user with proper error handling
func (s *PostService) GetPostNamesByUserIdWithErr(userId int) ([]string, error) {
	if userId <= 0 {
		return nil, errors.New("invalid user ID")
	}

	var postNames []string

	err := dal.Gorm.Model(model.SysPost{}).
		Joins("JOIN sys_user_post ON sys_user_post.post_id = sys_post.post_id").
		Where("sys_user_post.user_id = ? AND sys_post.status = ?", userId, constant.NORMAL_STATUS).
		Pluck("sys_post.post_name", &postNames).Error
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve post names for user")
	}

	return postNames, nil
}
