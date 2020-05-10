package main

import (
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	googleAPI "timecodes/cmd/google_api"
)

var TestDB *gorm.DB

func TestDBUserRepository_FindOrCreateByGoogleInfo(t *testing.T) {
	TestDB = setupTestDB()
	defer TestDB.Close()

	t.Run("when user doesn't exist", func(t *testing.T) {
		googleID := "10001"
		repo := DBUserRepository{DB: TestDB}
		userInfo := &googleAPI.UserInfo{ID: googleID}

		user := repo.FindOrCreateByGoogleInfo(userInfo)

		assert.Equal(t, "10001", user.GoogleID)
	})

	t.Run("when user exists", func(t *testing.T) {
		googleID := "10002"

		TestDB.Create(&User{GoogleID: googleID})

		repo := DBUserRepository{DB: TestDB}
		userInfo := &googleAPI.UserInfo{ID: googleID}

		user := repo.FindOrCreateByGoogleInfo(userInfo)

		assert.Equal(t, "10002", user.GoogleID)
	})
}
