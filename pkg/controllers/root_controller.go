package controllers

import (
	"fmt"
	"net/http"
)

// HandleRoot GET /
func HandleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "OK")
}
