package main

// Annotation struct
type Annotation struct {
	ID      int64  `gorm:"primary_key"`
	Seconds int    `json:"seconds" gorm:"not null"`
	Text    string `json:"text"`
	VideoID string `json:"videoId" gorm:"not null;index"`
}
