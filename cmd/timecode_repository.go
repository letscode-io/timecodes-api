package main

import (
	"strconv"
	timecodeParser "timecodes/cmd/timecode_parser"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
)

// Timecode represents timecode model
type Timecode struct {
	gorm.Model
	Description string         `json:"description" gorm:"not null;default:null"`
	Seconds     int            `json:"seconds" gorm:"not null"`
	VideoID     string         `json:"videoId" gorm:"not null;index"`
	Likes       []TimecodeLike `json:"likes" gorm:"foreignkey:TimecodeID"`
	UserID      uint           `json:"userId"`
}

type TimecodeRepository interface {
	FindByVideoId(string) *[]*Timecode
	Create(*Timecode) (*Timecode, error)
	CreateFromParsedCodes([]timecodeParser.ParsedTimeCode, string) *[]*Timecode
}

type DBTimecodeRepository struct {
	TimecodeRepository

	DB *gorm.DB
}

func (repo *DBTimecodeRepository) FindByVideoId(videoID string) *[]*Timecode {
	timecodes := &[]*Timecode{}

	repo.DB.Order("seconds asc").
		Preload("Likes").
		Where(&Timecode{VideoID: videoID}).
		Find(timecodes)

	return timecodes
}

func (repo *DBTimecodeRepository) Create(timecode *Timecode) (*Timecode, error) {
	err := repo.DB.Create(timecode).Error

	return timecode, err
}

func (repo *DBTimecodeRepository) CreateFromParsedCodes(parsedTimecodes []timecodeParser.ParsedTimeCode, videoId string) *[]*Timecode {
	seen := make(map[string]struct{})

	var collection []*Timecode
	for _, code := range parsedTimecodes {
		key := strconv.Itoa(code.Seconds) + strcase.ToCamel(code.Description)
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}

		user := getAdminUser(repo.DB)

		timecode := &Timecode{
			Seconds:     code.Seconds,
			VideoID:     videoId,
			Description: code.Description,
			UserID:      user.ID,
		}

		err := repo.DB.Create(timecode).Error
		if err != nil {
			continue
		}
		collection = append(collection, timecode)
	}

	return &collection
}
