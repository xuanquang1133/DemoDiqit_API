package models

import (
	"time"

	"gorm.io/gorm"
)

// Cart represents a user's shopping cart (stored in DB for persistence)
type Cart struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	UserID    uint           `gorm:"uniqueIndex;not null" json:"user_id"`
	CartItems []CartItem     `gorm:"foreignKey:CartID;constraint:OnDelete:CASCADE;" json:"cart_items"`
}

// CartItem represents an item inside a cart
type CartItem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CartID      uint      `gorm:"uniqueIndex:idx_cart_product;not null" json:"cart_id"`
	ProductID   uint      `gorm:"uniqueIndex:idx_cart_product;not null" json:"product_id"`
	ProductName string    `gorm:"type:varchar(255);not null" json:"product_name"`
	Thumbnail   string    `gorm:"type:text" json:"thumbnail"`
	Price       float64   `gorm:"type:decimal(12,2);not null" json:"price"`
	Quantity    int       `gorm:"default:1;not null" json:"quantity"`
	Description string    `gorm:"type:text" json:"description"`
}
