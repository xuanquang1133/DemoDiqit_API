package request

import (
	"time"

	"github.com/shopspring/decimal"
)

// CreateProductRequest is the DTO for creating a new product
type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	SKU         string `json:"sku"`
	Description string `json:"description"`
	Price       string `json:"price" binding:"required"`
	SalePrice   string `json:"sale_price"`
	Thumbnail   string `json:"thumbnail"`
}

// UpdateProductRequest is the DTO for updating an existing product
type UpdateProductRequest struct {
	Name        string `json:"name"`
	SKU         string `json:"sku"`
	Description string `json:"description"`
	Price       string `json:"price"`
	SalePrice   string `json:"sale_price"`
	Thumbnail   string `json:"thumbnail"`
}

// UpdateProductStatusRequest is the DTO for updating product status
type UpdateProductStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// ProductListQuery holds the query parameters for listing products
type ProductListQuery struct {
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
	Keyword string `form:"keyword"`
}

// ProductResponse is the response DTO for a single product
type ProductResponse struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	SKU         string            `json:"sku"`
	Description string            `json:"description"`
	Price       decimal.Decimal   `json:"price"`
	SalePrice   decimal.Decimal   `json:"sale_price"`
	Thumbnail   string            `json:"thumbnail"`
	IsActive    bool              `json:"is_active"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// ProductListResponse is the paginated response for product listing
type ProductListResponse struct {
	Items      []ProductResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
