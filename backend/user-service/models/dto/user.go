package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username     string    `json:"username" binding:"required,min=2"`
	Password     string    `json:"password" binding:"required,min=8"`
	EmailAddress string    `json:"emailAddress" binding:"required,email"`
	PhoneNumber  string    `json:"phoneNumber" binding:"required,min=10"`
	Name         string    `json:"name"`
	Address      string    `json:"address"`
	DateOfBirth  time.Time `json:"dateOfBirth" binding:"required"`
	Image        string    `json:"image" binding:"required"`
	RoleID       uint      `json:"roleId" binding:"required"`
}

type CreateUserResponse struct {
	UserID      uuid.UUID      `json:"userId"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phoneNumber"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	RoleID      uint      `json:"roleId"`
}

type UpdateUserRequest struct {
	Username     *string    `json:"username" binding:"omitempty,min=2"`
	Password     *string    `json:"password" binding:"omitempty,min=8"`
	EmailAddress *string    `json:"emailAddress" binding:"omitempty,email"`
	PhoneNumber  *string    `json:"phoneNumber" binding:"omitempty,min=10"`
	DateOfBirth  *time.Time `json:"dateOfBirth" binding:"omitempty"`
	Image        *string    `json:"image" binding:"omitempty"`
	Address      *string    `json:"address" binding:"omitempty"`
	RoleID       *uint      `json:"roleId" binding:"omitempty"`
}

type UpdateUserResponse struct {
	UserID      uint      `json:"userId"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	PhoneNumber string    `json:"phoneNumber"`
	Image       string    `json:"image"`
	DateOfBirth time.Time `json:"dateOfBirth"`
	Address     string    `json:"address"`
	RoleID      uint      `json:"roleId"`
}

type GetUserRequest struct {
	Username     *string    `json:"username" binding:"omitempty,min=2"`
	Password     *string    `json:"password" binding:"omitempty,min=8"`
	EmailAddress *string    `json:"emailAddress" binding:"omitempty,email"`
	PhoneNumber  *string    `json:"phoneNumber" binding:"omitempty,min=10"`
	DateOfBirth  *time.Time `json:"dateOfBirth" binding:"omitempty"`
	Image        *string    `json:"image" binding:"omitempty"`
	Address      *string    `json:"address" binding:"omitempty"`
	RoleID       *uint      `json:"roleId" binding:"omitempty"`
}

type GetUserResponse struct {
	UserID      uuid.UUID         `json:"userId"`
	Email       string       `json:"email"`
	Username    string       `json:"username"`
	PhoneNumber string       `json:"phoneNumber"`
	Image       string       `json:"image"`
	DateOfBirth time.Time    `json:"dateOfBirth"`
	Address     string       `json:"address"`
	RoleID      uint         `json:"roleId"`
	Role        RoleResponse `json:"role"`
}

type UserWithProfileResponse struct {
	GetUserResponse
	MemberProfile     *MemberProfileResponse     `json:"memberProfile,omitempty"`
	InstructorProfile *InstructorProfileResponse `json:"instructorProfile,omitempty"`
}