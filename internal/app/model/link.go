package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type Link struct {
	ID     int
	Link   string
	Code   string
	UserID string
}

func (l *Link) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Link, validation.Required, is.URL),
		validation.Field(&l.Code, validation.Required),
		validation.Field(&l.UserID, validation.Required),
	)
}
