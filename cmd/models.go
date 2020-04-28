package main

// Timecode struct
type Timecode struct {
	ID          int64  `gorm:"primary_key"`
	Seconds     int    `json:"seconds" gorm:"not null"`
	Description string `json:"description"`
	VideoID     string `json:"videoId" gorm:"not null;index"`
}
