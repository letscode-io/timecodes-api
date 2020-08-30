package testhelpers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"timecodes/internal/db"

	"github.com/gorilla/mux"
	"github.com/khaiql/dbcleaner"
	"github.com/khaiql/dbcleaner/engine"
)

// CreateDBCleaner create an instance of db cleaner
func CreateDBCleaner(t *testing.T, database *db.Database) dbcleaner.DbCleaner {
	t.Helper()

	cleaner := dbcleaner.New()
	dsn := database.DSN
	pg := engine.NewPostgresEngine(dsn.String())

	cleaner.SetEngine(pg)

	return cleaner
}

type iHandler interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// ExecuteRequest helper
func ExecuteRequest(ctx context.Context, t *testing.T, req *http.Request, handler iHandler, path string) *httptest.ResponseRecorder {
	t.Helper()

	req.Header.Set("Content-Type", "application/json")

	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	if path == "" {
		handler.ServeHTTP(rr, req)
	} else {
		rtr := mux.NewRouter()
		rtr.Handle(path, handler)
		rtr.ServeHTTP(rr, req)
	}

	return rr
}
