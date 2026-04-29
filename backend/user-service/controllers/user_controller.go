package controllers

import (
	"net/http"
	"strconv"

	"user-service/models"
	"user-service/models/dto"
	apperrors "user-service/pkg/errors"
	responseRes "user-service/pkg/response"
	"user-service/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	userService       services.IUserService
	authService       services.IAuthService
	memberService     services.IMemberService
	instructorService services.IInstructorService
	roleService       services.IRoleService
	emailService      services.IMailtrapEmailService
	mediaService      services.IMediaService
}

func NewUserController(
	userService services.IUserService,
	authService services.IAuthService,
	memberService services.IMemberService,
	instructorService services.IInstructorService,
	roleService services.IRoleService,
	emailService services.IMailtrapEmailService,
	mediaService services.IMediaService,
) IUserController {
	return &UserController{
		userService:       userService,
		authService:       authService,
		memberService:     memberService,
		instructorService: instructorService,
		roleService:       roleService,
		emailService:      emailService,
		mediaService:      mediaService,
	}
}

type IUserController interface {
	CreateUser(cxt *gin.Context)
	GetAllUsers(cxt *gin.Context)
	GetUserByID(cxt *gin.Context)
	UpdateUser(cxt *gin.Context)
	DeleteUser(cxt *gin.Context)
}

// @Summary Create User
// @Description Create a new user (admin only)
// @Tags Users
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "User data"
// @Success 201 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 409 {object} responseRes.Response
// @Router /users [post]
func (c *UserController) CreateUser(cxt *gin.Context) {
	var input dto.CreateUserRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	user, err := c.userService.CreateUser(input)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusCreated, "User created successfully", user)
}

// @Summary Get All Users
// @Description Get all users with their profiles
// @Tags Users
// @Produce json
// @Success 200 {object} responseRes.Response
// @Failure 500 {object} responseRes.Response
// @Router /users [get]
func (c *UserController) GetAllUsers(cxt *gin.Context) {
	users, err := c.userService.GetAllUsersWithProfiles()
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "Users retrieved successfully", users)
}

// @Summary Get User by ID
// @Description Get a user by their ID with profiles
// @Tags Users
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /users/{id} [get]
func (c *UserController) GetUserByID(cxt *gin.Context) {
	id, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	user, err := c.userService.GetUserByIDWithProfiles(id)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "User retrieved successfully", user)
}

// @Summary Update User
// @Description Update user information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Param request body dto.UpdateUserRequest true "Update user data"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /users/{id} [put]
func (c *UserController) UpdateUser(cxt *gin.Context) {
	id, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	var input dto.UpdateUserRequest
	if err := cxt.ShouldBindJSON(&input); err != nil {
		responseRes.ErrorFromAppError(cxt, apperrors.ErrBadRequest)
		return
	}

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	// Convert DTO to model for update
	userModel := &models.User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		EmailAddress: user.Email,
		PhoneNumber:  user.PhoneNumber,
		RoleID:       user.RoleID,
	}

	// Update fields if provided
	if input.Username != nil {
		userModel.Username = *input.Username
	}
	if input.EmailAddress != nil {
		userModel.Email = *input.EmailAddress
		userModel.EmailAddress = *input.EmailAddress
	}
	if input.PhoneNumber != nil {
		userModel.PhoneNumber = *input.PhoneNumber
	}
	if input.DateOfBirth != nil {
		userModel.DateOfBirth = *input.DateOfBirth
	}
	if input.Image != nil {
		userModel.Image = *input.Image
	}
	if input.Address != nil {
		userModel.Address = *input.Address
	}
	if input.RoleID != nil {
		userModel.RoleID = *input.RoleID
	}

	if err := c.userService.UpdateUser(userModel); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "User updated successfully", user)
}

// @Summary Delete User
// @Description Delete a user by ID
// @Tags Users
// @Produce json
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} responseRes.Response
// @Failure 400 {object} responseRes.Response
// @Failure 404 {object} responseRes.Response
// @Router /users/{id} [delete]
func (c *UserController) DeleteUser(cxt *gin.Context) {
	id, err := getUserIDFromPath(cxt, "id")
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	user, err := c.userService.GetUserByID(id)
	if err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	// Convert DTO to model for delete
	userModel := &models.User{
		ID: user.ID,
	}

	if err := c.userService.DeleteUser(userModel); err != nil {
		responseRes.ErrorFromGeneric(cxt, err)
		return
	}

	responseRes.Success(cxt, http.StatusOK, "User deleted successfully", nil)
}

// Helper to get user ID from path (returns UUID)
func getUserIDFromPath(c *gin.Context, paramName string) (uuid.UUID, error) {
	idStr := c.Param(paramName)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, apperrors.ErrBadRequest
	}
	return id, nil
}

// Helper to get uint ID from path
func getUintIDFromPath(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	if idStr == "" {
		return 0, apperrors.ErrBadRequest
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, apperrors.ErrBadRequest
	}
	return uint(id), nil
}