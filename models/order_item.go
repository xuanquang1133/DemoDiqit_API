package models

import (
	"time"
)

// OrderItem represents an item within an order
// Prices are snapshots at the time of order placement
type OrderItem struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	OrderID   uint      `gorm:"index;not null" json:"order_id"`
	ProductID *uint     `gorm:"index" json:"product_id"`

	// Snapshot fields from Product at order time
	ProductName      string `gorm:"type:varchar(255);not null" json:"product_name"`
	ProductThumbnail string `gorm:"type:text" json:"product_thumbnail"`
	Quantity         int    `gorm:"not null" json:"quantity"`
	Price            int64  `gorm:"not null" json:"price"`     // Price per unit (in VND)
	Subtotal         int64  `gorm:"not null" json:"subtotal"`  // Quantity * Price

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
