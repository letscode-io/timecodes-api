package controllers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	testHelpers "timecodes/internal/test_helpers"
	"timecodes/pkg/container"
	"timecodes/pkg/models"
	"timecodes/pkg/router"
	timecodeParser "timecodes/pkg/timecode_parser"
	"timecodes/pkg/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockTimecodeRepo = &MockTimecodeRepository{}
var mockYTAPI = &mockYT{}
var timecodesContainer = &container.Container{
	YoutubeAPI:         mockYTAPI,
	TimecodeRepository: mockTimecodeRepo,
}

type MockTimecodeRepository struct {
	mock.Mock
}

func (m *MockTimecodeRepository) FindByVideoID(videoID string) *[]*models.Timecode {
	args := m.Called(videoID)

	return args.Get(0).(*[]*models.Timecode)
}

func (m *MockTimecodeRepository) Create(timecode *models.Timecode) (*models.Timecode, error) {
	args := m.Called(timecode)

	return args.Get(0).(*models.Timecode), args.Error(1)
}

func (m *MockTimecodeRepository) CreateFromParsedCodes(parsedCodes []timecodeParser.ParsedTimeCode, videoID string) *[]*models.Timecode {
	args := m.Called(parsedCodes, videoID)

	return args.Get(0).(*[]*models.Timecode)
}

type mockYT struct {
	mock.Mock
}

func (m *mockYT) FetchVideoDescription(videoID string) string {
	args := m.Called(videoID)

	return args.Get(0).(string)
}

func (m *mockYT) FetchVideoComments(videoID string) []string {
	args := m.Called(videoID)

	return args.Get(0).([]string)
}

func Test_HandleGetTimecodes(t *testing.T) {
	currentUser := &models.User{}
	currentUser.ID = 1
	handler := router.Handler{Container: timecodesContainer, H: HandleGetTimecodes}
	path := "/timecodes/{videoId}"
	ctx := context.WithValue(context.Background(), users.CurrentUserKey{}, currentUser)

	t.Run("when timecodes exist", func(t *testing.T) {
		timecodes := &[]*models.Timecode{{}, {}, {}}
		req, _ := http.NewRequest(http.MethodGet, "/timecodes/video-id", nil)

		mockTimecodeRepo.On("FindByVideoID", "video-id").Return(timecodes, nil)

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when timecodes don't exist", func(t *testing.T) {
		timecodes := &[]*models.Timecode{}
		var emptyParsedCodes []timecodeParser.ParsedTimeCode
		req, _ := http.NewRequest(http.MethodGet, "/timecodes/no-items", nil)

		mockTimecodeRepo.On("FindByVideoID", "no-items").Return(timecodes, nil)
		mockYTAPI.On("FetchVideoDescription", "no-items").Return("")
		mockYTAPI.On("FetchVideoComments", "no-items").Return([]string{})
		mockTimecodeRepo.
			On("CreateFromParsedCodes", emptyParsedCodes, "no-items").
			Return(timecodes, nil)

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		time.Sleep(1 * time.Millisecond)

		mockYTAPI.AssertExpectations(t)
		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func Test_HandleCreateTimecode(t *testing.T) {
	currentUser := &models.User{}
	currentUser.ID = 1
	handler := router.Handler{Container: timecodesContainer, H: HandleCreateTimecode}
	path := ""
	ctx := context.WithValue(context.Background(), users.CurrentUserKey{}, currentUser)

	t.Run("when request params are valid", func(t *testing.T) {
		timecode := &models.Timecode{
			VideoID:     "video-id",
			Seconds:     71,
			Description: "ABC",
			UserID:      currentUser.ID,
		}
		params := []byte(`{ "videoId": "video-id", "seconds": "1:11", "description": "ABC" }`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/timecodes", bytes.NewBuffer(params))

		mockTimecodeRepo.On("Create", timecode).Return(timecode, nil)

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("when request params are invalid", func(t *testing.T) {
		timecode := &models.Timecode{VideoID: "video-id", Seconds: 71, Description: "", UserID: currentUser.ID}

		mockTimecodeRepo.On("Create", timecode).Return(&models.Timecode{}, errors.New(""))

		params := []byte(`{ "videoId": "video-id", "seconds": "1:11", "description": "" }`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/timecodes", bytes.NewBuffer(params))

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("when request params contain invalid JSON", func(t *testing.T) {
		params := []byte(`{ "Invalid json }`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/timecodes", bytes.NewBuffer(params))

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
