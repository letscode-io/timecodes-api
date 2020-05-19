package main

import (
	"context"
	"net/http"
	"net/http/httptest"
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

func executeRequest(t *testing.T, router http.Handler, req *http.Request, user *User) *httptest.ResponseRecorder {
	t.Helper()

	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(context.Background(), CurrentUserKey{}, user)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	return rr
}
