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

	/*
			c := &database.Challenge{
				Name: "Add",
				Description: `# Add two numbers

		Input: two decimal numbers seperated via newline on stdin
		Output: result of adding the two numbers as decimal on stdout`,
				Generator: `(function() {
		  var one = Math.random() * 10;
		  var two = Math.random() * 10;
		  return {
		    Input: {
		      Args: [],
		      Stdin: one.toString() + '\n' + two.toString(),
		    },
		    Output: {
		      Stdout: (one + two).toString(),
		      Stderr: '',
		    },
		  };
		})()
		`,
				MaxLiveAttempts: 3,
				Timeout:         2 * time.Second,
				PublishDate:     time.Now(),
				ResultsDate:     time.Now().Add(10 * time.Minute),
			}

			err := c.Save()
			if err != nil {
				panic(err)
			}

			fmt.Println(c.ID)
			fmt.Println(challenge.Generate(c))
	*/

	fmt.Println("5d21e24bf789f169d5548652")

	http.ListenAndServe(":8080", r)
}
