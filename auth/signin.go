package auth

import (
	"errors"

	"github.com/codelympicsdev/api/database"
)

// Signin an existing user
func Signin(email string, password string) (*database.User, error) {
	user, err := database.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !DoesPasswordMatch(user, password) {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}
