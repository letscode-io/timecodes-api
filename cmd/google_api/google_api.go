package googleapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const userInfoHost = "https://www.googleapis.com/oauth2/v2/userinfo"

type UserInfo struct {
	Id      string `json:"id"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

func FetchUserInfo(accessToken string) (userInfo *UserInfo, err error) {
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
