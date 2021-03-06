package main

import (
	"os"
	googleAPI "timecodes/cmd/google_api"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Email      string
	GoogleID   string
	PictureURL string
}

type UserRepository interface {
	FindOrCreateByGoogleInfo(*googleAPI.UserInfo) *User
}

type DBUserRepository struct {
	UserRepository

	DB *gorm.DB
}

func (repo *DBUserRepository) FindOrCreateByGoogleInfo(userInfo *googleAPI.UserInfo) *User {
	user := &User{}

	repo.DB.Where(User{GoogleID: userInfo.ID}).
		Assign(User{Email: userInfo.Email, PictureURL: userInfo.Picture}).
		FirstOrCreate(&user)

	return user
}

func getAdminUser(db *gorm.DB) *User {
	adminUser := &User{Email: os.Getenv("ADMIN_EMAIL"), GoogleID: os.Getenv("ADMIN_GOOGLE_ID")}
	db.FirstOrCreate(adminUser)

	return adminUser
}
