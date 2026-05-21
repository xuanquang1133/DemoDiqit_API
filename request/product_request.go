package request

import "time"

// CreateProductRequest is the DTO for creating a new product
type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Slug        string   `json:"slug" binding:"required"`
	SKU         string   `json:"sku" binding:"required"`
	Description string   `json:"description"`
	Price       float64  `json:"price" binding:"required"`
	SalePrice   *float64 `json:"sale_price"`
	Thumbnail   string   `json:"thumbnail" binding:"required"`
}

// UpdateProductRequest is the DTO for updating an existing product
type UpdateProductRequest struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	SKU         string   `json:"sku"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	SalePrice   *float64 `json:"sale_price"`
	Thumbnail   string   `json:"thumbnail"`
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
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	SKU         string    `json:"sku"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	SalePrice   float64   `json:"sale_price"`
	Thumbnail   string    `json:"thumbnail"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
