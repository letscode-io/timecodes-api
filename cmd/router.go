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

	// public
	router.HandleFunc("/", handleHome)
	router.HandleFunc("/timecodes/{videoId}", handleGetTimecodes)

	// auth
	auth := router.PathPrefix("/auth").Subrouter()
	auth.Use(authMiddleware)

	auth.HandleFunc("/login", handleLogin)

	auth.HandleFunc("/timecodes", handleCreateTimecode).Methods(http.MethodPost)
	auth.HandleFunc("/timecodes/{videoId}", handleGetTimecodes)

	auth.HandleFunc("/timecode_likes", handleCreateTimecodeLike).Methods(http.MethodPost)
	auth.HandleFunc("/timecode_likes", handleDeleteTimecodeLike).Methods(http.MethodDelete)

	handler := cors.Default().Handler(router)

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

func getCurrentUser(r *http.Request) *User {
	user := r.Context().Value(CurrentUserKey{})
	if user != nil {
		return user.(*User)
	}

	return nil
}
