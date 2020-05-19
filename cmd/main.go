package main

import (
	"log"
	"net/http"

	youtubeAPI "timecodes/cmd/youtube_api"
)

func main() {
	db := initDB()
	runMigrations(db)

	ytService, err := youtubeAPI.New()
	if err != nil {
		log.Fatal(err)
	}

	container := &Container{
		UserRepository:         &DBUserRepository{DB: db},
		TimecodeRepository:     &DBTimecodeRepository{DB: db},
		TimecodeLikeRepository: &DBTimecodeLikeRepository{DB: db},
		YoutubeAPI:             ytService,
	}

	router := createRouter(container)

	log.Println("Starting development server at http://127.0.0.1:8080/")
	log.Fatal(http.ListenAndServe(":8080", router))
}
