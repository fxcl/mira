package service

import (
	"testing"

	"mira/anima/dal"
	"mira/app/dto"
	"mira/app/model"

	"github.com/stretchr/testify/assert"
)

func TestPostService_CreatePost(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should create post successfully", func(t *testing.T) {
		// Prepare
		post := dto.SavePost{
			PostCode: "test_code",
			PostName: "Test Post",
			PostSort: 1,
			Status:   "0",
			Remark:   "Test Remark",
			CreateBy: "test_user",
		}

		// Execute
		err := s.CreatePostWithErr(post)
		assert.NoError(t, err)

		// Verify
		var result model.SysPost
		dal.Gorm.First(&result, "post_code = ?", "test_code")
		assert.Equal(t, "Test Post", result.PostName)
	})

	t.Run("should return error when post code is empty", func(t *testing.T) {
		// Prepare
		post := dto.SavePost{
			PostName: "Test Post",
		}

		// Execute
		err := s.CreatePostWithErr(post)
		assert.Error(t, err)
	})

	t.Run("should return error when post name is empty", func(t *testing.T) {
		// Prepare
		post := dto.SavePost{
			PostCode: "test_code",
		}

		// Execute
		err := s.CreatePostWithErr(post)
		assert.Error(t, err)
	})
}

func TestPostService_UpdatePost(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should update post successfully", func(t *testing.T) {
		// Prepare
		post := model.SysPost{PostId: 1, PostCode: "old_code", PostName: "Old Post"}
		dal.Gorm.Create(&post)
		update := dto.SavePost{
			PostId:   1,
			PostCode: "new_code",
			PostName: "New Post",
		}

		// Execute
		err := s.UpdatePostWithErr(update)
		assert.NoError(t, err)

		// Verify
		var result model.SysPost
		dal.Gorm.First(&result, 1)
		assert.Equal(t, "new_code", result.PostCode)
		assert.Equal(t, "New Post", result.PostName)
	})
}

func TestPostService_DeletePost(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should delete post successfully", func(t *testing.T) {
		// Prepare
		post := model.SysPost{PostId: 1, PostCode: "test_code", PostName: "Test Post"}
		dal.Gorm.Create(&post)

		// Execute
		err := s.DeletePostWithErr([]int{1})
		assert.NoError(t, err)

		// Verify
		var result model.SysPost
		err = dal.Gorm.First(&result, 1).Error
		assert.Error(t, err, "record not found")
	})
}

func TestPostService_GetPostList(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should get post list", func(t *testing.T) {
		// Prepare
		post1 := model.SysPost{PostId: 1, PostCode: "code1", PostName: "Post 1"}
		post2 := model.SysPost{PostId: 2, PostCode: "code2", PostName: "Post 2"}
		dal.Gorm.Create(&post1)
		dal.Gorm.Create(&post2)

		// Execute
		params := dto.PostListRequest{
			PageRequest: dto.PageRequest{PageNum: 1, PageSize: 10},
		}
		posts, count, err := s.GetPostListWithErr(params, true)

		// Verify
		assert.NoError(t, err)
		assert.Equal(t, 2, count)
		assert.Len(t, posts, 2)
	})
}

func TestPostService_GetPostByPostId(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should get post by post id", func(t *testing.T) {
		// Prepare
		post := model.SysPost{PostId: 1, PostCode: "test_code", PostName: "Test Post"}
		dal.Gorm.Create(&post)

		// Execute
		result, err := s.GetPostByPostIdWithErr(1)
		assert.NoError(t, err)
		assert.Equal(t, "Test Post", result.PostName)
	})
}

func TestPostService_GetPostByPostName(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should get post by post name", func(t *testing.T) {
		// Prepare
		post := model.SysPost{PostId: 1, PostCode: "test_code", PostName: "Test Post"}
		dal.Gorm.Create(&post)

		// Execute
		result, err := s.GetPostByPostNameWithErr("Test Post")
		assert.NoError(t, err)
		assert.Equal(t, "test_code", result.PostCode)
	})
}

func TestPostService_GetPostByPostCode(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should get post by post code", func(t *testing.T) {
		// Prepare
		post := model.SysPost{PostId: 1, PostCode: "test_code", PostName: "Test Post"}
		dal.Gorm.Create(&post)

		// Execute
		result, err := s.GetPostByPostCodeWithErr("test_code")
		assert.NoError(t, err)
		assert.Equal(t, "Test Post", result.PostName)
	})
}

func TestPostService_GetPostIdsByUserId(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should get post ids by user id", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysPost{PostId: 1, Status: "0"})
		dal.Gorm.Create(&model.SysPost{PostId: 2, Status: "0"})
		dal.Gorm.Create(&model.SysUserPost{UserId: 1, PostId: 1})
		dal.Gorm.Create(&model.SysUserPost{UserId: 1, PostId: 2})

		// Execute
		postIds, err := s.GetPostIdsByUserIdWithErr(1)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []int{1, 2}, postIds)
	})
}

func TestPostService_GetPostNamesByUserId(t *testing.T) {
	setup()
	defer teardown()
	s := NewPostService()

	t.Run("should get post names by user id", func(t *testing.T) {
		// Prepare
		dal.Gorm.Create(&model.SysUser{UserId: 1})
		dal.Gorm.Create(&model.SysPost{PostId: 1, PostName: "Post 1", Status: "0"})
		dal.Gorm.Create(&model.SysPost{PostId: 2, PostName: "Post 2", Status: "0"})
		dal.Gorm.Create(&model.SysUserPost{UserId: 1, PostId: 1})
		dal.Gorm.Create(&model.SysUserPost{UserId: 1, PostId: 2})

		// Execute
		postNames, err := s.GetPostNamesByUserIdWithErr(1)
		assert.NoError(t, err)
		assert.ElementsMatch(t, []string{"Post 1", "Post 2"}, postNames)
	})
}
