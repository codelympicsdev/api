package graphql

import (
	"log"
	"strings"
	"time"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/challenge"
	"github.com/codelympicsdev/api/database"
	gql "github.com/graphql-go/graphql"
)

func generateAttemptResolve(p gql.ResolveParams) (interface{}, error) {
	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	challengeID := p.Args["challenge"].(string)
	live := p.Args["live"].(bool)

	if live && !token.HasScope("challenge.attempt.write") {
		return nil, ErrMissingPermissions
	}

	c, err := database.GetChallengeByID(challengeID)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	if time.Now().Before(c.PublishDate) && !token.HasScope("admin.challenges") {
		return nil, ErrNotFound
	}

	if live {
		count, err := database.GetAttemptCount(c, token.ID)
		if err != nil {
			log.Println(err.Error())
			return nil, ErrInternalServer
		}

		if c.MaxLiveAttempts <= count {
			return nil, ErrMaxAttemptsReached
		}
	}

	input, expectedOutput, err := challenge.Generate(c)
	if err != nil {
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	now := time.Now()

	a := &database.Attempt{
		User:      token.Subject,
		Challenge: c.ID,

		Started: now,
		Timeout: now.Add(c.Timeout),

		Input:          input,
		ExpectedOutput: expectedOutput,
	}

	if live {
		a.Save()
	}

	resp := attempt{
		ID:   a.ID,
		User: a.User,
		Challenge: Challenge{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			PublishDate: c.PublishDate,
			ResultsDate: c.ResultsDate,
		},

		Started: a.Started,
		Timeout: a.Timeout,

		Input: &attemptInput{a.Input.Arguments, a.Input.Stdin},
	}

	if !live {
		resp.ExpectedOutput = &attemptOutput{a.ExpectedOutput.Stdout, a.ExpectedOutput.Stderr}
	}

	return resp, nil
}

func submitAttemptResolve(p gql.ResolveParams) (interface{}, error) {
	now := time.Now()

	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	attemptID := p.Args["attempt"].(string)
	stdout := p.Args["stdout"].(string)
	stderr := p.Args["stderr"].(string)

	a, err := database.GetAttemptByID(attemptID)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			return nil, ErrMissingPermissions
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	if a.User != token.ID {
		return nil, ErrInvalidAuth
	}

	if !token.HasScope("challenge.attempt.write") {
		return nil, ErrMissingPermissions
	}

	if a.Completed != (time.Time{}) || a.RecievedOutput != nil {
		return nil, ErrDuplicate
	}

	if a.Timeout.Before(now) {
		return nil, ErrTimedOut
	}

	a.Completed = now
	a.RecievedOutput = &database.AttemptOutput{Stdout: stdout, Stderr: stderr}

	err = a.Save()
	if err != nil {
		return nil, ErrInternalServer
	}

	c, err := database.GetChallengeByID(a.Challenge)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	var resp = &attempt{
		ID:   a.ID,
		User: a.User,
		Challenge: Challenge{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			PublishDate: c.PublishDate,
			ResultsDate: c.ResultsDate,
		},

		Started:   a.Started,
		Timeout:   a.Timeout,
		Completed: a.Completed,

		Input:          &attemptInput{a.Input.Arguments, a.Input.Stdin},
		RecievedOutput: &attemptOutput{a.RecievedOutput.Stdout, a.RecievedOutput.Stderr},
	}

	if time.Now().After(c.ResultsDate) || token.HasScope("admin.attempts") {
		resp.ExpectedOutput = &attemptOutput{a.ExpectedOutput.Stdout, a.ExpectedOutput.Stderr}
	}

	return resp, nil
}
