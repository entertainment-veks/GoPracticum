package sqlstore_test

import (
	"go_practicum/internal/app/model"
	"go_practicum/internal/app/store"
	"go_practicum/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t)
	defer teardown("links")
	s := sqlstore.New(db)

	l := model.TestLink()
	assert.NoError(t, s.Link().Create(l))
	assert.NotNil(t, l)
}

func TestLinkRepository_GetByCode(t *testing.T) {
	db, teardown := sqlstore.TestDB(t)
	defer teardown("links")
	s := sqlstore.New(db)

	code := "example"
	_, err := s.Link().GetByCode(code)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	l := model.TestLink()
	l.Code = code
	s.Link().Create(l)
	l, err = s.Link().GetByCode(code)
	assert.NoError(t, err)
	assert.NotNil(t, l)
}
