package main

import (
	"log"
	"net/http"

	"timecodes/internal/db"
	"timecodes/pkg/container"
	"timecodes/pkg/controllers"
	"timecodes/pkg/models"
	"timecodes/pkg/router"
	timecodeLikes "timecodes/pkg/timecode_likes"
	"timecodes/pkg/timecodes"
	"timecodes/pkg/users"
	youtubeAPI "timecodes/pkg/youtube_api"

	"github.com/gorilla/mux"
)

func main() {
	database := db.Init()

	models.Migrate(database.Connection)

	ytService, err := youtubeAPI.New()
	if err != nil {
		log.Fatal(err)
	}

	container := &container.Container{
		UserRepository:         &users.DBUserRepository{DB: database.Connection},
		TimecodeRepository:     &timecodes.DBTimecodeRepository{DB: database.Connection},
		TimecodeLikeRepository: &timecodeLikes.DBTimecodeLikeRepository{DB: database.Connection},
		YoutubeAPI:             ytService,
	}

	server := router.Create(container, routes)

	log.Println("Starting development server at http://127.0.0.1:8080/")
	log.Fatal(http.ListenAndServe(":8080", server))
}

func routes(mux *mux.Router, c *container.Container) {
	handlers := getHandlers(c)

	// public
	mux.HandleFunc("/", controllers.HandleRoot)
	mux.Handle("/timecodes/{videoId}", handlers["getTimecodes"])

	// auth
	auth := mux.PathPrefix("/auth").Subrouter()
	auth.Use(controllers.AuthMiddleware(c))

	auth.HandleFunc("/login", controllers.HandleLogin)

	auth.Handle("/timecodes", handlers["createTimecodes"]).Methods(http.MethodPost)
	auth.Handle("/timecodes/{videoId}", handlers["getTimecodes"])

	auth.Handle("/timecode_likes", handlers["createTimecodeLike"]).Methods(http.MethodPost)
	auth.Handle("/timecode_likes", handlers["deleteTimecodeLike"]).Methods(http.MethodDelete)
}

func getHandlers(c *container.Container) map[string]router.Handler {
	return map[string]router.Handler{
		"createTimecode":     {Container: c, H: controllers.HandleCreateTimecode},
		"createTimecodeLike": {Container: c, H: controllers.HandleCreateTimecodeLike},
		"deleteTimecodeLike": {Container: c, H: controllers.HandleDeleteTimecodeLike},
		"getTimecodes":       {Container: c, H: controllers.HandleGetTimecodes},
	}
}
