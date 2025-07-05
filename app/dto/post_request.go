package dto

// Save Post
type SavePost struct {
	PostId   int    `json:"postId"`
	PostCode string `json:"postCode"`
	PostName string `json:"postName"`
	PostSort int    `json:"postSort"`
	Status   string `json:"status"`
	CreateBy string `json:"createBy"`
	UpdateBy string `json:"updateBy"`
	Remark   string `json:"remark"`
}

// Post List
type PostListRequest struct {
	PageRequest
	PostCode string `query:"postCode" form:"postCode"`
	PostName string `query:"postName" form:"postName"`
	Status   string `query:"status" form:"status"`
}

// Create Post
type CreatePostRequest struct {
	PostCode string `json:"postCode"`
	PostName string `json:"postName"`
	PostSort int    `json:"postSort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// Update Post
type UpdatePostRequest struct {
	PostId   int    `json:"postId"`
	PostCode string `json:"postCode"`
	PostName string `json:"postName"`
	PostSort int    `json:"postSort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}
