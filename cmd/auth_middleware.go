package main

import (
	"context"
	"net/http"
	"regexp"
)

var authTokenRegExp = regexp.MustCompile(`Bearer (\S+$)`)

type CurrentUserKey struct{}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getAuthToken(r.Header.Get("Authorization"))

		userInfo, err := fetchUserInfo(token)
		if err != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		ctx := context.WithValue(r.Context(), CurrentUserKey{}, findOrCreateUser(userInfo))
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getAuthToken(authorizationHeader string) string {
	matches := authTokenRegExp.FindSubmatch([]byte(authorizationHeader))

	if len(matches) == 0 {
		return ""
	}

	token := matches[len(matches)-1]

	return string(token)
}

func findOrCreateUser(userInfo *UserInfo) *User {
	currentUser := &User{}

	db.Where(User{GoogleID: userInfo.Id}).
		Assign(User{Email: userInfo.Email, PictureURL: userInfo.Picture}).
		FirstOrCreate(&currentUser)

	return currentUser
}
