package main

import (
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func init() {
	initDB()
	createTables()
	runMigrations()
}

func main() {
	startHttpServer()
}
