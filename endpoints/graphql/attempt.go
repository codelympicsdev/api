package graphql

import (
	"log"
	"strings"
	"time"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	gql "github.com/graphql-go/graphql"
)

type attempt struct {
	ID        string    `json:"id"`
	User      string    `json:"user"`
	Challenge Challenge `json:"challenge"`

	CreationDate   time.Time `json:"creation_date"`
	TimeoutDate    time.Time `json:"timeout_date"`
	SubmissionDate time.Time `json:"submission_date"`

	Input          *attemptInput  `json:"input"`
	ExpectedOutput *attemptOutput `json:"expected_output"`
	RecievedOutput *attemptOutput `json:"recieved_output"`
}

type attemptInput struct {
	Arguments []string `json:"arguments"`
	Stdin     string   `json:"stdin"`
}

type attemptOutput struct {
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
}

var attemptType = gql.NewObject(
	gql.ObjectConfig{
		Name: "Attempt",
		Fields: gql.Fields{
			"id": &gql.Field{
				Type: gql.ID,
			},
			"user": &gql.Field{
				Type: userType,
				Resolve: func(p gql.ResolveParams) (interface{}, error) {
					p.Args["id"] = p.Source.(*attempt).User
					return userResolve(p)
				},
			},
			"challenge": &gql.Field{
				Type: challengeType,
			},
			"creation_date": &gql.Field{
				Type: gql.DateTime,
			},
			"timeout_date": &gql.Field{
				Type: gql.DateTime,
			},
			"submission_date": &gql.Field{
				Type: gql.DateTime,
			},
			"input": &gql.Field{
				Type: attemptInputType,
			},
			"expected_output": &gql.Field{
				Type: attemptOutputType,
			},
			"recieved_output": &gql.Field{
				Type: attemptOutputType,
			},
		},
	},
)

var attemptInputType = gql.NewObject(
	gql.ObjectConfig{
		Name: "AttemptInput",
		Fields: gql.Fields{
			"arguments": &gql.Field{
				Type: gql.NewList(gql.String),
			},
			"stdin": &gql.Field{
				Type: gql.String,
			},
		},
	},
)

var attemptOutputType = gql.NewObject(
	gql.ObjectConfig{
		Name: "AttemptOutput",
		Fields: gql.Fields{
			"stdout": &gql.Field{
				Type: gql.String,
			},
			"stderr": &gql.Field{
				Type: gql.String,
			},
		},
	},
)

func attemptResolve(p gql.ResolveParams) (interface{}, error) {
	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	a, err := database.GetAttemptByID(p.Args["id"].(string))
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	if !((a.User == token.ID && token.HasScope("challenge.attempt.read")) || token.HasScope("admin.attempts")) {
		return nil, ErrMissingPermissions
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

		CreationDate:   a.CreationDate,
		TimeoutDate:    a.TimeoutDate,
		SubmissionDate: a.SubmissionDate,

		Input: &attemptInput{a.Input.Arguments, a.Input.Stdin},
	}

	if a.RecievedOutput != nil {
		resp.RecievedOutput = &attemptOutput{a.RecievedOutput.Stdout, a.RecievedOutput.Stderr}
	}

	if time.Now().After(c.ResultsDate) || token.HasScope("admin.attempts") {
		resp.ExpectedOutput = &attemptOutput{a.ExpectedOutput.Stdout, a.ExpectedOutput.Stderr}
	}

	return resp, nil
}

func attemptsByUserResolve(p gql.ResolveParams) (interface{}, error) {
	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	user := p.Args["id"].(string)

	if !((user == token.ID && token.HasScope("challenge.attempt.read")) || token.HasScope("admin.attempts")) {
		return nil, ErrMissingPermissions
	}

	var challengeID string
	if challengeID, ok = p.Args["challenge"].(string); !ok {
		challengeID = ""
	}
	attempts, err := database.GetAttemptsByUser(user, challengeID, p.Args["limit"].(int), p.Args["skip"].(int))
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	var resp []*attempt
	for _, a := range attempts {
		c, err := database.GetChallengeByID(a.Challenge)
		if err != nil {
			if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
				return nil, ErrNotFound
			}
			log.Println(err.Error())
			return nil, ErrInternalServer
		}

		var r = &attempt{
			ID:   a.ID,
			User: a.User,
			Challenge: Challenge{
				ID:          c.ID,
				Name:        c.Name,
				Description: c.Description,
				PublishDate: c.PublishDate,
				ResultsDate: c.ResultsDate,
			},

			CreationDate:   a.CreationDate,
			TimeoutDate:    a.TimeoutDate,
			SubmissionDate: a.SubmissionDate,

			Input: &attemptInput{a.Input.Arguments, a.Input.Stdin},
		}

		if a.RecievedOutput != nil {
			r.RecievedOutput = &attemptOutput{a.RecievedOutput.Stdout, a.RecievedOutput.Stderr}
		}

		if time.Now().After(c.ResultsDate) || token.HasScope("admin.attempts") {
			r.ExpectedOutput = &attemptOutput{a.ExpectedOutput.Stdout, a.ExpectedOutput.Stderr}
		}

		resp = append(resp, r)
	}

	return resp, nil
}
