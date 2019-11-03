package graphql

import (
	"log"
	"strings"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	gql "github.com/graphql-go/graphql"
)

type user struct {
	ID         string `json:"id"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	AvatarURL  string `json:"avatar_url"`
	OTPEnabled bool   `json:"otp_enabled"`
}

var userType = gql.NewObject(
	gql.ObjectConfig{
		Name: "User",
		Fields: gql.Fields{
			"id": &gql.Field{
				Type: gql.ID,
			},
			"full_name": &gql.Field{
				Type: gql.String,
			},
			"email": &gql.Field{
				Type: gql.String,
			},
			"avatar_url": &gql.Field{
				Type: gql.String,
			},
			"otp_enabled": &gql.Field{
				Type: gql.Boolean,
			},
		},
	},
)

func init() {
	userType.AddFieldConfig("attempts", &gql.Field{
		Type: gql.NewList(attemptType),
		Args: gql.FieldConfigArgument{
			"challenge": &gql.ArgumentConfig{
				Type:         gql.ID,
				DefaultValue: nil,
			},
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
			p.Args["id"] = p.Source.(*user).ID
			return attemptsByUserResolve(p)
		},
	})
}

func meResolve(p gql.ResolveParams) (interface{}, error) {
	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	if !(token.HasScope("user.basic") || token.HasScope("user.read") || token.HasScope("admin")) {
		return nil, ErrMissingPermissions
	}

	a, err := database.GetUserByID(token.ID)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	var resp = &user{
		ID:       a.ID,
		FullName: a.FullName,
		Email:    a.Email,

		AvatarURL:  a.AvatarURL,
		OTPEnabled: a.OTPEnabled,
	}

	return resp, nil
}

func userResolve(p gql.ResolveParams) (interface{}, error) {
	var token *auth.Token
	var ok bool
	if token, ok = p.Context.Value(keyToken).(*auth.Token); !ok {
		return nil, ErrInvalidAuth
	}

	if !((p.Args["id"] == token.ID && (token.HasScope("user.basic") || token.HasScope("user.read"))) || token.HasScope("admin.users")) {
		return nil, ErrMissingPermissions
	}

	a, err := database.GetUserByID(p.Args["id"].(string))
	if err != nil {
		if strings.Contains(err.Error(), "no documents") || strings.Contains(err.Error(), "valid ObjectID") {
			return nil, ErrNotFound
		}
		log.Println(err.Error())
		return nil, ErrInternalServer
	}

	var resp = &user{
		ID:       a.ID,
		FullName: a.FullName,
		Email:    a.Email,

		AvatarURL:  a.AvatarURL,
		OTPEnabled: a.OTPEnabled,
	}

	return resp, nil
}
