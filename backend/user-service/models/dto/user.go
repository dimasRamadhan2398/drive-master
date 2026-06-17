package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName    string    `json:"firstName" binding:"required,min=2"`
	LastName     string    `json:"lastName" binding:"required,min=2"`
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
	UserID      uuid.UUID `json:"userId"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	PhoneNumber string    `json:"phoneNumber"`
	DateOfBirth string    `json:"dateOfBirth"`
	RoleID      uint      `json:"roleId"`
}

type UpdateUserRequest struct {
	Username     *string    `json:"username" binding:"omitempty,min=2"`
	FirstName    string     `json:"firstName" binding:"required,min=2"`
	LastName     string     `json:"lastName" binding:"required,min=2"`
	Password     *string    `json:"password" binding:"omitempty,min=8"`
	EmailAddress *string    `json:"emailAddress" binding:"omitempty,email"`
	PhoneNumber  *string    `json:"phoneNumber" binding:"omitempty,min=10"`
	DateOfBirth  *time.Time `json:"dateOfBirth" binding:"omitempty"`
	Image        *string    `json:"image" binding:"omitempty"`
	Address      *string    `json:"address" binding:"omitempty"`
	RoleID       *uint      `json:"roleId" binding:"omitempty"`
}

type UpdateUserResponse struct {
	UserID      uint   `json:"userId"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	PhoneNumber string `json:"phoneNumber"`
	Image       string `json:"image"`
	DateOfBirth string `json:"dateOfBirth"`
	Address     string `json:"address"`
	RoleID      uint   `json:"roleId"`
}

type GetUserRequest struct {
	Username     *string    `json:"username" binding:"omitempty,min=2"`
	FirstName    string     `json:"firstName" binding:"required,min=2"`
	LastName     string     `json:"lastName" binding:"required,min=2"`
	Password     *string    `json:"password" binding:"omitempty,min=8"`
	EmailAddress *string    `json:"emailAddress" binding:"omitempty,email"`
	PhoneNumber  *string    `json:"phoneNumber" binding:"omitempty,min=10"`
	DateOfBirth  *time.Time `json:"dateOfBirth" binding:"omitempty"`
	Image        *string    `json:"image" binding:"omitempty"`
	Address      *string    `json:"address" binding:"omitempty"`
	RoleID       *uint      `json:"roleId" binding:"omitempty"`
}

type GetUserResponse struct {
	UserID      uuid.UUID    `json:"userId"`
	Email       string       `json:"email"`
	Username    string       `json:"username"`
	FirstName   string       `json:"firstName"`
	LastName    string       `json:"lastName"`
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

type RegistrationFilters struct {
	FromDate *time.Time `form:"fromDate" binding:"omitempty"`
	ToDate   *time.Time `form:"toDate" binding:"omitempty"`
}

type DashboardStatsResponse struct {
	TotalUsers          int64 `json:"totalUsers"`
	TotalMembers        int64 `json:"totalMembers"`
	TotalInstructors    int64 `json:"totalInstructors"`
	RecentRegistrations int64 `json:"recentRegistrations"`
}
