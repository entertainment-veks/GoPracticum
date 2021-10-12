package teststore

import (
	"go_practicum/internal/app/model"
	"go_practicum/internal/app/store"
)

type LinkRepository struct {
	store *Store
	links map[string]*model.Link
}

func (r *LinkRepository) Create(l *model.Link) error {
	if err := l.Validate(); err != nil {
		return err
	}

	r.links[l.Code] = l
	l.ID = len(r.links)

	return nil
}

func (r *LinkRepository) GetByCode(c string) (*model.Link, error) {
	l, ok := r.links[c]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return l, nil
}
