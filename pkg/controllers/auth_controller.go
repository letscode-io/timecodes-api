package controllers

import (
	"encoding/json"
	"net/http"

	"timecodes/pkg/users"
)

// HandleLogin GET /login
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users.GetCurrentUser(r))
}
