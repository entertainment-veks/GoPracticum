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

func (r *LinkRepository) GetAllByUserID(id string) ([]*model.Link, error) {
	links := make([]*model.Link, 100)

	rows, err := r.store.db.Query(
		"SELECT id, link, code, userid FROM links WHERE userid = $1 LIMIT 100",
		id,
	)
	if err != nil {
		return nil, err
	}

	for i := 0; rows.Next(); i++ {
		l := &model.Link{}
		err := rows.Scan(
			&l.ID,
			&l.Link,
			&l.Code,
			&l.UserID,
		)
		if err != nil {
			return nil, err
		}

		links[i] = l
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if links[0] == nil {
		return nil, store.ErrRecordNotFound
	}

	return links, nil
}
