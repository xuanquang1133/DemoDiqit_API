package request

import (
	"time"

	"github.com/lib/pq"
)

type CreateUserRequest struct {
	Username string         `json:"username" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Email    string         `json:"email" binding:"required,email"`
	FullName string         `json:"full_name"`
	Roles    pq.StringArray `json:"roles"`
	Status   string         `json:"status"`
}

type UpdateUserRequest struct {
	Username string         `json:"username"`
	Email    string         `json:"email"`
	FullName string         `json:"full_name"`
	Roles    pq.StringArray `json:"roles"`
	Status   string         `json:"status"`
}



type UserResponse struct {
	ID        uint           `json:"id"`
	Username  string         `json:"username"`
	Email     string         `json:"email"`
	FullName  string         `json:"full_name"`
	Roles     pq.StringArray `json:"roles"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
}
