package timecodes

import (
	"strconv"

	m "timecodes/pkg/models"
	timecodeParser "timecodes/pkg/timecode_parser"
	"timecodes/pkg/users"

	"github.com/iancoleman/strcase"
	"github.com/jinzhu/gorm"
)

// TimecodeRepository represents repository interface
type TimecodeRepository interface {
	FindByVideoID(string) *[]*m.Timecode
	Create(*m.Timecode) (*m.Timecode, error)
	CreateFromParsedCodes([]timecodeParser.ParsedTimeCode, string) *[]*m.Timecode
}

// DBTimecodeRepository implements TimecodeRepository
type DBTimecodeRepository struct {
	TimecodeRepository

	DB *gorm.DB
}

// FindByVideoID finds timecode by given video id
func (repo *DBTimecodeRepository) FindByVideoID(videoID string) *[]*m.Timecode {
	timecodes := &[]*m.Timecode{}

	repo.DB.Order("seconds asc").
		Preload("Likes").
		Where(&m.Timecode{VideoID: videoID}).
		Find(timecodes)

	return timecodes
}

// Create creates a new timecode record
func (repo *DBTimecodeRepository) Create(timecode *m.Timecode) (*m.Timecode, error) {
	err := repo.DB.Create(timecode).Error

	return timecode, err
}

// CreateFromParsedCodes creates timecodes from parsed codes
func (repo *DBTimecodeRepository) CreateFromParsedCodes(parsedTimecodes []timecodeParser.ParsedTimeCode, videoID string) *[]*m.Timecode {
	seen := make(map[string]struct{})

	var collection []*m.Timecode
	for _, code := range parsedTimecodes {
		key := strconv.Itoa(code.Seconds) + strcase.ToCamel(code.Description)
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}

		user := users.GetAdminUser(repo.DB)

		timecode := &m.Timecode{
			Seconds:     code.Seconds,
			VideoID:     videoID,
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
