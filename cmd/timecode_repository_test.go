package main

import (
	"testing"
	timecodeParser "timecodes/cmd/timecode_parser"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const videoID = "armenian-dram"
const anotherVideoID = "strategist"

type TimecodeRepositorySuite struct {
	suite.Suite
	DB   *gorm.DB
	Repo *DBTimecodeRepository
}

func (suite *TimecodeRepositorySuite) ResetDB() {
	suite.DB.Exec("TRUNCATE TABLE timecodes;")
}

func (suite *TimecodeRepositorySuite) SetupTest() {
	db := setupTestDB()
	suite.DB = db
	suite.Repo = &DBTimecodeRepository{DB: db}
}

func TestTimecodeRepositorySuite(t *testing.T) {
	suite.Run(t, new(TimecodeRepositorySuite))
}

func (suite *TimecodeRepositorySuite) TestDBTimecodeRepository_FindByVideoId() {
	t := suite.T()

	t.Run("when matching records exist", func(t *testing.T) {
		suite.DB.Create(&Timecode{VideoID: videoID, Seconds: 55, Description: "ABC"})
		suite.DB.Create(&Timecode{VideoID: videoID, Seconds: 23, Description: "DEFG"})
		suite.DB.Create(&Timecode{VideoID: anotherVideoID, Seconds: 77, Description: "FGHJ"})
		defer suite.ResetDB()

		timecodes := *suite.Repo.FindByVideoId(videoID)

		assert.Equal(t, 2, len(timecodes))
		assert.Equal(t, 23, timecodes[0].Seconds)
		assert.Equal(t, 55, timecodes[1].Seconds)
	})

	t.Run("when there are no matching records", func(t *testing.T) {
		suite.DB.Create(&Timecode{VideoID: anotherVideoID, Seconds: 77, Description: "FGHJ"})
		defer suite.ResetDB()

		timecodes := *suite.Repo.FindByVideoId(videoID)

		assert.Equal(t, 0, len(timecodes))
	})
}

func (suite *TimecodeRepositorySuite) TestDBTimecodeRepository_Create() {
	t := suite.T()

	t.Run("when record has been created", func(t *testing.T) {
		defer suite.ResetDB()

		timecode, err := suite.Repo.Create(&Timecode{VideoID: videoID, Seconds: 55, Description: "ABC"})

		assert.Nil(t, err)
		assert.NotNil(t, timecode.ID)
		assert.Equal(t, videoID, timecode.VideoID)
	})

	t.Run("when db returns an error", func(t *testing.T) {
		seconds := 10
		description := "ABC"

		suite.DB.Create(&Timecode{VideoID: videoID, Seconds: seconds, Description: description})
		defer suite.ResetDB()

		timecode, err := suite.Repo.Create(&Timecode{VideoID: videoID, Seconds: seconds, Description: description})

		assert.True(t, suite.DB.NewRecord(timecode))
		assert.EqualError(t, err, `pq: duplicate key value violates unique constraint "idx_timecodes_seconds_text_video_id"`)
	})
}

func (suite *TimecodeRepositorySuite) TestDBTimecodeRepository_CreateFromParsedCodes() {
	suite.T().Run("when valid parsed codes has been given", func(t *testing.T) {
		parsedTimecodes := []timecodeParser.ParsedTimeCode{
			{Seconds: 24, Description: "ABC"},
			{Seconds: 24, Description: "ABC"},
			{Seconds: 56, Description: "DFG"},
			{Seconds: 56, Description: "DFG"},
			{Seconds: 56, Description: "DFG"},
		}
		defer suite.ResetDB()

		timecodes := *suite.Repo.CreateFromParsedCodes(parsedTimecodes, videoID)

		assert.Equal(t, 2, len(timecodes))
		assert.Equal(t, 24, timecodes[0].Seconds)
		assert.Equal(t, 56, timecodes[1].Seconds)
	})
}
