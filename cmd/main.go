package main

import (
	"github.com/jinzhu/gorm"

	youtubeAPI "timecodes/cmd/youtube_api"
)

var (
	db             *gorm.DB
	youtubeService *youtubeAPI.Service
)

func init() {
	initDB()
	createTables()
	runMigrations()
	initYoutubeService()
}

func initYoutubeService() {
	youtubeService = youtubeAPI.New()
}

func main() {
	startHttpServer()
}
