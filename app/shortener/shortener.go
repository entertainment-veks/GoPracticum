package shortener

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate"
	"go_practicum/app/config"
	mos "go_practicum/app/os"
	"go_practicum/app/store/sqlstore"
	"golang.org/x/crypto/acme/autocert"
	"net/http"
)

const driverName = "postgres"

func Start(ctx context.Context, config config.Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	r := newRouter(store, config)

	server := &http.Server{
		Addr:    config.ServerAddress,
		Handler: r,
	}

	go mos.ListenShutdownSignals(ctx, server)

	if config.EnableHTTPS {
		return server.Serve(autocert.NewListener(config.ServerAddress))
	} else {
		return server.ListenAndServe()
	}
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
