package structs

import (
	"strings"
	"testing"

	validator "gopkg.in/go-playground/validator.v9"
)

// The profile update writes traits straight to the Kratos table, so the API must hold the
// identity schema's name limits (GitHub #65).
func TestReqUpdateUserNameLimits(t *testing.T) {
	v := validator.New()
	cases := []struct {
		first, last string
		ok          bool
	}{
		{strings.Repeat("a", 12), strings.Repeat("b", 50), true},
		{strings.Repeat("a", 13), "Last", false},
		{"First", strings.Repeat("b", 51), false},
		{"", "Last", false},
	}
	for _, c := range cases {
		err := v.Struct(ReqUpdateUser{FirstName: c.first, LastName: c.last})
		if (err == nil) != c.ok {
			t.Errorf("first=%d last=%d chars: got err=%v, want ok=%v", len(c.first), len(c.last), err, c.ok)
		}
	}
}
