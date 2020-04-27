package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func initDB() {
	var err error

	dsn := url.URL{
		User:     url.UserPassword(os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD")),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", os.Getenv("PG_HOST"), os.Getenv("PG_PORT")),
		Path:     os.Getenv("PG_DB"),
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}

	db, err = gorm.Open("postgres", dsn.String())
	if err != nil {
		log.Println("Failed to connect to database")
		panic(err)
	}

	log.Println("DB connection has been established.")
}

func createTables() {
	if db.HasTable(&Annotation{}) {
		return
	}

	err := db.CreateTable(&Annotation{})
	if err != nil {
		log.Println("Table already exists")
	}
}

func runMigrations() {
	db.AutoMigrate(&Annotation{})
	db.Model(&Annotation{}).AddUniqueIndex("idx_seconds_text_video_id", "seconds", "text", "video_id")
}
