package sqlstore

import (
	"database/sql"
	"go_practicum/internal/app/model"
	"go_practicum/internal/app/store"
)

type LinkRepository struct {
	store *Store
}

func (r *LinkRepository) Create(l *model.Link) error {
	if err := l.Validate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO links (link, code, userid) VALUES ($1, $2, $3) RETURNING id",
		l.Link,
		l.Code,
		l.UserID,
	).Scan(&l.ID)
}

func (r *LinkRepository) GetByCode(c string) (*model.Link, error) {
	l := &model.Link{}
	if err := r.store.db.QueryRow(
		"SELECT id, link, code, userid FROM links WHERE code = $1",
		c,
	).Scan(
		&l.ID,
		&l.Link,
		&l.Code,
		&l.UserID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return l, nil
}
