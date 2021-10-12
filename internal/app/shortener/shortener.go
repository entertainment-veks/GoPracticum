package shortener

import (
	"database/sql"
	"go_practicum/internal/app/store/sqlstore"
	"net/http"
)

const DRIVER_NAME = "postgres"

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	s := newServer(store, config.BaseURL)

	return http.ListenAndServe(config.ServerAddress, s)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open(DRIVER_NAME, databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
