package timecodelikes

import (
	m "timecodes/pkg/models"

	"github.com/jinzhu/gorm"
)

// TimecodeLikeRepository represents repository interface
type TimecodeLikeRepository interface {
	Create(*m.TimecodeLike, uint) (*m.TimecodeLike, error)
	Delete(*m.TimecodeLike, uint) (*m.TimecodeLike, error)
}

// DBTimecodeLikeRepository implements timecode like repository
type DBTimecodeLikeRepository struct {
	TimecodeLikeRepository

	DB *gorm.DB
}

// Create creates a new timecode like by given parameters
func (repo *DBTimecodeLikeRepository) Create(timecodeLike *m.TimecodeLike, userID uint) (*m.TimecodeLike, error) {
	timecodeLike.UserID = userID

	err := repo.DB.Create(timecodeLike).Error

	return timecodeLike, err
}

// Delete deletes timecode like
func (repo *DBTimecodeLikeRepository) Delete(timecodeLike *m.TimecodeLike, userID uint) (*m.TimecodeLike, error) {
	err := repo.DB.Where(&m.TimecodeLike{UserID: userID, TimecodeID: timecodeLike.TimecodeID}).First(timecodeLike).Error
	if err != nil {
		return nil, err
	}

	repo.DB.Unscoped().Delete(timecodeLike)

	return timecodeLike, nil
}
