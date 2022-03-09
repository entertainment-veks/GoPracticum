package sqlstore

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate"
	_ "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	"strings"
	"testing"
)

const (
	driverName         = "postgres"
	testingDatabaseURL = "postgres://postgres:postgres@localhost:5432/shortener?sslmode=disable"
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

	m, err := migrate.New("file://./migration", testingDatabaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", ")))
		}

		db.Close()
	}
}
