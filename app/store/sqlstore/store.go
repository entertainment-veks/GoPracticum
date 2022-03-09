package sqlstore

import (
	"database/sql"
	_ "github.com/lib/pq"
	"go_practicum/app/store"
)

type Store struct {
	db             *sql.DB
	linkRepository *LinkRepository
}

func New(db *sql.DB) *Store {
	s := &Store{
		db: db,
	}

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
