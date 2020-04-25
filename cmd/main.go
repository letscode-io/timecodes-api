package main

import (
	"context"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

var (
	db             *gorm.DB
	youtubeService *youtube.Service
)

func init() {
	initDB()
	createTables()
	runMigrations()
	initYoutubeService()
}

func initYoutubeService() {
	var err error
	ctx := context.Background()
	youtubeService, err = youtube.NewService(ctx, option.WithAPIKey(os.Getenv("GOOGLE_API_KEY")))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	startHttpServer()
}
