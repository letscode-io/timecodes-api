package main

import "github.com/jinzhu/gorm"

// Timecode struct
type Timecode struct {
	gorm.Model
	Description string `json:"description"`
	Seconds     int    `json:"seconds" gorm:"not null"`
	VideoID     string `json:"videoId" gorm:"not null;index"`
}

type User struct {
	gorm.Model
}

type TimecodeLike struct {
	gorm.Model
	TimecodeID uint   `gorm:"not null"`
	UserID     uint   `gorm:"not null"`
	VideoID    string `gorm:"not null"`
}
