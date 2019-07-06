package database

import "time"

// Challenge is the data stored about a challenge in the database
type Challenge struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Generator   string     `json:"generator"`
	PublishDate time.Time `json:"publish_date"`
	ResultsDate time.Time `json:"results_date"`
}

// GetChallengeByID a challenge from the database by id
func GetChallengeByID(id string) (*Challenge, error) {
	data, err := db.Collection("challenges").Document(id).Get()
	if err != nil {
		return nil, err
	}

	var challenge = new(Challenge)
	err = data.DataTo(challenge)
	if err != nil {
		return nil, err
	}

	challenge.ID = data.Document.ID

	return challenge, nil
}

// Save user data to the database
func (challenge *Challenge) Save() error {
	if challenge.ID == "" {
		resp, err := db.Collection("challenges").Add(challenge)
		if err == nil {
			challenge.ID = resp.ID
		}
		return err
	}

	return db.Collection("users").Document(challenge.ID).Set(challenge)
}
