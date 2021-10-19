package sqlstore

import (
	"database/sql"
	"go_practicum/internal/app/store"

	_ "github.com/lib/pq"
)

type Store struct {
	db             *sql.DB
	linkRepository *LinkRepository
}

func New(db *sql.DB) *Store {
	s := &Store{
		db: db,
	}
	//creating table 'links' if not exist. should move to another method?
	s.db.Exec("CREATE TABLE IF NOT EXISTS links (id bigserial NOT NULL PRIMARY KEY, link text NOT NULL UNIQUE, code text NOT NULL, userid text NOT NULL);")
	return s
}

func (s *Store) Link() store.LinkRepository {
	if s.linkRepository == nil {
		s.linkRepository = &LinkRepository{
			store: s,
		}
	}

	return s.linkRepository
}
