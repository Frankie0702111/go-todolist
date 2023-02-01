package request

type CategoryGetListRequest struct {
	Id   int    `form:"id" json:"id,omitempty"`
	Name string `form:"name" json:"name,omitempty" binding:"max=100"`
	Pagination
}