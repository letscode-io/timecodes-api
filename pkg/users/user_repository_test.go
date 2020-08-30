package users

import (
	"testing"

	"timecodes/internal/db"
	testHelpers "timecodes/internal/test_helpers"
	googleAPI "timecodes/pkg/google_api"
	"timecodes/pkg/models"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
)

type UserRepositorySuite struct {
	suite.Suite

	Cleaner dbcleaner.DbCleaner
	DB      *gorm.DB
	Repo    *DBUserRepository
}

func (suite *UserRepositorySuite) SetupSuite() {
	database := db.Init()
	cleaner := testHelpers.CreateDBCleaner(suite.T(), database)

	suite.Cleaner = cleaner
	suite.DB = database.Connection
	suite.Repo = &DBUserRepository{DB: database.Connection}

	models.Migrate(database.Connection)
}

func (suite *UserRepositorySuite) SetupTest() {
	suite.Cleaner.Acquire("users")
}

func (suite *UserRepositorySuite) TearDownTest() {
	suite.Cleaner.Clean("users")
}

func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositorySuite))
}

func (suite *UserRepositorySuite) TestDBUserRepository_FindOrCreateByGoogleInfo() {
	t := suite.T()

	t.Run("when user doesn't exist", func(t *testing.T) {
		googleID := "10001"
		userInfo := &googleAPI.UserInfo{ID: googleID}

		user := suite.Repo.FindOrCreateByGoogleInfo(userInfo)

		assert.Equal(t, "10001", user.GoogleID)
	})

	t.Run("when user exists", func(t *testing.T) {
		googleID := "10002"

		suite.DB.Create(&models.User{GoogleID: googleID})

		userInfo := &googleAPI.UserInfo{ID: googleID}

		user := suite.Repo.FindOrCreateByGoogleInfo(userInfo)

		assert.Equal(t, "10002", user.GoogleID)
	})
}
