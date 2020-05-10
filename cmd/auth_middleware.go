package main

import (
	"context"
	"net/http"
	"regexp"

	googleAPI "timecodes/cmd/google_api"
)

var authTokenRegExp = regexp.MustCompile(`Bearer (\S+$)`)

type CurrentUserKey struct{}

func authMiddleware(c *Container) (mw func(http.Handler) http.Handler) {
	mw = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := getAuthToken(r.Header.Get("Authorization"))

			userInfo, err := googleAPI.FetchUserInfo(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			user := c.UserRepository.FindOrCreateByGoogleInfo(userInfo)
			ctx := context.WithValue(r.Context(), CurrentUserKey{}, user)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
	return
}

func getAuthToken(authorizationHeader string) string {
	matches := authTokenRegExp.FindSubmatch([]byte(authorizationHeader))

	if len(matches) == 0 {
		return ""
	}

	token := matches[len(matches)-1]

	return string(token)
}
