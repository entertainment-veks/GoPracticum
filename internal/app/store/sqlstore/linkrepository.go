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
	err := r.store.db.QueryRow(
		"INSERT INTO links (link, code, userid) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id",
		l.Link,
		l.Code,
		l.UserID,
	).Scan(&l.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return store.ErrConflict
		}
		return err
	}
	return nil
}

func (r LinkRepository) CreateAll(ls []*model.Link) error {
	var err error
	for _, l := range ls {
		err = r.store.db.QueryRow(
			"INSERT INTO links (link, code, userid) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING RETURNING id",
			l.Link,
			l.Code,
			l.UserID,
		).Scan(&l.ID)
	}
	return err
}

func (r *LinkRepository) GetByCode(c string) (*model.Link, error) {
	l := &model.Link{}
	var isDeleted *bool
	if err := r.store.db.QueryRow(
		"SELECT id, link, code, userid, deleted FROM links WHERE code = $1",
		c,
	).Scan(
		&l.ID,
		&l.Link,
		&l.Code,
		&l.UserID,
		&isDeleted,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	if *isDeleted {
		return nil, store.ErrURLDeleted
	}

	return l, nil
}

func (r *LinkRepository) GetAllByUserID(id string) ([]*model.Link, error) {
	links := []*model.Link{}

	rows, err := r.store.db.Query(
		"SELECT id, link, code, userid FROM links WHERE userid = $1",
		id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
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
		links = append(links, l)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(links) == 0 {
		return nil, store.ErrRecordNotFound
	}

	return links, nil
}

func (r *LinkRepository) DeleteAllByCode(codes []string) error {
	_, err := r.store.db.Exec("begin;")
	if err != nil {
		return err
	}

	for _, currentCode := range codes {
		_, err = r.store.db.Exec(
			"UPDATE links SET deleted = true WHERE code = $1",
			currentCode,
		)
	}
	if err != nil {
		return err
	}

	_, err = r.store.db.Exec("end;")

	return err
}
