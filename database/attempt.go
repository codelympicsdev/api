package database

import (
	"errors"
	"time"

	"github.com/lucacasonato/wrap"
	"github.com/lucacasonato/wrap/filter"
)

// Attempt is the data stored about an attempt in the database
type Attempt struct {
	ID        string `bson:"id"`
	User      string `bson:"user"`
	Challenge string `bson:"challenge"`

	Started   time.Time `bson:"started"`
	Timeout   time.Time `bson:"timeout"`
	Completed time.Time `bson:"completed"`

	Input          *AttemptInput  `bson:"input"`
	ExpectedOutput *AttemptOutput `bson:"expected_output"`
	RecievedOutput *AttemptOutput `bson:"recieved_output"`
}

// AttemptInput is the input data for a specifc attempt
type AttemptInput struct {
	Arguments []string `json:"arguments"`
	Stdin     string   `json:"stdin"`
}

// AttemptOutput is the output data for a specifc attempt
type AttemptOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

// GetAttemptByID an attempt from the database by id
func GetAttemptByID(id string) (*Attempt, error) {
	data, err := db.Collection("attempts").Document(id).Get()
	if err != nil {
		return nil, err
	}

	var attempt = new(Attempt)
	err = data.DataTo(attempt)
	if err != nil {
		return nil, err
	}

	attempt.ID = data.Document.ID

	return attempt, nil
}

// Save user data to the database
func (attempt *Attempt) Save() error {
	if attempt.ID == "" {
		resp, err := db.Collection("attempts").Add(attempt)
		if err == nil {
			attempt.ID = resp.ID
		}
		return err
	}

	return db.Collection("attempts").Document(attempt.ID).Set(attempt)
}

// GetAttemptCount gets amount of live attempts for a challenge and user
func GetAttemptCount(challenge *Challenge, userID string) (int, error) {
	iterator, err := db.Collection("attempts").Where(filter.AND(filter.Equal("challenge", challenge.ID), filter.Equal("user", userID))).Count("count").DocumentIterator()
	if err != nil {
		return 0, err
	}
	defer iterator.Close()

	if !iterator.Next() {
		return 0, nil
	}

	var data map[string]interface{}
	err = iterator.DataTo(&data)
	if err != nil {
		return 0, err
	}

	count, ok := data["count"].(int32)
	if ok == false {
		return 0, errors.New("count not int32")
	}

	return int(count), nil
}

// GetAttemptsByUser gets all attempts by a certain user
func GetAttemptsByUser(userID string, onlyForChallenge string, limit int, skip int) ([]*Attempt, error) {
	f := []filter.Filter{filter.Equal("user", userID)}

	if onlyForChallenge != "" {
		f = append(f, filter.Equal("challenge", onlyForChallenge))
	}

	iterator, err := db.Collection("attempts").Where(filter.AND(f...)).Sort(wrap.Descending("started")).Skip(skip).Limit(limit).DocumentIterator()
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var attempts []*Attempt
	for iterator.Next() {
		var attempt = new(Attempt)
		err = iterator.DataTo(attempt)
		if err != nil {
			return nil, err
		}

		attempt.ID = iterator.ID()
		attempts = append(attempts, attempt)
	}

	return attempts, nil
}
