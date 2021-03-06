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
	ID      string `json:"id"`
	Picture string `json:"picture"`
	Email   string `json:"email"`
}

type APIError struct {
	ErrorData struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *APIError) Error() string {
	return e.ErrorData.Message
}

func FetchUserInfo(accessToken string) (userInfo *UserInfo, err error) {
	if len(accessToken) == 0 {
		return nil, errors.New("accessToken cannot be empty")
	}

	url := fmt.Sprintf("%s?access_token=%s", userInfoHost, accessToken)
	response, err := http.Get(url) // #nosec G107
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
		apiError := &APIError{}
		unmarshalErr := json.Unmarshal(contents, apiError)
		if unmarshalErr != nil {
			return nil, unmarshalErr
		}

		err = apiError
	default:
		userInfo = &UserInfo{}
		err = nil
		err = json.Unmarshal(contents, userInfo)
		if err != nil {
			return nil, err
		}
	}

	return userInfo, err
}
