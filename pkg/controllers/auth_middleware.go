package controllers

import (
	"context"
	"net/http"
	"regexp"

	"timecodes/pkg/container"
	googleAPI "timecodes/pkg/google_api"
	"timecodes/pkg/users"
)

var authTokenRegExp = regexp.MustCompile(`Bearer (\S+$)`)

// AuthMiddleware middleware
func AuthMiddleware(c *container.Container) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			currentUser := users.GetCurrentUser(r)
			if currentUser != nil {
				next.ServeHTTP(w, r)
				return
			}

			token := getAuthToken(r.Header.Get("Authorization"))

			userInfo, err := googleAPI.FetchUserInfo(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			user := c.UserRepository.FindOrCreateByGoogleInfo(userInfo)
			ctx := context.WithValue(r.Context(), users.CurrentUserKey{}, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func getAuthToken(authorizationHeader string) string {
	matches := authTokenRegExp.FindSubmatch([]byte(authorizationHeader))

	if len(matches) == 0 {
		return ""
	}

	token := matches[len(matches)-1]

	return string(token)
}
