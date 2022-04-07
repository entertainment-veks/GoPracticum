package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"regexp"
)

type Link struct {
	ID     int
	Link   string
	Code   string
	UserID string
}

func (l *Link) Validate() error {
	//lint:ignore S1007 errors in regex if change " -to> `
	regex, err := regexp.Compile("^(http(s)?://)[\\w.-]+(?:\\.[\\w.-]+)+[\\w\\-._~:/?#[\\]@!$&'()*+,;=]+$")
	if err != nil {
		return err
	}

	return validation.ValidateStruct(
		l,
		validation.Field(&l.Link, validation.Required, validation.Match(regex)),
		validation.Field(&l.Code, validation.Required),
		validation.Field(&l.UserID, validation.Required),
	)
}
