package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents the 'orders' table in the database
type Order struct {
	ID              uint        `gorm:"primaryKey" json:"id"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	UserID          *uint       `gorm:"index" json:"user_id"`
	OrderNumber     string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"`
	Status          string      `gorm:"type:varchar(20);default:'pending';not null" json:"status"`
	Subtotal        float64     `gorm:"type:decimal(12,2);not null" json:"subtotal"`
	ShippingFee     float64     `gorm:"type:decimal(12,2);not null" json:"shipping_fee"`
	TotalAmount     float64     `gorm:"type:decimal(12,2);not null" json:"total_amount"`
	CustomerName    string      `gorm:"type:varchar(255);not null" json:"customer_name"`
	CustomerEmail   string      `gorm:"type:varchar(255);not null" json:"customer_email"`
	CustomerPhone   string      `gorm:"type:varchar(20);not null" json:"customer_phone"`
	ShippingAddress string      `gorm:"type:text;not null" json:"shipping_address"`
	Notes           string      `gorm:"type:text" json:"notes"`

	// Relations
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE;" json:"order_items"`
}
