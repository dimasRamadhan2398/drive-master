package dto

type RoleResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type UpdateRoleInput struct {
	RoleID uint `json:"roleId" binding:"required"`
	
}