package request

import (
	"mime/multipart"
	"time"
)

type TaskGetListRequest struct {
	Id              int64      `form:"id" json:"id,omitempty"`
	UserID          int64      `form:"user_id" json:"user_id,omitempty"`
	Title           string     `form:"title" json:"title,omitempty" binding:"max=255"`
	SpecifyDatetime *time.Time `form:"specify_datetime" json:"specify_datetime,omitempty" time_format:"2006-01-02 15:04:05"`
	IsSpecifyTime   *bool      `form:"is_specify_time" json:"is_specify_time,omitempty"`
	IsComplete      *bool      `form:"is_complete" json:"is_complete,omitempty"`
	Pagination
}

type TaskCreateRequest struct {
	UserID          int64                 `form:"user_id" json:"user_id" binding:"required"`
	CategoryID      int64                 `form:"category_id" json:"category_id" binding:"required"`
	Title           string                `form:"title" json:"title" binding:"required,max=100"`
	Note            string                `form:"note" json:"note,omitempty"`
	Url             string                `form:"url" json:"url,omitempty"`
	Image           *multipart.FileHeader `form:"image" json:"image,omitempty"`
	SpecifyDatetime *time.Time            `form:"specify_datetime" json:"specify_datetime,omitempty" time_format:"2006-01-02 15:04:05"`
	IsSpecifyTime   bool                  `form:"is_specify_time" json:"is_specify_time,omitempty"`
	Priority        int8                  `form:"priority" json:"priority" binding:"required,oneof=1 2 3"`
	IsComplete      bool                  `form:"is_complete" json:"is_complete,omitempty"`
}

type TaskUpdateRequest struct {
	CategoryID      int64                 `form:"category_id" json:"category_id,omitempty"`
	Title           string                `form:"title" json:"title,omitempty" binding:"max=100"`
	Note            string                `form:"note" json:"note,omitempty"`
	Url             string                `form:"url" json:"url,omitempty"`
	Image           *multipart.FileHeader `form:"image" json:"image,omitempty"`
	SpecifyDatetime *time.Time            `form:"specify_datetime" json:"specify_datetime,omitempty" time_format:"2006-01-02 15:04:05"`
	IsSpecifyTime   bool                  `form:"is_specify_time" json:"is_specify_time,omitempty"`
	Priority        int8                  `form:"priority" json:"priority,omitempty" binding:"required,oneof=1 2 3"`
	IsComplete      bool                  `form:"is_complete" json:"is_complete,omitempty"`
}

type TaskGetRequest struct {
	TableID
}
