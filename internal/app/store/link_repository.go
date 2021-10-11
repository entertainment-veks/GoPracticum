package store

import "go_practicum/internal/app/model"

type LinkRepository struct {
	store *Store
}

func (r *LinkRepository) Create(l *model.Link) (*model.Link, error) {
	if err := l.Validate(); err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(
		"INSERT INTO links (link, code) VALUES ($1, $2) RETURNING id",
		l.Link,
		l.Code,
	).Scan(&l.ID); err != nil {
		return nil, err
	}

	return l, nil
}

func (r *LinkRepository) GetByCode(c string) (*model.Link, error) {
	l := &model.Link{}
	if err := r.store.db.QueryRow(
		"SELECT id, link, code FROM links WHERE code = $1",
		c,
	).Scan(
		&l.ID,
		&l.Link,
		&l.Code,
	); err != nil {
		return nil, err
	}

	return l, nil
}
