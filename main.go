package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/codelympicsdev/api/endpoints/apiclients"
	"github.com/codelympicsdev/api/endpoints/auth"
	"github.com/codelympicsdev/api/endpoints/graphql"
)

func main() {
	r := mux.NewRouter()

	v0 := r.PathPrefix("/v0").Subrouter()

	//v0.Use(HTTPRedirect)

	apiclients.Route(v0.PathPrefix("/apiclients").Subrouter())
	auth.Route(v0.PathPrefix("/auth").Subrouter())

	v0.HandleFunc("/graphql", graphql.Handle)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, cors.New(cors.Options{AllowedHeaders: []string{"Authorization", "Content-Type"}}).Handler(r))
	if err != nil {
		panic(err.Error())
	}
}

// HTTPRedirect redirects http traffic to https
func HTTPRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["X-Forwarded-Proto"][0] == "http" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusMovedPermanently)
			return
		}
		next.ServeHTTP(w, r)
	})
}
