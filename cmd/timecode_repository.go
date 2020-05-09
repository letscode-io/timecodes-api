package main

import (
	"github.com/jinzhu/gorm"
)

// Timecode represents timecode model
type Timecode struct {
	gorm.Model
	Description string         `json:"description"`
	Seconds     int            `json:"seconds" gorm:"not null"`
	VideoID     string         `json:"videoId" gorm:"not null;index"`
	Likes       []TimecodeLike `json:"likes" gorm:"foreignkey:TimecodeID"`
}
