package main

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email      string
	GoogleID   string
	PictureURL string
}
