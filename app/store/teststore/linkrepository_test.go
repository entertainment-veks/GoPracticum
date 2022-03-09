package teststore_test

import (
	"go_practicum/app/model"
	"go_practicum/app/store"
	"go_practicum/app/store/teststore"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkRepository_Create(t *testing.T) {
	s := teststore.New()
	l := model.TestLink()
	assert.NoError(t, s.Link().Create(l))
	assert.NotNil(t, l)
}

func TestLinkRepository_CreateAll(t *testing.T) {
	s := teststore.New()

	l1 := model.TestLink()
	l2 := model.TestLink()
	l3 := model.TestLink()

	l2.Code = "ABC"
	l3.Code = "123"

	ls := []*model.Link{l1, l2, l3}

	assert.NoError(t, s.Link().CreateAll(ls))
	assert.NotNil(t, ls)
}

func TestLinkRepository_GetByCode(t *testing.T) {
	s := teststore.New()
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

func TestLinkRepository_GetAllByUseID(t *testing.T) {
	s := teststore.New()

	userID := "example"
	_, err := s.Link().GetAllByUserID(userID)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	l := model.TestLink()
	l.UserID = userID
	s.Link().Create(l)
	ls, err := s.Link().GetAllByUserID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, ls)
}

func TestLinkRepository_DeleteAllByCode(t *testing.T) {
	s := teststore.New()
	l := model.TestLink()
	assert.NoError(t, s.Link().Create(l))
	assert.NotNil(t, l)

	assert.NoError(t, s.Link().DeleteAllByCode([]string{l.Code}))

	l, _ = s.Link().GetByCode(l.Code)
	assert.Nil(t, l)
}
