package users

import (
	"net/http"

	"timecodes/pkg/models"
)

// CurrentUserKey used to store user struct in http context
type CurrentUserKey struct{}

// GetCurrentUser gets the current user from request context
func GetCurrentUser(r *http.Request) *models.User {
	user := r.Context().Value(CurrentUserKey{})
	if user != nil {
		return user.(*models.User)
	}

	return nil
}
