package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"sync"

	"github.com/DATA-DOG/go-txdb"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const TXDB_NAME = "txdb_postgres"
const TXDB_CONNECTION_POOL_NAME = "timecodes_connection_pool"

var doOnce sync.Once

func initDB() *gorm.DB {
	envDsn := getDsn(os.Getenv("PG_DB"))

	db, err := gorm.Open("postgres", envDsn.String())
	if err != nil {
		db = createDatabase(os.Getenv("PG_DB"))
	}

	return db
}

func createDatabase(dbName string) *gorm.DB {
	defaultDsn := getDsn("postgres")
	db, err := gorm.Open("postgres", defaultDsn.String())
	if err != nil {
		handleDBConnectionError(err)
	}

	err = db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)).Error
	if err != nil {
		handleDBConnectionError(err)
	}

	dsn := getDsn(dbName)

	db, err = gorm.Open("postgres", dsn.String())
	if err != nil {
		handleDBConnectionError(err)
	}

	return db
}

func handleDBConnectionError(err error) {
	log.Println("Unable to connect to db")
	panic(err)
}

func getDsn(path string) url.URL {
	return url.URL{
		User:     url.UserPassword(os.Getenv("PG_USER"), os.Getenv("PG_PASSWORD")),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%s", os.Getenv("PG_HOST"), os.Getenv("PG_PORT")),
		Path:     path,
		RawQuery: (&url.Values{"sslmode": []string{"disable"}}).Encode(),
	}
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&Timecode{})
	db.Model(&Timecode{}).AddUniqueIndex(
		"idx_timecodes_seconds_text_video_id",
		"seconds", "description", "video_id",
	)
	db.AutoMigrate(&User{})
	db.AutoMigrate(&TimecodeLike{})
	db.Model(&TimecodeLike{}).AddUniqueIndex(
		"idx_timecodes_likes_user_id_timecode_id_video_id",
		"user_id", "timecode_id",
	)
}

func setupTestDB() *gorm.DB {
	initDB()
	dsn := getDsn(os.Getenv("PG_DB"))
	doOnce.Do(func() {
		txdb.Register(TXDB_NAME, "postgres", dsn.String())
	})

	db, err := gorm.Open("postgres", TXDB_NAME, TXDB_CONNECTION_POOL_NAME)
	if err != nil {
		handleDBConnectionError(err)
	}

	runMigrations(db)

	return db
}
