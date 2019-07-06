package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/endpoints/errors"
	"github.com/gorilla/context"
)

// TokenValidationMiddleware validates the JWT
func TokenValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")

		authorizationParts := strings.Split(authorizationHeader, " ")

		if len(authorizationParts) != 2 || authorizationParts[0] != "Bearer" || authorizationParts[1] == "" {
			errors.InvalidCredentials(w)
			return
		}

		token, err := auth.Validate(authorizationParts[1], auth.Issuer)
		if err != nil {
			log.Println(err)
			errors.InvalidCredentials(w)
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
