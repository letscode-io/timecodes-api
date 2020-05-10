package main

import (
	"fmt"
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
	DB *gorm.DB
}

func (repo *DBUserRepository) FindOrCreateByGoogleInfo(userInfo *googleAPI.UserInfo) *User {
	user := &User{}

	err := repo.DB.Where(User{GoogleID: userInfo.ID}).
		Assign(User{Email: userInfo.Email, PictureURL: userInfo.Picture}).
		FirstOrCreate(&user).Error
	if err != nil {
		fmt.Println("ERR")
		fmt.Println(err)
	}
	fmt.Println(user)
	return user
}
