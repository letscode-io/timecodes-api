package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khaiql/dbcleaner/engine"
	"gopkg.in/khaiql/dbcleaner.v2"
)

func createDBCleaner(t *testing.T) dbcleaner.DbCleaner {
	t.Helper()

	cleaner := dbcleaner.New()
	dsn := getEnvDSN()
	pg := engine.NewPostgresEngine(dsn.String())
	cleaner.SetEngine(pg)

	return cleaner
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
