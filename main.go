package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/codelympicsdev/api/endpoints/apiclients"
	"github.com/codelympicsdev/api/endpoints/attempts"
	"github.com/codelympicsdev/api/endpoints/auth"
	"github.com/codelympicsdev/api/endpoints/challenges"
	"github.com/codelympicsdev/api/endpoints/users"
)

func main() {
	r := mux.NewRouter()

	v0 := r.PathPrefix("/v0").Subrouter()

	apiclients.Route(v0.PathPrefix("/apicĺients").Subrouter())
	auth.Route(v0.PathPrefix("/auth").Subrouter())
	users.Route(v0.PathPrefix("/users").Subrouter())
	challenges.Route(v0.PathPrefix("/challenges").Subrouter())
	attempts.Route(v0.PathPrefix("/attempts").Subrouter())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, r)
}
