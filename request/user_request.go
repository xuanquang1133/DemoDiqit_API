package request

import "github.com/lib/pq"

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
}

type UpdateUserStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type UserResponse struct {
	ID       uint           `json:"id"`
	Username string         `json:"username"`
	Email    string         `json:"email"`
	FullName string         `json:"full_name"`
	Roles    pq.StringArray `json:"roles"`
	IsActive bool           `json:"is_active"`
}
