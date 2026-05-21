package models

import (
	"gorm.io/gorm"
)

// Category represents the 'categories' table in the database
type Category struct {
	gorm.Model
	Name    string `gorm:"uniqueIndex:idx_name_deleted_at,where:deleted_at IS NULL;not null" json:"name"`
	Code    string `gorm:"uniqueIndex:idx_code_deleted_at,where:deleted_at IS NULL;not null" json:"code"`
	IsActive bool  `gorm:"default:true" json:"is_active"`
}
