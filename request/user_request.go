package request

import (
	"time"

	"github.com/lib/pq"
)

type ListUserRequest struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit    int    `form:"limit,default=10" binding:"omitempty,min=1"`
	Keyword  string `form:"keyword"`
	Role     string `form:"role"`
	IsActive *bool  `form:"is_active"`
}

type CreateUserRequest struct {
	Username string         `json:"username" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Email    string         `json:"email" binding:"required,email"`
	FullName string         `json:"full_name"`
	Roles    pq.StringArray `json:"roles"`
	IsActive *bool          `json:"is_active"`
}

type UpdateUserRequest struct {
	Username string         `json:"username"`
	Email    string         `json:"email"`
	FullName string         `json:"full_name"`
	Roles    pq.StringArray `json:"roles"`
	IsActive *bool          `json:"is_active"`
}

type UpdateUserStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type UserResponse struct {
	ID        uint           `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	FullName  string         `json:"full_name"`
	Roles     pq.StringArray `json:"roles"`
	IsActive  bool           `json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
}
