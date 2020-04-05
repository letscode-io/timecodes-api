package main

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func main() {
	initDB()
	createTables()
	runMigrations()
	handleRequests()
}
