package shortener

import (
	"database/sql"
	"github.com/golang-migrate/migrate"
	"go_practicum/app/config"
	"go_practicum/app/store/sqlstore"
	"net/http"
)

const driverName = "postgres"

func Start(config config.Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	s := newServer(store, config)

	return http.ListenAndServe(config.ServerAddress, s)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open(driverName, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	m, err := migrate.New("file://app/store/sqlstore/migration", databaseURL)
	if err != nil {
		return nil, err
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, err
	}

	return db, nil
}
