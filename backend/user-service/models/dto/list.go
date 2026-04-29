package dto

type PaginationQuery struct {
	Page  int `form:"page"  binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}
