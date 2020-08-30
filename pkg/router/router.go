package router

import (
	"net/http"

	"timecodes/pkg/container"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Handler represents a helper structure for holding
type Handler struct {
	*container.Container
	H func(c *container.Container, w http.ResponseWriter, r *http.Request)
}

// ServeHTTP allows our Handler type to satisfy http.Handler.
func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.H(h.Container, w, r)
}

type routesApplier = func(rt *mux.Router, ctn *container.Container)

// Create creates a new mux.Router
func Create(c *container.Container, applyRoutes routesApplier) http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(contentTypeJSON)

	applyRoutes(router, c)

	return cors.AllowAll().Handler(router)
}

func contentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		next.ServeHTTP(w, r)
	})
}
