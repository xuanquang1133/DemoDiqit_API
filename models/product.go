package models

import (
	"gorm.io/gorm"
)

// Product represents the 'products' table in the database
type Product struct {
	gorm.Model
	CategoryID  *uint     `gorm:"index" json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Name        string    `gorm:"not null" json:"name"`
	Slug        string    `gorm:"uniqueIndex:idx_slug,where:deleted_at IS NULL" json:"slug"`
	SKU         string    `gorm:"uniqueIndex:idx_sku,where:deleted_at IS NULL" json:"sku"`
	Description string    `json:"description"`
	Price       float64   `gorm:"type:double precision" json:"price"`
	SalePrice   float64   `gorm:"type:double precision" json:"sale_price"`
	Thumbnail   string    `json:"thumbnail"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
}

// ExistsBySlug checks whether a non-deleted product with the given slug already exists.
// If excludeID is provided, that product is excluded from the check (used during updates).
func (p *Product) ExistsBySlug(tx *gorm.DB, slug string, excludeID ...uint) bool {
	q := tx.Where("slug = ?", slug)
	if len(excludeID) > 0 && excludeID[0] > 0 {
		q = q.Where("id != ?", excludeID[0])
	}
	var count int64
	q.Model(&Product{}).Count(&count)
	return count > 0
}

// ExistsBySKU checks whether a non-deleted product with the given SKU already exists.
// Empty SKU always returns false (SKU is optional).
// If excludeID is provided, that product is excluded from the check (used during updates).
func (p *Product) ExistsBySKU(tx *gorm.DB, sku string, excludeID ...uint) bool {
	if sku == "" {
		return false
	}
	q := tx.Where("sku = ?", sku)
	if len(excludeID) > 0 && excludeID[0] > 0 {
		q = q.Where("id != ?", excludeID[0])
	}
	var count int64
	q.Model(&Product{}).Count(&count)
	return count > 0
}
