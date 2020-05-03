package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func startHttpServer() {
	log.Println("Starting development server at http://127.0.0.1:8080/")

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	handler := cors.Default().Handler(router)

	// home
	router.HandleFunc("/", handleHome)

	// timecodes
	router.HandleFunc("/timecodes", createTimecode).Methods(http.MethodPost)
	router.HandleFunc("/timecodes/{videoId}", getTimecodes)

	// timecode_likes
	router.HandleFunc("/timecode_likes", handleCreateTimecodeLike).Methods(http.MethodPost)
	router.HandleFunc("/timecode_likes", handleDeleteTimecodeLike).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
