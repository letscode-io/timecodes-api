package models

import "github.com/jinzhu/gorm"

// TimecodeLike represents timecode like model
type TimecodeLike struct {
	gorm.Model
	TimecodeID uint `json:"timecodeId" gorm:"not null"`
	UserID     uint `json:"userId" gorm:"not null"`
}
