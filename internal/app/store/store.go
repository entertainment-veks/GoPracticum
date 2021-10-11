package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

const POSTGRES_KEY string = "postgres"

type Store struct {
	db             *sql.DB
	linkRepository *LinkRepository
}

func New(db *sql.DB) *Store {
	s := &Store{
		db: db,
	}
	//creating table 'links' if not exist. should move to another method?
	s.db.Exec("CREATE TABLE links (id bigserial NOT NULL PRIMARY KEY, link text NOT NULL,code text NOT NULL);")
	return s
}

// func (s *Store) Open() error {
// 	db, err := sql.Open(POSTGRES_KEY, s.config.DatabseURL)
// 	if err != nil {
// 		return err
// 	}

// 	if err := db.Ping(); err != nil {
// 		return err
// 	}

// 	s.db = db

// 	//creating table 'links' if not exist. should move to another method?
// 	db.Exec("CREATE TABLE links (id bigserial NOT NULL PRIMARY KEY, link text NOT NULL,code text NOT NULL);")

// 	return nil
// }

// func (s *Store) Close() {
// 	s.db.Close()
// }

func (s *Store) Link() *LinkRepository {
	if s.linkRepository == nil {
		s.linkRepository = &LinkRepository{
			store: s,
		}
	}

	return s.linkRepository
}
