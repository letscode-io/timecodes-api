package googleapi

import (
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("reader error")
}

func (errReader) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	registerMockResponders()

	os.Exit(m.Run())
}

func registerMockResponders() {
	httpmock.RegisterResponder(
		"GET", "https://www.googleapis.com/oauth2/v2/userinfo?access_token=tokenWithInvalidData",
		httpmock.ResponderFromResponse(
			&http.Response{
				Status:        "200",
				StatusCode:    200,
				Body:          errReader(0),
				Header:        http.Header{},
				ContentLength: 1,
			},
		),
	)

	unauthorizedError := `
		{
			"error": {
				"code": 401,
				"message": "Request is missing required authentication credential.",
				"status": "UNAUTHENTICATED"
			}
		}
	`
	httpmock.RegisterResponder(
		"GET", "https://www.googleapis.com/oauth2/v2/userinfo?access_token=invalidToken",
		httpmock.NewStringResponder(401, unauthorizedError),
	)

	userInfoResponse := `
		{
			"id": "92752",
			"email": "better_call_your_local_dealer@gmail.com",
			"verified_email": true,
			"picture": "https://lh3.googleusercontent.com/a-/0f78ce08-876e-4ffa-a227-84f52a69063f"
		}
	`
	httpmock.RegisterResponder(
		"GET", "https://www.googleapis.com/oauth2/v2/userinfo?access_token=validToken",
		httpmock.NewStringResponder(200, userInfoResponse),
	)
}

func TestFetchUserInfo(t *testing.T) {
	t.Run("when given accessToken is empty", func(t *testing.T) {
		token := ""
		userInfo, err := FetchUserInfo(token)

		assert.Nil(t, userInfo)
		assert.EqualError(t, err, "accessToken cannot be empty")
	})

	t.Run("when http client error occurs", func(t *testing.T) {
		token := "tokenWithClientError"
		userInfo, err := FetchUserInfo(token)

		assert.Nil(t, userInfo)
		assert.EqualError(
			t,
			err,
			`failed getting user info: Get "https://www.googleapis.com/oauth2/v2/userinfo?access_token=tokenWithClientError": no responder found`,
		)
	})

	t.Run("when response body contains invalid data", func(t *testing.T) {
		token := "tokenWithInvalidData"
		userInfo, err := FetchUserInfo(token)

		assert.Nil(t, userInfo)
		assert.EqualError(t, err, "failed reading response body: reader error")
	})

	t.Run("when given accessToken return unauthorized response", func(t *testing.T) {
		token := "invalidToken"
		userInfo, err := FetchUserInfo(token)

		assert.Nil(t, userInfo)
		assert.EqualError(t, err, "Request is missing required authentication credential.")
	})

	t.Run("when valid access token given", func(t *testing.T) {
		token := "validToken"
		userInfo, err := FetchUserInfo(token)

		assert.NotNil(t, userInfo)
		assert.Equal(t, "92752", userInfo.ID)
		assert.Nil(t, err)
	})
}
