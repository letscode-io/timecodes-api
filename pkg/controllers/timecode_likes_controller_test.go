package controllers

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	testHelpers "timecodes/internal/test_helpers"
	"timecodes/pkg/container"
	m "timecodes/pkg/models"
	"timecodes/pkg/router"
	"timecodes/pkg/users"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const invalidTimecodeID = 1984

var invalidJSON = []byte("Invalid JSON")

var mockTLRepo = &MockTimecodeLikeRepository{}
var mockTLContainer = &container.Container{
	TimecodeLikeRepository: mockTLRepo,
}

type MockTimecodeLikeRepository struct {
	mock.Mock
}

func (mr *MockTimecodeLikeRepository) Create(like *m.TimecodeLike, userID uint) (*m.TimecodeLike, error) {
	args := mr.Called(like, userID)

	return args.Get(0).(*m.TimecodeLike), args.Error(1)
}

func (mr *MockTimecodeLikeRepository) Delete(like *m.TimecodeLike, userID uint) (*m.TimecodeLike, error) {
	args := mr.Called(like, userID)

	return args.Get(0).(*m.TimecodeLike), args.Error(1)
}

func Test_HandleCreateTimecodeLike(t *testing.T) {
	url := "/auth/timecode_likes"
	currentUser := &m.User{}
	currentUser.ID = 1
	handler := router.Handler{Container: mockTLContainer, H: HandleCreateTimecodeLike}
	ctx := context.WithValue(context.Background(), users.CurrentUserKey{}, currentUser)
	path := ""

	t.Run("when creation is successful", func(t *testing.T) {
		likeParams := &m.TimecodeLike{TimecodeID: 5}
		like := &m.TimecodeLike{TimecodeID: 5, UserID: currentUser.ID}
		params := []byte(`{ "timecodeId": 5 }`)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(params))

		mockTLRepo.On("Create", likeParams, currentUser.ID).Return(like, nil)

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		mockTLRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("when creation has been failed", func(t *testing.T) {
		like := &m.TimecodeLike{TimecodeID: invalidTimecodeID}
		params := []byte(fmt.Sprintf(`{ "timecodeId": %d }`, invalidTimecodeID))
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(params))

		mockTLRepo.On("Create", like, currentUser.ID).Return(&m.TimecodeLike{}, errors.New(""))

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		mockTLRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("when invalid request params has been given", func(t *testing.T) {
		ctx = context.WithValue(context.Background(), users.CurrentUserKey{}, currentUser)
		req, _ := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(invalidJSON))

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func Test_HandleDeleteTimecodeLike(t *testing.T) {
	url := "/auth/timecode_likes"
	currentUser := &m.User{}
	currentUser.ID = 1
	handler := router.Handler{Container: mockTLContainer, H: HandleDeleteTimecodeLike}
	ctx := context.WithValue(context.Background(), users.CurrentUserKey{}, currentUser)
	path := ""

	t.Run("when deletion is successful", func(t *testing.T) {
		like := &m.TimecodeLike{TimecodeID: 5}
		params := []byte(`{ "timecodeId": 5 }`)
		req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(params))

		mockTLRepo.On("Delete", like, currentUser.ID).Return(like, nil)

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		mockTLRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("when deletion has been failed", func(t *testing.T) {
		like := &m.TimecodeLike{TimecodeID: invalidTimecodeID}
		params := []byte(fmt.Sprintf(`{ "timecodeId": %d }`, invalidTimecodeID))
		req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(params))

		mockTLRepo.On("Delete", like, currentUser.ID).Return(&m.TimecodeLike{}, errors.New(""))

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		mockTLRepo.AssertExpectations(t)
		assert.Equal(t, http.StatusUnprocessableEntity, response.Code)
	})

	t.Run("when invalid request params has been given", func(t *testing.T) {
		ctx = context.WithValue(context.Background(), users.CurrentUserKey{}, &m.User{})
		req, _ := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(invalidJSON))

		response := testHelpers.ExecuteRequest(ctx, t, req, handler, path)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}
