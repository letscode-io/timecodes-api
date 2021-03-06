package youtubeapi

import (
	"net/http"
	"os"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var mockClient = http.DefaultClient

func TestMain(m *testing.M) {
	httpmock.ActivateNonDefault(mockClient)
	defer httpmock.DeactivateAndReset()

	registerMockResponders()

	os.Exit(m.Run())
}

func registerMockResponders() {
	legitVideoResponse := `
		{
			"items": [
				{
					"id": "legitVideoId",
					"snippet": {
						"description": "Legit description"
					}
				}
			]
		}
	`
	httpmock.RegisterResponder(
		"GET", "https://youtube.googleapis.com/youtube/v3/videos?alt=json&id=legitVideoId&part=snippet&prettyPrint=false",
		httpmock.NewStringResponder(200, legitVideoResponse),
	)

	notExistingVideoResponse := `{ "items": [] }`
	httpmock.RegisterResponder(
		"GET", "https://youtube.googleapis.com/youtube/v3/videos?alt=json&id=notExistingVideo&part=snippet&prettyPrint=false",
		httpmock.NewStringResponder(200, notExistingVideoResponse),
	)

	httpmock.RegisterResponder(
		"GET", "https://youtube.googleapis.com/youtube/v3/videos?alt=json&id=wrongResponse&part=snippet&prettyPrint=false",
		httpmock.NewStringResponder(500, "Something went wrong :("),
	)

	legitCommentsResponse := `
		{
			"items": [
				{
					"snippet": {
						"topLevelComment": {
							"snippet": {
								"textOriginal": "Just a comment."
							}
						}
					}
				}
			]
		}
	`
	httpmock.RegisterResponder(
		"GET", "https://youtube.googleapis.com/youtube/v3/commentThreads?alt=json&maxResults=100&order=relevance&part=snippet&prettyPrint=false&videoId=legitVideoId",
		httpmock.NewStringResponder(200, legitCommentsResponse),
	)
}

func TestNew(t *testing.T) {
	t.Run("when youtube.NewService returns an error", func(t *testing.T) {
		os.Setenv(GOOGLE_API_KEY, "")
		defer os.Unsetenv(GOOGLE_API_KEY)

		service, err := New()

		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "google: could not find default credentials")
	})

	t.Run("when returns correct service", func(t *testing.T) {
		os.Setenv(GOOGLE_API_KEY, "CORRECT_KEY")
		defer os.Unsetenv(GOOGLE_API_KEY)

		service, err := New()

		assert.NotNil(t, service)
		assert.Equal(t, "https://youtube.googleapis.com/", service.client.BasePath)
		assert.Nil(t, err)
	})
}

func TestNewWithClient(t *testing.T) {
	t.Run("when youtube.NewService returns an error", func(t *testing.T) {
		os.Setenv(GOOGLE_API_KEY, "")
		defer os.Unsetenv(GOOGLE_API_KEY)

		service, err := NewWithClient(nil)

		assert.Nil(t, service)
		assert.Contains(t, err.Error(), "google: could not find default credentials")
	})

	t.Run("when returns correct service", func(t *testing.T) {
		os.Setenv(GOOGLE_API_KEY, "CORRECT_KEY")
		defer os.Unsetenv(GOOGLE_API_KEY)

		service, err := NewWithClient(http.DefaultClient)

		assert.NotNil(t, service)
		assert.Equal(t, "https://youtube.googleapis.com/", service.client.BasePath)
		assert.Nil(t, err)
	})
}

func TestService_FetchVideoDescription(t *testing.T) {
	os.Setenv(GOOGLE_API_KEY, "CORRECT_KEY")
	defer os.Unsetenv(GOOGLE_API_KEY)
	service, _ := NewWithClient(mockClient)

	t.Run("when video exists", func(t *testing.T) {
		videoID := "legitVideoId"
		description := service.FetchVideoDescription(videoID)

		assert.Equal(t, "Legit description", description)
	})

	t.Run("when video doesn't exist", func(t *testing.T) {
		videoID := "notExistingVideo"
		description := service.FetchVideoDescription(videoID)

		assert.Equal(t, "", description)
	})

	t.Run("when http client returns an error", func(t *testing.T) {
		videoID := "wrongResponse"
		description := service.FetchVideoDescription(videoID)

		assert.Equal(t, "", description)
	})
}

func TestService_FetchVideoComments(t *testing.T) {
	os.Setenv(GOOGLE_API_KEY, "CORRECT_KEY")
	defer os.Unsetenv(GOOGLE_API_KEY)
	service, _ := NewWithClient(mockClient)

	t.Run("when video exists", func(t *testing.T) {
		videoID := "legitVideoId"
		comments := service.FetchVideoComments(videoID)

		assert.Equal(t, "Just a comment.", comments[0])
	})

	t.Run("when client returns an error", func(t *testing.T) {
		videoID := "wrongResponse"
		comments := service.FetchVideoComments(videoID)

		assert.Empty(t, comments)
	})
}
