package auth

import (
	"github.com/codelympicsdev/api/database"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword takes a password and hashes it
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	return string(hash), err
}

// DoesPasswordMatch checks if the hashed stored password and the supplied plain text password match
func DoesPasswordMatch(user *database.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		return true
	}

	return false
}

