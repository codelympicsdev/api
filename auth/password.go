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
	return DoesHashMatch(user.Password, password)
}

// DoesHashMatch checks if the hash matches the password
func DoesHashMatch(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}

	return false
}
