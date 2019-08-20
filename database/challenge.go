package database

import (
	"time"

	"github.com/lucacasonato/wrap/filter"
)

// Challenge is the data stored about a challenge in the database
type Challenge struct {
	ID              string        `bson:"id"`
	Name            string        `bson:"name"`
	Description     string        `bson:"description"`
	Generator       string        `bson:"generator"`
	MaxLiveAttempts int           `bson:"max_live_attempts"`
	Timeout         time.Duration `bson:"timeout"`
	PublishDate     time.Time     `bson:"publish_date"`
	ResultsDate     time.Time     `bson:"results_date"`
}

// GetChallenges all challenges from the database
func GetChallenges(onlyPublished bool, limit int, skip int) ([]*Challenge, error) {
	c := db.Collection("challenges")

	query := c.All()

	if onlyPublished {
		query = c.Where(filter.GreaterThanOrEqual("publish_date", time.Now())).Skip(skip).Limit(limit)
	}

	data, err := query.DocumentIterator()
	if err != nil {
		return nil, err
	}

	var challenges []*Challenge

	for data.Next() {
		var challenge = new(Challenge)
		err = data.DataTo(challenge)
		if err != nil {
			return nil, err
		}

		challenge.ID = data.ID()

		challenges = append(challenges, challenge)
	}

	return challenges, nil
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
