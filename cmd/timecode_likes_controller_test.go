package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const invalidTimecodeID = 1984

var invalidJSON = []byte("Invalid JSON")

var mockTKRepo = &MockTimecodeLikeRepository{}
var container = &Container{
	TimecodeLikeRepository: mockTKRepo,
}
var timecodeLikeRouter = createRouter(container)

type MockTimecodeLikeRepository struct {
	mock.Mock
}

func (mr *MockTimecodeLikeRepository) Create(like *TimecodeLike, userID uint) (*TimecodeLike, error) {
	args := mr.Called(like, userID)

	return args.Get(0).(*TimecodeLike), args.Error(1)
}

func (mr *MockTimecodeLikeRepository) Delete(like *TimecodeLike, userID uint) (*TimecodeLike, error) {
	args := mr.Called(like, userID)

	return args.Get(0).(*TimecodeLike), args.Error(1)
}

func Test_handleCreateTimecodeLike(t *testing.T) {
	url := "/auth/timecode_likes"
	currentUser := &User{}
	currentUser.ID = 1

	t.Run("when creation is successful", func(t *testing.T) {
		likeParams := &TimecodeLike{TimecodeID: 5}
		like := &TimecodeLike{TimecodeID: 5, UserID: currentUser.ID}

		mockTKRepo.On("Create", likeParams, currentUser.ID).Return(like, nil)

		params := []byte(`{ "timecodeId": 5 }`)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(params))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		mockTKRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("when creation has been failed", func(t *testing.T) {
		like := &TimecodeLike{TimecodeID: invalidTimecodeID}
		mockTKRepo.On("Create", like, currentUser.ID).Return(&TimecodeLike{}, errors.New(""))

		params := []byte(fmt.Sprintf(`{ "timecodeId": %d }`, invalidTimecodeID))
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(params))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		mockTKRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("when invalid request params has been given", func(t *testing.T) {
		currentUser := &User{}

		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(invalidJSON))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func Test_handleDeleteTimecodeLike(t *testing.T) {
	url := "/auth/timecode_likes"
	currentUser := &User{}
	currentUser.ID = 1

	t.Run("when deletion is successful", func(t *testing.T) {

		like := &TimecodeLike{TimecodeID: 5}
		mockTKRepo.On("Delete", like, currentUser.ID).Return(like, nil)

		params := []byte(`{ "timecodeId": 5 }`)
		req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(params))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		mockTKRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when deletion has been failed", func(t *testing.T) {
		like := &TimecodeLike{TimecodeID: invalidTimecodeID}
		mockTKRepo.On("Delete", like, currentUser.ID).Return(&TimecodeLike{}, errors.New(""))

		params := []byte(fmt.Sprintf(`{ "timecodeId": %d }`, invalidTimecodeID))
		req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(params))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		mockTKRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("when invalid request params has been given", func(t *testing.T) {
		currentUser := &User{}

		req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(invalidJSON))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
