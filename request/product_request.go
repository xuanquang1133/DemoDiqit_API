package request

import "time"

// CreateProductRequest is the DTO for creating a new product
type CreateProductRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	SKU         string  `json:"sku" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required"`
	Thumbnail   string  `json:"thumbnail" binding:"required"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateProductRequest is the DTO for updating an existing product
type UpdateProductRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	SKU         string  `json:"sku"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Thumbnail   string  `json:"thumbnail"`
	IsActive    *bool   `json:"is_active"`
}

// UpdateProductStatusRequest is the DTO for updating product status
type UpdateProductStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// ProductListQuery holds the query parameters for listing products
type ProductListQuery struct {
	Page       int    `form:"page"`
	Limit      int    `form:"limit"`
	Keyword    string `form:"keyword"`
	IsCategory string `form:"is_category"`
	IsActive   *bool  `form:"is_active"`
}

// ProductCategoryInfo holds minimal category info for product responses
type ProductCategoryInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ProductResponse is the response DTO for a single product
type ProductResponse struct {
	ID          uint                 `json:"id"`
	CategoryID  *uint                `json:"category_id"`
	Category    *ProductCategoryInfo `json:"category,omitempty"`
	Name        string               `json:"name"`
	Slug        string               `json:"slug"`
	SKU         string               `json:"sku"`
	Description string               `json:"description"`
	Price       float64              `json:"price"`
	Thumbnail   string               `json:"thumbnail"`
	IsActive    bool                 `json:"is_active"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
}
