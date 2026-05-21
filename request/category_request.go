package request

import (
	"time"
)

type ListCategoryRequest struct {
	Page     int    `form:"page,default=1" binding:"omitempty,min=1"`
	Limit    int    `form:"limit,default=10" binding:"omitempty,min=1"`
	Keyword  string `form:"keyword"`
	IsActive *bool  `form:"is_active"`
}

type ListCommonCategoryRequest struct {
	IsActive *bool `form:"is_active"`
}

type CreateCategoryRequest struct {
	Name     string `json:"name" binding:"required"`
	Code     string `json:"code" binding:"required"`
	IsActive *bool  `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name     string `json:"name"`
	Code     string `json:"code"`
	IsActive *bool  `json:"is_active"`
}

type UpdateCategoryStatusRequest struct {
	IsActive *bool `json:"is_active" binding:"required"`
}

type CategoryResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
