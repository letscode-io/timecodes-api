package main

import (
	"bytes"
	"net/http"
	"testing"
	timecodeParser "timecodes/cmd/timecode_parser"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var mockTimecodeRepo = &MockTimecodeRepository{}
var timecodesContainer = &Container{
	TimecodeRepository: mockTimecodeRepo,
}
var timecodesRouter = createRouter(timecodesContainer)

type MockTimecodeRepository struct {
	mock.Mock
}

func (m *MockTimecodeRepository) FindByVideoId(videoID string) *[]*Timecode {
	args := m.Called(videoID)

	collection := args.Get(0).(*[]*Timecode)
	if videoID == "no-items" {
		return collection
	}

	collection = &[]*Timecode{{}, {}, {}}

	return collection
}

func (m *MockTimecodeRepository) Create(timecode *Timecode) (*Timecode, error) {
	args := m.Called(timecode)

	err := args.Error(1)

	if len(timecode.Description) == 0 {
		return nil, err
	}

	_ = args.Get(0).(*Timecode)

	return timecode, nil
}

func (m *MockTimecodeRepository) CreateFromParsedCodes(parsedCodes []timecodeParser.ParsedTimeCode, videoID string) *[]*Timecode {
	args := m.Called(parsedCodes, videoID)

	return args.Get(0).(*[]*Timecode)
}

func Test_handleGetTimecodes(t *testing.T) {
	currentUser := &User{}
	currentUser.ID = 1

	t.Run("when timecodes exist", func(t *testing.T) {
		timecodes := &[]*Timecode{{}, {}, {}}

		mockTimecodeRepo.On("FindByVideoId", "video-id").Return(timecodes, nil)

		req, _ := http.NewRequest(http.MethodGet, "/timecodes/video-id", nil)

		response := executeRequest(t, timecodesRouter, req, currentUser)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when timecodes don't exist", func(t *testing.T) {
		t.Skip()
		timecodes := &[]*Timecode{}

		mockTimecodeRepo.On("FindByVideoId", "no-items").Return(timecodes, nil)

		req, _ := http.NewRequest(http.MethodGet, "/timecodes/video-id", nil)

		response := executeRequest(t, timecodesRouter, req, currentUser)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func Test_handleCreateTimecode(t *testing.T) {
	currentUser := &User{}
	currentUser.ID = 1

	t.Run("when request params are valid", func(t *testing.T) {
		timecode := &Timecode{VideoID: "video-id", Seconds: 111, Description: "ABC"}

		mockTimecodeRepo.On("Create", timecode).Return(timecode, nil)

		params := []byte(`{ "videoId": "video-id", "seconds": 111, "description": "ABC" }`)
		req, _ := http.NewRequest(http.MethodPost, "/auth/timecodes", bytes.NewBuffer(params))

		response := executeRequest(t, timecodesRouter, req, currentUser)

		mockTimecodeRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})
}
