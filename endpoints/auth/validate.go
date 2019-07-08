package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
	"github.com/gorilla/context"
)

// TokenValidationMiddleware validates the JWT
func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.TokenFromHeader(r)
		if err != nil {
			log.Println(err)
			errors.InvalidCredentials(w)
			return
		}

		if token.RequiresUpgrade == true {
			errors.UnverifiedToken(w)
			return
		}

		context.Set(r, "token", token)

		next.ServeHTTP(w, r)
	})
}

// ScopeValidationMiddleware checks this token has the required scopes
func ScopeValidationMiddleware(scopes []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := context.Get(r, "token").(*auth.Token)
			if ok == false {
				errors.InternalServerError(w)
				return
			}

			for _, scope := range scopes {
				if !token.HasScope(scope) {
					errors.MissingPermission(w, scope)
					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}

// TrustedClientIDValidationMiddleware validates if the client id is trusted and the token is valid - used for auth only
func TrustedClientIDValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		authorizationParts := strings.Split(authorizationHeader, " ")

		if len(authorizationParts) != 2 || authorizationParts[0] == "" || authorizationParts[1] == "" {
			errors.InvalidAPIClient(w)
			return
		}

		client, err := database.GetAPIClientByID(authorizationParts[0])
		if err != nil {
			log.Println(err.Error())
			errors.InvalidAPIClient(w)
			return
		}

		if !auth.DoesHashMatch(client.Secret, authorizationParts[1]) || !client.Trusted {
			errors.InvalidAPIClient(w)
			return
		}

		next.ServeHTTP(w, r)
	})
}
