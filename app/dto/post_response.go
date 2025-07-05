package dto

import "mira/anima/datetime"

// Post List
type PostListResponse struct {
	PostId     int               `json:"postId"`
	PostCode   string            `json:"postCode"`
	PostName   string            `json:"postName"`
	PostSort   int               `json:"postSort"`
	Status     string            `json:"status"`
	CreateTime datetime.Datetime `json:"createTime"`
}

// Post Details
type PostDetailResponse struct {
	PostId   int    `json:"postId"`
	PostCode string `json:"postCode"`
	PostName string `json:"postName"`
	PostSort int    `json:"postSort"`
	Status   string `json:"status"`
	Remark   string `json:"remark"`
}

// Post Export
type PostExportResponse struct {
	PostId   int    `excel:"name:Post ID;"`
	PostCode string `excel:"name:Post Code;"`
	PostName string `excel:"name:Post Name;"`
	PostSort int    `excel:"name:Post Sort;"`
	Status   string `excel:"name:Status;replace:0_Normal,1_Disabled;"`
}
