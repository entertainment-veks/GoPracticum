package store

import (
	"go_practicum/app/model"
)

type LinkRepository interface {
	Create(l *model.Link) error
	CreateAll([]*model.Link) error
	GetByCode(c string) (*model.Link, error)
	GetAllByUserID(id string) ([]*model.Link, error)
	DeleteAllByCode([]string) error
}
