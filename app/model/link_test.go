package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLink_Validate(t *testing.T) {
	testCases := []struct {
		name    string
		l       func() *Link
		isValid bool
	}{
		{
			name: "Valid",
			l: func() *Link {
				return TestLink()
			},
			isValid: true,
		},
		{
			name: "Empty link",
			l: func() *Link {
				l := TestLink()
				l.Link = ""
				return l
			},
			isValid: false,
		},
		{
			name: "Invalid link #1",
			l: func() *Link {
				l := TestLink()
				l.Link = "www.goo gle.com"
				return l
			},
			isValid: false,
		},
		{
			name: "Invalid link #2",
			l: func() *Link {
				l := TestLink()
				l.Link = "invalid"
				return l
			},
			isValid: false,
		},
		{
			name: "Empty code",
			l: func() *Link {
				l := TestLink()
				l.Code = ""
				return l
			},
			isValid: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isValid {
				assert.NoError(t, tc.l().Validate())
			} else {
				assert.Error(t, tc.l().Validate())
			}
		})
	}
}
