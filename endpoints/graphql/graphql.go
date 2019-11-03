package graphql

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/codelympicsdev/api/auth"
	gql "github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// Handler is the thing that handles the gql endpoint
var Handler *handler.Handler

type key int

const (
	keyToken key = iota
)

var (
	ErrInvalidAuth        = errors.New("invalid authentication token")
	ErrMissingPermissions = errors.New("missing permissions for operation")
	ErrInternalServer     = errors.New("internal server error")
	ErrNotFound           = errors.New("not found")
	ErrMaxAttemptsReached = errors.New("max attempts reached")
	ErrDuplicate          = errors.New("duplicate submission")
	ErrTimedOut           = errors.New("timed out")
)

// Handle a gql connection with authorization
func Handle(w http.ResponseWriter, r *http.Request) {
	token, err := auth.TokenFromHeader(r)
	if err != nil {
		token = nil
	}

	Handler.ContextHandler(context.WithValue(r.Context(), keyToken, token), w, r)
}

func init() {
	Handler =
		handler.New(&handler.Config{
			Schema:     schema(),
			Pretty:     true,
			Playground: true,
		})
}

func schema() *gql.Schema {
	query := gql.Fields{
		"me": &gql.Field{
			Type:    userType,
			Resolve: meResolve,
		},
		"user": &gql.Field{
			Type: userType,
			Args: gql.FieldConfigArgument{
				"id": &gql.ArgumentConfig{
					Type: gql.ID,
				},
			},
			Resolve: userResolve,
		},
		"challenge": &gql.Field{
			Type: challengeType,
			Args: gql.FieldConfigArgument{
				"id": &gql.ArgumentConfig{
					Type: gql.ID,
				},
			},
			Resolve: challengeResolve,
		},
		"challenges": &gql.Field{
			Type: gql.NewList(challengeType),
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
			Resolve: challengesResolve,
		},
		"attempt": &gql.Field{
			Type: attemptType,
			Args: gql.FieldConfigArgument{
				"id": &gql.ArgumentConfig{
					Type: gql.ID,
				},
			},
			Resolve: attemptResolve,
		},
	}
	mutation := gql.Fields{
		"generateAttempt": &gql.Field{
			Type: attemptType,
			Args: gql.FieldConfigArgument{
				"challenge": &gql.ArgumentConfig{
					Type: gql.NewNonNull(gql.ID),
				},
				"live": &gql.ArgumentConfig{
					Type:         gql.Boolean,
					DefaultValue: false,
				},
			},
			Resolve: generateAttemptResolve,
		},
		"submitAttempt": &gql.Field{
			Type: attemptType,
			Args: gql.FieldConfigArgument{
				"attempt": &gql.ArgumentConfig{
					Type: gql.ID,
				},
				"stdout": &gql.ArgumentConfig{
					Type: gql.String,
				},
				"stderr": &gql.ArgumentConfig{
					Type: gql.String,
				},
			},
			Resolve: submitAttemptResolve,
		},
	}
	rootQuery := gql.ObjectConfig{Name: "RootQuery", Fields: query}
	rootMutation := gql.ObjectConfig{Name: "RootMutation", Fields: mutation}
	schemaConfig := gql.SchemaConfig{Query: gql.NewObject(rootQuery), Mutation: gql.NewObject(rootMutation)}
	schema, err := gql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}
	return &schema
}
