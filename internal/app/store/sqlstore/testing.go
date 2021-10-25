package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

const (
	driverName         = "postgres"
	testingDatabaseURL = "host=localhost dbname=shortener_db sslmode=disable user=postgres password=postgres"
)

func TestDB(t *testing.T) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open(driverName, testingDatabaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}

		db.Close()
	}
}
