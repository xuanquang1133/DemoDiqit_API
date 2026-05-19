package models

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Product represents the 'products' table in the database
type Product struct {
	gorm.Model
	Name        string          `gorm:"not null" json:"name"`
	Slug        string          `gorm:"unique;not null" json:"slug"`
	SKU         string          `gorm:"unique" json:"sku"`
	Description string          `json:"description"`
	Price       decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"price"`
	SalePrice   decimal.Decimal `gorm:"type:decimal(15,2)" json:"sale_price"`
	Thumbnail   string          `json:"thumbnail"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
}
