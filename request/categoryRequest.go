package request

type CategoryGetListRequest struct {
	Id   int64  `form:"id" json:"id,omitempty"`
	Name string `form:"name" json:"name,omitempty" binding:"max=100"`
	Pagination
}

type CategoryCreateOrUpdateRequest struct {
	Name string `form:"name" json:"name,omitempty" binding:"max=100"`
}

type CategoryGetRequest struct {
	TableID
}
