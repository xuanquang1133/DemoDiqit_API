package models

import (
	"gorm.io/gorm"
	"github.com/lib/pq"
)

// User represents the 'users' table in the database
type User struct {
	gorm.Model           // Contains ID, CreatedAt, UpdatedAt, DeletedAt
	Username string `gorm:"unique;not null" json:"username"`
	Password string `gorm:"not null" json:"-"` // Do not return password in JSON
	Email    string `gorm:"unique;not null" json:"email"`
	FullName string `json:"full_name"`
	Roles    pq.StringArray `gorm:"type:text[];default:'{}'" json:"roles"` // User authorization roles (multiple)
}
