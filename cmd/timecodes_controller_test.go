package main

import (
	"bytes"
	"errors"
	"net/http"
	"testing"
	"time"

	timecodeParser "timecodes/pkg/timecode_parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockTimecodeRepo = &MockTimecodeRepository{}
var mockYTAPI = &mockYT{}
var timecodesContainer = &Container{
	YoutubeAPI:         mockYTAPI,
	TimecodeRepository: mockTimecodeRepo,
}
var timecodesRouter = createRouter(timecodesContainer)

type MockTimecodeRepository struct {
	mock.Mock
}

func (m *MockTimecodeRepository) FindByVideoID(videoID string) *[]*Timecode {
	args := m.Called(videoID)

	return args.Get(0).(*[]*Timecode)
}

func (m *MockTimecodeRepository) Create(timecode *Timecode) (*Timecode, error) {
	args := m.Called(timecode)

	return args.Get(0).(*Timecode), args.Error(1)
}

func (m *MockTimecodeRepository) CreateFromParsedCodes(parsedCodes []timecodeParser.ParsedTimeCode, videoID string) *[]*Timecode {
	args := m.Called(parsedCodes, videoID)

	return args.Get(0).(*[]*Timecode)
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

func Test_handleGetTimecodes(t *testing.T) {
	currentUser := &User{}
	currentUser.ID = 1

	t.Run("when timecodes exist", func(t *testing.T) {
		timecodes := &[]*Timecode{{}, {}, {}}

		mockTimecodeRepo.On("FindByVideoID", "video-id").Return(timecodes, nil)

		req, _ := http.NewRequest(http.MethodGet, "/timecodes/video-id", nil)

		response := executeRequest(t, timecodesRouter, req, currentUser)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when timecodes don't exist", func(t *testing.T) {
		timecodes := &[]*Timecode{}
		var emptyParsedCodes []timecodeParser.ParsedTimeCode

		mockTimecodeRepo.On("FindByVideoID", "no-items").Return(timecodes, nil)
		mockYTAPI.On("FetchVideoDescription", "no-items").Return("")
		mockYTAPI.On("FetchVideoComments", "no-items").Return([]string{})
		mockTimecodeRepo.
			On("CreateFromParsedCodes", emptyParsedCodes, "no-items").
			Return(timecodes, nil)

		req, _ := http.NewRequest(http.MethodGet, "/timecodes/no-items", nil)

		response := executeRequest(t, timecodesRouter, req, currentUser)

		time.Sleep(1 * time.Millisecond)

		mockYTAPI.AssertExpectations(t)
		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func Test_handleCreateTimecode(t *testing.T) {
	currentUser := &User{}
	currentUser.ID = 1

	t.Run("when request params are valid", func(t *testing.T) {
		timecode := &Timecode{
			VideoID:     "video-id",
			Seconds:     71,
			Description: "ABC",
			UserID:      currentUser.ID,
		}

		mockTimecodeRepo.On("Create", timecode).Return(timecode, nil)

		params := []byte(`{ "videoId": "video-id", "seconds": "1:11", "description": "ABC" }`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/timecodes", bytes.NewBuffer(params))

		response := executeRequest(t, timecodesRouter, req, currentUser)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("when request params are invalid", func(t *testing.T) {
		timecode := &Timecode{VideoID: "video-id", Seconds: 71, Description: "", UserID: currentUser.ID}

		mockTimecodeRepo.On("Create", timecode).Return(&Timecode{}, errors.New(""))

		params := []byte(`{ "videoId": "video-id", "seconds": "1:11", "description": "" }`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/timecodes", bytes.NewBuffer(params))

		response := executeRequest(t, timecodesRouter, req, currentUser)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("when request params contain invalid JSON", func(t *testing.T) {
		params := []byte(`{ "Invalid json }`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/timecodes", bytes.NewBuffer(params))

		response := executeRequest(t, timecodesRouter, req, currentUser)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
