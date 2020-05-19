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
	err := args.Error(1)

	like.UserID = userID

	if like.TimecodeID == invalidTimecodeID {
		return nil, err
	}

	argLike := args.Get(0).(*TimecodeLike)
	argLike = like

	return argLike, nil
}

func (mr *MockTimecodeLikeRepository) Delete(like *TimecodeLike, userID uint) (*TimecodeLike, error) {
	args := mr.Called(like, userID)
	err := args.Error(1)

	like.UserID = userID

	if like.TimecodeID == invalidTimecodeID {
		return nil, err
	}

	argLike := args.Get(0).(*TimecodeLike)
	argLike = like

	return argLike, nil
}

func Test_handleCreateTimecodeLike(t *testing.T) {
	url := "/auth/timecode_likes"

	t.Run("when creation is successful", func(t *testing.T) {
		currentUser := &User{}
		currentUser.ID = 1

		like := &TimecodeLike{TimecodeID: 5}
		mockTKRepo.On("Create", like, currentUser.ID).Return(like, nil)

		params := []byte(`{ "timecodeId": 5 }`)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(params))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		mockTKRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("when creation has been failed", func(t *testing.T) {
		currentUser := &User{}
		currentUser.ID = 1

		like := &TimecodeLike{TimecodeID: invalidTimecodeID}
		mockTKRepo.On("Create", like, currentUser.ID).Return(nil, errors.New(""))

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

	t.Run("when deletion is successful", func(t *testing.T) {
		currentUser := &User{}
		currentUser.ID = 1

		like := &TimecodeLike{TimecodeID: 5}
		mockTKRepo.On("Delete", like, currentUser.ID).Return(like, nil)

		params := []byte(`{ "timecodeId": 5 }`)
		req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(params))

		response := executeRequest(t, timecodeLikeRouter, req, currentUser)

		mockTKRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when deletion has been failed", func(t *testing.T) {
		currentUser := &User{}
		currentUser.ID = 1

		like := &TimecodeLike{TimecodeID: invalidTimecodeID}
		mockTKRepo.On("Delete", like, currentUser.ID).Return(nil, errors.New(""))

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
