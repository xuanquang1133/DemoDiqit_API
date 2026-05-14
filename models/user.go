package models

import (
	"demodiqit_api/helpers/crypt"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// User represents the 'users' table in the database
type User struct {
	gorm.Model
	Username string         `gorm:"unique;not null" json:"username"`
	Password string         `gorm:"not null" json:"-"` // Do not return password in JSON
	Email     string         `gorm:"unique;not null" json:"email"`
	FullName  string         `json:"full_name"`
	Roles     pq.StringArray `gorm:"type:text[];default:'{}'" json:"roles"` // User authorization roles (multiple)
	UserToken string         `gorm:"type:text" json:"user_token"`           // Store the latest JWT token
}

// BeforeCreate hook to hash password before saving to the database
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Password != "" {
		hashedPassword, err := crypt.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
	}
	return
}
