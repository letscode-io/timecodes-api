package main

import (
	"log"

	youtubeAPI "timecodes/cmd/youtube_api"
)

func main() {
	db := initDB()
	runMigrations(db)

	youtubeAPI, err := youtubeAPI.New()
	if err != nil {
		log.Fatal(err)
	}

	container := &Container{
		UserRepository:         &DBUserRepository{DB: db},
		TimecodeRepository:     &DBTimecodeRepository{DB: db},
		TimecodeLikeRepository: &DBTimecodeLikeRepository{DB: db},
		YoutubeAPI:             youtubeAPI,
	}

	startHttpServer(container)
}
