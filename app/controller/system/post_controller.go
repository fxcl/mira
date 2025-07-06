package systemcontroller

import (
	"strconv"
	"time"

	"mira/anima/response"
	"mira/app/dto"
	"mira/app/security"
	"mira/app/service"
	"mira/app/validator"
	"mira/common/utils"

	"gitee.com/hanshuangjianke/go-excel/excel"
	"github.com/gin-gonic/gin"
)

// PostController handles post-related operations.
type PostController struct {
	PostService *service.PostService
}

// NewPostController creates a new PostController.
func NewPostController(postService *service.PostService) *PostController {
	return &PostController{PostService: postService}
}

// List retrieves a paginated list of posts.
// @Summary Get post list
// @Description Retrieves a paginated list of posts based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.PostListRequest true "Query parameters"
// @Success 200 {object} response.Response{data=response.PageData{list=[]dto.PostListResponse}} "Success"
// @Router /system/post/list [get]
func (c *PostController) List(ctx *gin.Context) {
	var param dto.PostListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	posts, total := c.PostService.GetPostList(param, true)

	response.NewSuccess().SetPageData(posts, total).Json(ctx)
}

// Detail retrieves the details of a specific post.
// @Summary Get post details
// @Description Retrieves the details of a post by its ID.
// @Tags System
// @Accept json
// @Produce json
// @Param postId path int true "Post ID"
// @Success 200 {object} response.Response{data=dto.PostDetailResponse} "Success"
// @Router /system/post/{postId} [get]
func (c *PostController) Detail(ctx *gin.Context) {
	postId, _ := strconv.Atoi(ctx.Param("postId"))

	post := c.PostService.GetPostByPostId(postId)

	response.NewSuccess().SetData("data", post).Json(ctx)
}

// Create adds a new post.
// @Summary Add post
// @Description Adds a new post to the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.CreatePostRequest true "Post data"
// @Success 200 {object} response.Response "Success"
// @Router /system/post [post]
func (c *PostController) Create(ctx *gin.Context) {
	var param dto.CreatePostRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.CreatePostValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if post := c.PostService.GetPostByPostName(param.PostName); post.PostId > 0 {
		response.NewError().SetMsg("Failed to add post " + param.PostName + ", post name already exists").Json(ctx)
		return
	}

	if post := c.PostService.GetPostByPostCode(param.PostCode); post.PostId > 0 {
		response.NewError().SetMsg("Failed to add post " + param.PostName + ", post code already exists").Json(ctx)
		return
	}

	if err := c.PostService.CreatePost(dto.SavePost{
		PostCode: param.PostCode,
		PostName: param.PostName,
		PostSort: param.PostSort,
		Status:   param.Status,
		CreateBy: security.GetAuthUserName(ctx),
		Remark:   param.Remark,
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Update modifies an existing post.
// @Summary Update post
// @Description Modifies an existing post in the system.
// @Tags System
// @Accept json
// @Produce json
// @Param body body dto.UpdatePostRequest true "Post data"
// @Success 200 {object} response.Response "Success"
// @Router /system/post [put]
func (c *PostController) Update(ctx *gin.Context) {
	var param dto.UpdatePostRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err := validator.UpdatePostValidator(param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if post := c.PostService.GetPostByPostName(param.PostName); post.PostId > 0 && post.PostId != param.PostId {
		response.NewError().SetMsg("Failed to modify post " + param.PostName + ", post name already exists").Json(ctx)
		return
	}

	if post := c.PostService.GetPostByPostCode(param.PostCode); post.PostId > 0 && post.PostId != param.PostId {
		response.NewError().SetMsg("Failed to modify post " + param.PostName + ", post code already exists").Json(ctx)
		return
	}

	if err := c.PostService.UpdatePost(dto.SavePost{
		PostId:   param.PostId,
		PostCode: param.PostCode,
		PostName: param.PostName,
		PostSort: param.PostSort,
		Status:   param.Status,
		UpdateBy: security.GetAuthUserName(ctx),
		Remark:   param.Remark,
	}); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	response.NewSuccess().Json(ctx)
}

// Remove deletes one or more posts.
// @Summary Delete post
// @Description Deletes posts by their IDs.
// @Tags System
// @Accept json
// @Produce json
// @Param postIds path string true "Post IDs, comma-separated"
// @Success 200 {object} response.Response "Success"
// @Router /system/post/{postIds} [delete]
func (c *PostController) Remove(ctx *gin.Context) {
	postIds, err := utils.StringToIntSlice(ctx.Param("postIds"), ",")
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	if err = c.PostService.DeletePost(postIds); err != nil {
		response.NewError().SetMsg(err.Error())
		return
	}

	response.NewSuccess().Json(ctx)
}

// Export exports post data to an Excel file.
// @Summary Export posts
// @Description Exports post data to an Excel file based on query parameters.
// @Tags System
// @Accept json
// @Produce json
// @Param query body dto.PostListRequest true "Query parameters"
// @Success 200 {file} file "Excel file"
// @Router /system/post/export [post]
func (c *PostController) Export(ctx *gin.Context) {
	var param dto.PostListRequest

	if err := ctx.ShouldBind(&param); err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	list := make([]dto.PostExportResponse, 0)

	posts, _ := c.PostService.GetPostList(param, false)
	for _, post := range posts {
		list = append(list, dto.PostExportResponse{
			PostId:   post.PostId,
			PostCode: post.PostCode,
			PostName: post.PostName,
			PostSort: post.PostSort,
			Status:   post.Status,
		})
	}

	file, err := excel.NormalDynamicExport("Sheet1", "", "", false, false, list, nil)
	if err != nil {
		response.NewError().SetMsg(err.Error()).Json(ctx)
		return
	}

	excel.DownLoadExcel("post_"+time.Now().Format("20060102150405"), ctx.Writer, file)
}
