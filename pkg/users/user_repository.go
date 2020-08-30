package users

import (
	"os"

	googleAPI "timecodes/pkg/google_api"
	"timecodes/pkg/models"

	"github.com/jinzhu/gorm"
)

// UserRepository represents an interface for user repository
type UserRepository interface {
	FindOrCreateByGoogleInfo(*googleAPI.UserInfo) *models.User
}

// DBUserRepository represents database repository
type DBUserRepository struct {
	UserRepository

	DB *gorm.DB
}

// FindOrCreateByGoogleInfo finds user by given google information or creates a new user if it doesn't exist
func (repo *DBUserRepository) FindOrCreateByGoogleInfo(userInfo *googleAPI.UserInfo) *models.User {
	user := &models.User{}

	repo.DB.Where(models.User{GoogleID: userInfo.ID}).
		Assign(models.User{Email: userInfo.Email, PictureURL: userInfo.Picture}).
		FirstOrCreate(&user)

	return user
}

// GetAdminUser finds or creates an admin user
func GetAdminUser(db *gorm.DB) *models.User {
	adminUser := &models.User{Email: os.Getenv("ADMIN_EMAIL"), GoogleID: os.Getenv("ADMIN_GOOGLE_ID")}
	db.FirstOrCreate(adminUser)

	return adminUser
}
