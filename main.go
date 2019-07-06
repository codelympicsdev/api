package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/codelympicsdev/api/endpoints/auth"
	"github.com/codelympicsdev/api/endpoints/user"
)

func main() {
	r := mux.NewRouter()

	v0 := r.PathPrefix("/v0").Subrouter()

	auth.Route(v0.PathPrefix("/auth").Subrouter())
	user.Route(v0.PathPrefix("/user").Subrouter())

	http.ListenAndServe(":8080", r)
}
