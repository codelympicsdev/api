package attempts

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// GetAttemptResponse is the response for getting a single attempt
type GetAttemptResponse struct {
	ID        string `json:"id"`
	User      string `json:"user"`
	Challenge string `json:"challenge"`

	Started   int64 `json:"started"`
	Timeout   int64 `json:"timeout"`
	Completed int64 `json:"completed"`

	Input          *database.AttemptInput  `json:"input"`
	ExpectedOutput *database.AttemptOutput `json:"expected_output"`
	RecievedOutput *database.AttemptOutput `json:"recieved_output"`
}

// get a challenge
func get(w http.ResponseWriter, r *http.Request) {
	token, ok := context.Get(r, "token").(*auth.Token)
	if ok == false {
		errors.InternalServerError(w)
		return
	}

	vars := mux.Vars(r)

	attempt, err := database.GetAttemptByID(vars["id"])
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			errors.NotFound(w)
			return
		}
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	if !((attempt.User == token.ID && token.HasScope("challenge.attempt.read")) || token.HasScope("admin.attempts")) {
		errors.MissingPermission(w, "challenge.attempt.read")
		return
	}

	challenge, err := database.GetChallengeByID(attempt.Challenge)
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			errors.NotFound(w)
			return
		}
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	var resp = &GetAttemptResponse{
		ID:        attempt.ID,
		User:      attempt.User,
		Challenge: attempt.Challenge,

		Started:   attempt.Started.UnixNano(),
		Timeout:   attempt.Timeout.UnixNano(),
		Completed: attempt.Completed.UnixNano(),

		Input:          attempt.Input,
		RecievedOutput: attempt.RecievedOutput,
	}

	if time.Now().After(challenge.ResultsDate) || token.HasScope("admin.attempts") {
		resp.ExpectedOutput = attempt.ExpectedOutput
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}
