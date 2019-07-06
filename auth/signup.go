package auth

import (
	gravatar "github.com/Automattic/go-gravatar"
	"github.com/codelympicsdev/api/database"
)

// Signup a new user
func Signup(name string, email string, password string) (*database.User, error) {
	g := gravatar.NewGravatarFromEmail(email)
	url := g.GetURL()

	pw, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &database.User{
		FullName:  name,
		Email:     email,
		AvatarURL: url,

		Password: pw,
	}

	err = user.Save()
	if err != nil {
		return nil, err
	}

	return user, nil
}
