package main

import (
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TimecodeLikeRepositorySuite struct {
	suite.Suite
	DB   *gorm.DB
	Repo *DBTimecodeLikeRepository
}

func (suite *TimecodeLikeRepositorySuite) SetupSuite() {
	suite.DB = TestDB
	suite.Repo = &DBTimecodeLikeRepository{DB: TestDB}
}

func (suite *TimecodeLikeRepositorySuite) SetupTest() {
	for _, table := range []string{
		"timecode_likes",
		"timecodes",
		"users",
	} {
		Cleaner.Acquire(table)
	}
}

func (suite *TimecodeLikeRepositorySuite) TearDownTest() {
	for _, table := range []string{
		"timecode_likes",
		"timecodes",
		"users",
	} {
		Cleaner.Clean(table)
	}
}

func TestTimecodeLikeRepositorySuite(t *testing.T) {
	suite.Run(t, new(TimecodeLikeRepositorySuite))
}

func (suite *TimecodeLikeRepositorySuite) TestDBTimecodeLikeRepository_Create() {
	st := suite.T()

	st.Run("when valid parameters given", func(t *testing.T) {
		user := &User{Email: "user1@example.com"}
		timecode := &Timecode{VideoID: "video-id-1", Description: "test"}
		suite.DB.Create(user)
		suite.DB.Create(timecode)

		timecodeLike := &TimecodeLike{TimecodeID: timecode.ID}
		timecodeLike, err := suite.Repo.Create(timecodeLike, user.ID)

		assert.Equal(t, user.ID, timecodeLike.UserID)
		assert.Equal(t, timecode.ID, timecodeLike.TimecodeID)
		assert.Nil(t, err)
	})

	st.Run("when invalid parameters given", func(t *testing.T) {
		user := &User{Email: "user2@example.com"}
		suite.DB.Create(user)

		timecodeLike := &TimecodeLike{TimecodeID: uint(1984)}
		timecodeLike, err := suite.Repo.Create(timecodeLike, user.ID)

		assert.True(t, suite.DB.NewRecord(timecodeLike))
		assert.EqualError(
			t, err,
			`pq: insert or update on table "timecode_likes" violates foreign key constraint "timecode_likes_timecode_id_timecodes_id_foreign"`)
	})
}

func (suite *TimecodeLikeRepositorySuite) TestDBTimecodeLikeRepository_Delete() {
	st := suite.T()

	st.Run("when record exists", func(t *testing.T) {
		user := &User{Email: "user3@example.com"}
		timecode := &Timecode{VideoID: "video-id-2", Description: "test"}
		suite.DB.Create(user)
		suite.DB.Create(timecode)
		timecodeLike := &TimecodeLike{TimecodeID: timecode.ID, UserID: user.ID}
		suite.DB.Create(timecodeLike)

		timecodeLike, err := suite.Repo.Delete(timecodeLike, user.ID)

		assert.True(t, suite.DB.First(timecodeLike).RecordNotFound())
		assert.Nil(t, err)
	})

	st.Run("when record doesn't exist", func(t *testing.T) {
		timecodeLike := &TimecodeLike{TimecodeID: 4, UserID: 5}

		timecodeLike, err := suite.Repo.Delete(timecodeLike, 6)

		assert.Nil(t, timecodeLike)
		assert.EqualError(t, err, "record not found")
	})
}
