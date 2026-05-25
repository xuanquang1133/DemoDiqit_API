package request

import "time"

// OrderStatus constants
const (
	OrderStatusPending    = "pending"
	OrderStatusProcessing = "processing"
	OrderStatusCompleted  = "completed"
	OrderStatusCancelled  = "cancelled"
)

// OrderListQuery holds query parameters for listing orders
type OrderListQuery struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Keyword  string `form:"keyword"`
	Status   string `form:"status"`
	DateFrom string `form:"date_from"`
	DateTo   string `form:"date_to"`
}

// UpdateOrderStatusRequest is the DTO for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// OrderItemResponse is the response DTO for order items
type OrderItemResponse struct {
	ID              uint   `json:"id"`
	ProductID       *uint  `json:"product_id"`
	ProductName     string `json:"product_name"`
	ProductThumbnail string `json:"product_thumbnail"`
	Quantity        int    `json:"quantity"`
	Price           int64  `json:"price"`
	Subtotal        int64  `json:"subtotal"`
}

// OrderResponse is the response DTO for a single order
type OrderResponse struct {
	ID              uint                 `json:"id"`
	OrderNumber     string               `json:"order_number"`
	Status          string               `json:"status"`
	Subtotal        float64              `json:"subtotal"`
	ShippingFee     float64              `json:"shipping_fee"`
	TotalAmount     float64              `json:"total_amount"`
	CustomerName    string               `json:"customer_name"`
	CustomerEmail   string               `json:"customer_email"`
	CustomerPhone   string               `json:"customer_phone"`
	ShippingAddress string               `json:"shipping_address"`
	Notes           string               `json:"notes"`
	ItemsCount      int                  `json:"items_count"`
	OrderItems      []OrderItemResponse  `json:"order_items,omitempty"`
	CreatedAt       time.Time            `json:"created_at"`
	UpdatedAt       time.Time            `json:"updated_at"`
}

// OrderListResponse is the paginated response for order listing
type OrderListResponse struct {
	Items      []OrderResponse `json:"items"`
	Total      int64            `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}

// IsValidStatus checks if a status value is valid
func IsValidStatus(status string) bool {
	switch status {
	case OrderStatusPending, OrderStatusProcessing, OrderStatusCompleted, OrderStatusCancelled:
		return true
	default:
		return false
	}
}

// CreateOrderItemRequest is the DTO for an item in create order request
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// CreateOrderRequest is the DTO for creating a new order
type CreateOrderRequest struct {
	CustomerName    string                    `json:"customer_name" binding:"required"`
	CustomerEmail   string                    `json:"customer_email" binding:"required,email"`
	CustomerPhone   string                    `json:"customer_phone" binding:"required"`
	ShippingAddress string                    `json:"shipping_address" binding:"required"`
	ShippingFee     float64                   `json:"shipping_fee"`
	Notes           string                    `json:"notes"`
	Items           []CreateOrderItemRequest   `json:"items" binding:"required,min=1,dive"`
}
