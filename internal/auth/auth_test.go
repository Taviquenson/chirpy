package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

// Boot.dev's example of how to properly test
func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password",
			password:      password1,
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Incorrect password",
			password:      "wrongPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match different hash",
			password:      password1,
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash",
			password:      password1,
			hash:          "invalidhash",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("CheckPasswordHash() expects %v, got %v", tt.matchPassword, match)
			}
		})
	}
}

// My tests
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

func TestJWT(t *testing.T) {
	id := uuid.New()
	tokenString, err := MakeJWT(id, "AllYourBase", time.Second*5)
	if err != nil {
		t.Errorf(`%v`, err)
	}

	returnedID, err := ValidateJWT(tokenString, "AllYourBase")
	if err != nil {
		t.Errorf(`%v`, err)
	}
	if returnedID != id {
		t.Errorf("The IDs don't match")
	}
}
