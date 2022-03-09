package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"regexp"
)

type Link struct {
	ID     int
	Link   string
	Code   string
	UserID string
}

func (l *Link) Validate() error {
	regex, err := regexp.Compile("^(http).*")
	if err != nil {
		return err
	}

	return validation.ValidateStruct(
		l,
		validation.Field(&l.Link, validation.Required, is.URL, validation.Match(regex)),
		validation.Field(&l.Code, validation.Required),
		validation.Field(&l.UserID, validation.Required),
	)
}
