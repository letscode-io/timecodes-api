package main

import "time"

// Timecode struct
type Timecode struct {
	ID          int64     `gorm:"primary_key"`
	Seconds     int       `json:"seconds" gorm:"not null"`
	Description string    `json:"description"`
	VideoID     string    `json:"videoId" gorm:"not null;index"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null"`
}
