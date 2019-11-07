package graphql

import (
	"log"
	"strings"
	"time"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	gql "github.com/graphql-go/graphql"
)

type Challenge struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxAttempts int       `json:"max_attempts"`
	PublishDate time.Time `json:"publish_date"`
	ResultsDate time.Time `json:"results_date"`
}

var challengeType = gql.NewObject(
	gql.ObjectConfig{
		Name: "Challenge",
		Fields: gql.Fields{
			"id": &gql.Field{
				Type: gql.ID,
			},
			"name": &gql.Field{
				Type: gql.String,
			},
			"description": &gql.Field{
				Type: gql.String,
			},
			"max_attempts": &gql.Field{
				Type: gql.Int,
			},
			"publish_date": &gql.Field{
				Type: gql.DateTime,
			},
			"results_date": &gql.Field{
				Type: gql.DateTime,
			},
		},
	},
)

func init() {
	challengeType.AddFieldConfig("attempts", &gql.Field{
		Type: gql.NewList(attemptType),
		Args: gql.FieldConfigArgument{
			"limit": &gql.ArgumentConfig{
				Type:         gql.Int,
				DefaultValue: 10,
			},
			"skip": &gql.ArgumentConfig{
				Type:         gql.Int,
				DefaultValue: 0,
			},
		},
		Resolve: func(p gql.ResolveParams) (interface{}, error) {
			var token *auth.Token
			var ok bool
			if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok || token == nil {
				return nil, ErrInvalidAuth
			}
			p.Args["id"] = token.Subject
			p.Args["challenge"] = p.Source.(Challenge).ID
			return attemptsByUserResolve(p)
		},
	})
}

func challengeResolve(p gql.ResolveParams) (interface{}, error) {
	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	c, err := database.GetChallengeByID(p.Args["id"].(string))
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	if time.Now().Before(c.PublishDate) && !token.HasScope("admin.challenges") {
		return nil, ErrMissingPermissions
	}

	resp := Challenge{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		MaxAttempts: c.MaxLiveAttempts,
		PublishDate: c.PublishDate,
		ResultsDate: c.ResultsDate,
	}

	return resp, nil
}

func challengesResolve(p gql.ResolveParams) (interface{}, error) {
	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	challenges, err := database.GetChallenges(!token.HasScope("admin.challenges"), p.Args["limit"].(int), p.Args["skip"].(int))
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	var resp []Challenge
	for _, c := range challenges {
		resp = append(resp, Challenge{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			PublishDate: c.PublishDate,
			ResultsDate: c.ResultsDate,
		})
	}

	return resp, nil
}
