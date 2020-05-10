package main

import (
	youtubeAPI "timecodes/cmd/youtube_api"
)

func main() {
	db := initDB()
	createTables(db)
	runMigrations(db)

	container := &Container{
		UserRepository:         &DBUserRepository{DB: db},
		TimecodeRepository:     &DBTimecodeRepository{DB: db},
		TimecodeLikeRepository: &DBTimecodeLikeRepository{DB: db},

		YoutubeAPI: youtubeAPI.New(),
	}

	startHttpServer(container)
}
