package teststore

import (
	"go_practicum/app/model"
	"go_practicum/app/store"
	"sync"
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
			mu:    &sync.Mutex{},
		}
	}

	return s.linkRepository
}
