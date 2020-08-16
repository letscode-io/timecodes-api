package main

import (
	"testing"

	timecodeParser "timecodes/pkg/timecode_parser"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gopkg.in/khaiql/dbcleaner.v2"
)

type TimecodeRepositorySuite struct {
	suite.Suite

	VideoID        string
	AnotherVideoID string
	Cleaner        dbcleaner.DbCleaner
	DB             *gorm.DB
	Repo           *DBTimecodeRepository
}

func (suite *TimecodeRepositorySuite) SetupSuite() {
	cleaner := createDBCleaner(suite.T())
	db := initDB()
	runMigrations(db)

	suite.VideoID = "armenian-dram"
	suite.AnotherVideoID = "strategist"
	suite.Cleaner = cleaner
	suite.DB = db
	suite.Repo = &DBTimecodeRepository{DB: db}
}

func (suite *TimecodeRepositorySuite) SetupTest() {
	suite.Cleaner.Acquire("timecodes")
}

func (suite *TimecodeRepositorySuite) TearDownTest() {
	suite.Cleaner.Clean("timecodes")
}

func (suite *TimecodeRepositorySuite) TearDownSuite() {
	suite.DB.Close()
}

func TestTimecodeRepositorySuite(t *testing.T) {
	suite.Run(t, new(TimecodeRepositorySuite))
}

func (suite *TimecodeRepositorySuite) TestDBTimecodeRepository_FindByVideoID() {
	t := suite.T()

	t.Run("when matching records exist", func(t *testing.T) {
		suite.DB.Create(&Timecode{VideoID: suite.VideoID, Seconds: 55, Description: "ABC"})
		suite.DB.Create(&Timecode{VideoID: suite.VideoID, Seconds: 23, Description: "DEFG"})
		suite.DB.Create(&Timecode{VideoID: suite.AnotherVideoID, Seconds: 77, Description: "FGHJ"})
		defer suite.Cleaner.Clean("timecodes")

		timecodes := *suite.Repo.FindByVideoID(suite.VideoID)

		assert.Equal(t, 2, len(timecodes))
		assert.Equal(t, 23, timecodes[0].Seconds)
		assert.Equal(t, 55, timecodes[1].Seconds)
	})

	t.Run("when there are no matching records", func(t *testing.T) {
		suite.DB.Create(&Timecode{VideoID: suite.AnotherVideoID, Seconds: 77, Description: "FGHJ"})

		timecodes := *suite.Repo.FindByVideoID(suite.VideoID)

		assert.Equal(t, 0, len(timecodes))
	})
}

func (suite *TimecodeRepositorySuite) TestDBTimecodeRepository_Create() {
	t := suite.T()

	t.Run("when record has been created", func(t *testing.T) {
		timecode, err := suite.Repo.Create(&Timecode{VideoID: suite.VideoID, Seconds: 55, Description: "ABC"})

		assert.Nil(t, err)
		assert.NotNil(t, timecode.ID)
		assert.Equal(t, suite.VideoID, timecode.VideoID)
	})

	t.Run("when db returns an error", func(t *testing.T) {
		seconds := 10
		description := "ABC"

		suite.DB.Create(&Timecode{VideoID: suite.VideoID, Seconds: seconds, Description: description})

		timecode, err := suite.Repo.Create(&Timecode{VideoID: suite.VideoID, Seconds: seconds, Description: description})

		assert.True(t, suite.DB.NewRecord(timecode))
		assert.EqualError(t, err, `pq: duplicate key value violates unique constraint "idx_timecodes_seconds_text_video_id"`)
	})
}

func (suite *TimecodeRepositorySuite) TestDBTimecodeRepository_CreateFromParsedCodes() {
	suite.T().Run("when valid parsed codes have been given", func(t *testing.T) {
		parsedTimecodes := []timecodeParser.ParsedTimeCode{
			{Seconds: 24, Description: "Mariya Takeuchi - Oh Yes Oh No"},
			{Seconds: 24, Description: "Mariya Takeuchi Oh Yes Oh No"},
			{Seconds: 56, Description: "Toshiki Kadomatsu - Summer Emotions"},
			{Seconds: 56, Description: "Toshiki Kadomatsu Summer Emotions"},
			{Seconds: 88, Description: ""},
		}

		timecodes := *suite.Repo.CreateFromParsedCodes(parsedTimecodes, suite.VideoID)

		assert.Equal(t, 2, len(timecodes))
		assert.Equal(t, 24, timecodes[0].Seconds)
		assert.Equal(t, 56, timecodes[1].Seconds)
	})
}
