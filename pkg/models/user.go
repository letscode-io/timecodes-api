package models

import "github.com/jinzhu/gorm"

// User represents user model
type User struct {
	gorm.Model
	Email      string
	GoogleID   string
	PictureURL string
}
