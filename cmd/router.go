package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Handler struct {
	*Container
	H func(c *Container, w http.ResponseWriter, r *http.Request)
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.H(h.Container, w, r)
}

func createRouter(container *Container) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	// public
	router.HandleFunc("/", handleHome)
	router.Handle("/timecodes/{videoId}", Handler{container, handleGetTimecodes})

	// auth
	auth := router.PathPrefix("/auth").Subrouter()
	auth.Use(authMiddleware(container))

	auth.HandleFunc("/login", handleLogin)

	auth.Handle("/timecodes", Handler{container, handleCreateTimecode}).Methods(http.MethodPost)
	auth.Handle("/timecodes/{videoId}", Handler{container, handleGetTimecodes})

	auth.Handle("/timecode_likes", Handler{container, handleCreateTimecodeLike}).Methods(http.MethodPost)
	auth.Handle("/timecode_likes", Handler{container, handleDeleteTimecodeLike}).Methods(http.MethodDelete)

	return cors.AllowAll().Handler(router)
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
