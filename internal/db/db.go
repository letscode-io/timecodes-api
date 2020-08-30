package db

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Database holds database configuration
type Database struct {
	Connection *gorm.DB
	Name       string
	DSN        url.URL
}

// Init initializes a new database configuration
func Init() *Database {
	database := &Database{
		Connection: initDB(),
		Name:       getDbName(),
		DSN:        getEnvDSN(),
	}

	return database
}

func getDbName() string {
	return fmt.Sprintf("timecodes_%s", os.Getenv("APP_ENV"))
}

func getEnvDSN() url.URL {
	return getDsn(getDbName())
}

func initDB() *gorm.DB {
	envDsn := getEnvDSN()

	db, err := gorm.Open("postgres", envDsn.String())
	if err != nil {
		dbName := getDbName()
		db = createDatabase(dbName)
	}

	log.Println("Connection to the database has been established.")

	return db
}

func createDatabase(dbName string) *gorm.DB {
	db := getDefaultConnection()

	err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", dbName)).Error
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

func getDefaultConnection() *gorm.DB {
	defaultDsn := getDsn("postgres")
	db, err := gorm.Open("postgres", defaultDsn.String())
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
