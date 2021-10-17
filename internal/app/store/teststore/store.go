package teststore

import (
	"go_practicum/internal/app/model"
	"go_practicum/internal/app/store"
)

type Store struct {
	linkRepository *LinkRepository
}

func New() *Store {
	return &Store{}
}

func (s *Store) Link() store.LinkRepository {
	if s.linkRepository == nil {
		s.linkRepository = &LinkRepository{
			store: s,
			links: make(map[string]*model.Link),
		}
	}

	return s.linkRepository
}
