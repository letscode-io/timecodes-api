package main

import "github.com/jinzhu/gorm"

// TimecodeLike represents timecode like model
type TimecodeLike struct {
	gorm.Model
	TimecodeID uint `json:"timecodeId" gorm:"not null"`
	UserID     uint `json:"userId" gorm:"not null"`
}

// TimecodeLikeRepository represents repository interface
type TimecodeLikeRepository interface {
	Create(*TimecodeLike, uint) (*TimecodeLike, error)
	Delete(*TimecodeLike, uint) (*TimecodeLike, error)
}

// DBTimecodeLikeRepository implements timecode like repository
type DBTimecodeLikeRepository struct {
	TimecodeLikeRepository

	DB *gorm.DB
}

// Create creates a new timecode like by given parameters
func (repo *DBTimecodeLikeRepository) Create(timecodeLike *TimecodeLike, userID uint) (*TimecodeLike, error) {
	timecodeLike.UserID = userID

	err := repo.DB.Create(timecodeLike).Error

	return timecodeLike, err
}

// Delete deletes timecode like
func (repo *DBTimecodeLikeRepository) Delete(timecodeLike *TimecodeLike, userID uint) (*TimecodeLike, error) {
	err := repo.DB.Where(&TimecodeLike{UserID: userID, TimecodeID: timecodeLike.TimecodeID}).First(timecodeLike).Error
	if err != nil {
		return nil, err
	}

	repo.DB.Unscoped().Delete(timecodeLike)

	return timecodeLike, nil
}
