package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/codelympicsdev/api/endpoints/auth"
	"github.com/codelympicsdev/api/endpoints/challenges"
	"github.com/codelympicsdev/api/endpoints/users"
)

func main() {
	r := mux.NewRouter()

	v0 := r.PathPrefix("/v0").Subrouter()

	auth.Route(v0.PathPrefix("/auth").Subrouter())
	users.Route(v0.PathPrefix("/users").Subrouter())
	challenges.Route(v0.PathPrefix("/challenges").Subrouter())

	fmt.Println("5d20e10c52cdddde3c4c21a5")

	http.ListenAndServe(":8080", r)
}
