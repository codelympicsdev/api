package database

import (
	"errors"
	"time"

	"github.com/lucacasonato/wrap/filter"
)

// Attempt is the data stored about an attempt in the database
type Attempt struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Challenge string `json:"challenge"`

	Started   time.Time `json:"started"`
	Timeout   time.Time `json:"timeout"`
	Completed time.Time `json:"completed"`

	Input          *AttemptInput  `json:"input"`
	ExpectedOutput *AttemptOutput `json:"expected_output"`
	RecievedOutput *AttemptOutput `json:"recieved_output"`
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
