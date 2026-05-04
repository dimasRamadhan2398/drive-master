package controllers

import (
	"net/http"
	"strconv"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
)

type RoleController struct {
	roleSvc services.IRoleService
}

// GetAll implements [IRoleController].
func (r *RoleController) GetAll(ctx *gin.Context) {
	roles, err := r.roleSvc.FindAllRoles(ctx)
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusCreated, "List all roles", roles)
}

// GetRole implements [IRoleController].
func (r *RoleController) GetRole(ctx *gin.Context) {
	id := ctx.Param("id")

	num, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		responseRes.Error(ctx, http.StatusBadRequest, http.StatusText(http.StatusBadRequest), "Invalid role ID format", "")
		return
	}

	role, err := r.roleSvc.GetRole(ctx, uint(num))
	if err != nil {
		responseRes.ErrorFromGeneric(ctx, err)
		return
	}

	responseRes.Success(ctx, http.StatusOK, "Role retrieved successfully", role)
}

type IRoleController interface {
	GetAll(ctx *gin.Context)
	GetRole(ctx *gin.Context)
}

func NewRoleController(roleSvc services.IRoleService) IRoleController {
	return &RoleController{
		roleSvc: roleSvc,
	}
}
