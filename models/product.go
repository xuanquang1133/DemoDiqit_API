package models

import (
	"regexp"
	"strings"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Product represents the 'products' table in the database
type Product struct {
	gorm.Model
	Name        string          `gorm:"not null" json:"name"`
	Slug        string          `gorm:"uniqueIndex:idx_slug,where:deleted_at IS NULL" json:"slug"`
	SKU         string          `gorm:"uniqueIndex:idx_sku,where:deleted_at IS NULL" json:"sku"`
	Description string          `json:"description"`
	Price       decimal.Decimal `gorm:"type:decimal(15,2);not null" json:"price"`
	SalePrice   decimal.Decimal `gorm:"type:decimal(15,2)" json:"sale_price"`
	Thumbnail   string          `json:"thumbnail"`
	IsActive    bool            `gorm:"default:true" json:"is_active"`
}

// BeforeSave normalizes Slug and SKU before persisting to the database.
// Slug: lowercase, replace non-alphanum runs with dash, collapse consecutive dashes.
// SKU:  uppercase, replace non-alphanum runs with dash, collapse consecutive dashes.
func (p *Product) BeforeSave(tx *gorm.DB) error {
	p.Slug = cleanSlug(p.Slug)
	p.SKU = cleanSKU(p.SKU)
	return nil
}

// cleanSlug normalizes a slug value: lowercase, trim, replace non-alphanum
// runs with a single dash, collapse consecutive dashes, trim leading/trailing dashes.
func cleanSlug(raw string) string {
	if raw == "" {
		return ""
	}
	cleaned := strings.ToLower(strings.TrimSpace(raw))
	cleaned = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(cleaned, "-")
	cleaned = regexp.MustCompile(`-+`).ReplaceAllString(cleaned, "-")
	return strings.Trim(cleaned, "-")
}

// cleanSKU normalizes a SKU value: uppercase, trim, replace non-alphanum
// runs with a single dash, collapse consecutive dashes, trim leading/trailing dashes.
func cleanSKU(raw string) string {
	if raw == "" {
		return ""
	}
	cleaned := strings.ToUpper(strings.TrimSpace(raw))
	cleaned = regexp.MustCompile(`[^A-Z0-9]+`).ReplaceAllString(cleaned, "-")
	cleaned = regexp.MustCompile(`-+`).ReplaceAllString(cleaned, "-")
	return strings.Trim(cleaned, "-")
}

// ExistsBySlug checks whether a non-deleted product with the given slug already exists.
// If excludeID is provided, that product is excluded from the check (used during updates).
func (p *Product) ExistsBySlug(tx *gorm.DB, slug string, excludeID ...uint) bool {
	normalized := cleanSlug(slug)
	if normalized == "" {
		return false
	}
	q := tx.Unscoped().Where("slug = ? AND deleted_at IS NULL", normalized)
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
	normalized := cleanSKU(sku)
	if normalized == "" {
		return false
	}
	q := tx.Unscoped().Where("sku = ? AND deleted_at IS NULL", normalized)
	if len(excludeID) > 0 && excludeID[0] > 0 {
		q = q.Where("id != ?", excludeID[0])
	}
	var count int64
	q.Model(&Product{}).Count(&count)
	return count > 0
}
