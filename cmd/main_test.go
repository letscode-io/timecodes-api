package main

import (
	"os"
	"testing"

	"github.com/khaiql/dbcleaner/engine"
	"gopkg.in/khaiql/dbcleaner.v2"
)

var Cleaner = dbcleaner.New()
var TestDB = initDB()

func TestMain(m *testing.M) {
	dsn := getEnvDSN()
	pg := engine.NewPostgresEngine(dsn.String())
	Cleaner.SetEngine(pg)

	runMigrations(TestDB)
	defer TestDB.Close()

	os.Exit(m.Run())
}
