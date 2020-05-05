package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

const userInfoHost = "https://www.googleapis.com/oauth2/v2/userinfo"

var authTokenRegExp = regexp.MustCompile(`Bearer (\S+$)`)

type UserInfo struct {
	Id      string `json:"id"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}
type CurrentUserKey struct{}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getAuthToken(r.Header.Get("Authorization"))

		userInfo, err := getUserInfo(token)
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

func getUserInfo(accessToken string) (userInfo *UserInfo, err error) {
	if len(accessToken) == 0 {
		return nil, errors.New("accessToken cannot be empty")
	}

	url := fmt.Sprintf("%s?access_token=%s", userInfoHost, accessToken)
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}

	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed reading response body: %s", err.Error())
	}

	switch response.StatusCode {
	case 401:
		userInfo, err = nil, errors.New(string(contents))
	default:
		userInfo = &UserInfo{}
		err = nil
		json.Unmarshal(contents, userInfo)
	}

	return userInfo, err
}

func findOrCreateUser(userInfo *UserInfo) *User {
	currentUser := &User{}

	db.Where(User{GoogleID: userInfo.Id}).
		Assign(User{Email: userInfo.Email, PictureURL: userInfo.Picture}).
		FirstOrCreate(&currentUser)

	return currentUser
}
