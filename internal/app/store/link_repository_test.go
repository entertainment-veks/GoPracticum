package store

import (
	"go_practicum/internal/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkRepository_Create(t *testing.T) {
	s, teardown := TestStore(t)
	defer teardown("links")

	l, err := s.Link().Create(model.TestLink())
	assert.NoError(t, err)
	assert.NotNil(t, l)
}

func TestLinkRepository_GetByCode(t *testing.T) {
	s, teardown := TestStore(t)
	defer teardown("links")

	code := "NOT_EX"
	_, err := s.Link().GetByCode(code)
	assert.Error(t, err)

	s.Link().Create(model.TestLink())
	l, err := s.Link().GetByCode(model.TestLink().Code)
	assert.NoError(t, err)
	assert.NotNil(t, l)
}
