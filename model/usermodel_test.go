package model

import (
	"testing"
)

func TestToString(t *testing.T) {
	user := User{
		FirstName:    "Test",
		LastName:     "User",
		Email:        "test@user.com",
		Organization: "sales",
	}

	want := "Test\tUser\ttest@user.com\tsales"
	if got := user.String(); got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
