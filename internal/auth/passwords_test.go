package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("techie$4Good")
	if err != nil {
		t.Errorf(`%v`, err)
	}
	match, err := CheckPasswordHash("techie$4Good", hash)
	if !match || err != nil {
		t.Errorf(`%v`, err)
	}
}

func TestHashPasswordDefault(t *testing.T) {
	hash, err := HashPassword("unset")
	if err != nil {
		t.Errorf(`%v`, err)
	}
	match, err := CheckPasswordHash("unset", hash)
	if !match || err != nil {
		t.Errorf(`%v`, err)
	}
}
