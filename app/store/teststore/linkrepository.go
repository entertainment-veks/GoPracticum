package teststore

import (
	"go_practicum/app/model"
	"go_practicum/app/store"
	"sync"
)

type LinkRepository struct {
	store *Store
	links map[string]*model.Link
	mu    *sync.Mutex
}

func (r *LinkRepository) Create(l *model.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.links[l.Code] = l
	l.ID = len(r.links)

	return nil
}

func (r *LinkRepository) CreateAll(ls []*model.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, l := range ls {
		r.links[l.Code] = l
		l.ID = len(r.links)
	}
	return nil
}

func (r *LinkRepository) GetByCode(c string) (*model.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	l, ok := r.links[c]
	if !ok {
		return nil, store.ErrRecordNotFound
	}

	return l, nil
}

func (r *LinkRepository) GetAllByUserID(id string) ([]*model.Link, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var result []*model.Link
	for _, current := range r.links {
		if current.UserID == id {
			result = append(result, current)
		}
	}

	if len(result) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return result, nil
}

func (r *LinkRepository) DeleteAllByCode(codes []string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, currentCode := range codes {
		delete(r.links, currentCode)
	}
	return nil
}
