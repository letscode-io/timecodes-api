package main

import (
	"encoding/json"
	"net/http"
)

func handleLogin(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(getCurrentUser(r))
}
