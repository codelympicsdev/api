package database

import (
	"errors"

	"github.com/lucacasonato/wrap/filter"
)

// User is the data stored about a single user in the database
type User struct {
	ID         string   `bson:"id"`
	FullName   string   `bson:"full_name"`
	Email      string   `bson:"email"`
	AvatarURL  string   `bson:"avatar_url"`
	OTPEnabled bool     `bson:"otp_enabled"`
	Password   string   `bson:"password"`
	OTPSecret  string   `bson:"otp_secret"`
	Scopes     []string `bson:"scopes"`
}

// GetUserByID a user from the database by id
func GetUserByID(id string) (*User, error) {
	data, err := db.Collection("users").Document(id).Get()
	if err != nil {
		return nil, err
	}

	var user = new(User)
	err = data.DataTo(user)
	if err != nil {
		return nil, err
	}

	user.ID = data.Document.ID

	return user, nil
}

// GetUserByEmail a user from the database by email
func GetUserByEmail(email string) (*User, error) {
	iterator, err := db.Collection("users").Where(filter.Equal("email", email)).DocumentIterator()
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	if !iterator.Next() {
		return nil, errors.New("invalid username or password")
	}

	var user = new(User)
	err = iterator.DataTo(user)
	if err != nil {
		return nil, err
	}

	user.ID = iterator.ID()

	return user, nil
}

// Save user data to the database
func (user *User) Save() error {
	if user.ID == "" {
		resp, err := db.Collection("users").Add(user)
		if err == nil {
			user.ID = resp.ID
		}
		return err
	}

	return db.Collection("users").Document(user.ID).Set(user)
}
