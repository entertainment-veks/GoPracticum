package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

const DRIVER_NAME = "postgres"
const TESTING_DATABASE_URL = "host=localhost dbname=shortener_db sslmode=disable user=postgres password=postgres"

func TestDB(t *testing.T) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open(DRIVER_NAME, TESTING_DATABASE_URL)
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
