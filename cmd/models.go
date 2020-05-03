package main

import "github.com/jinzhu/gorm"

// Timecode struct
type Timecode struct {
	gorm.Model
	Description string         `json:"description"`
	Seconds     int            `json:"seconds" gorm:"not null"`
	VideoID     string         `json:"videoId" gorm:"not null;index"`
	Likes       []TimecodeLike `json:"likes" gorm:"foreignkey:TimecodeID"`
}

type User struct {
	gorm.Model
	AccessToken AccessToken
}

type AccessToken struct {
	gorm.Model
	Value  string `gorm:"not null;unique_index"`
	UserID uint   `gorm:"not null;index"`
}

type TimecodeLike struct {
	gorm.Model
	TimecodeID uint `json:"timecodeId" gorm:"not null"`
	UserID     uint `json:"userId" gorm:"not null"`
}
