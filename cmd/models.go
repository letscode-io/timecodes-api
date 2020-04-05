package main

// Annotation struct
type Annotation struct {
	ID      string  `gorm:"primary_key"`
	Seconds float64 `json:"seconds"`
	Text    string  `json:"text"`
	VideoID string  `json:"videoId"`
}
