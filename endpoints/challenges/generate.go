package challenges

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/challenge"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

// GenerateTestResponse is the response for a test challenge response attempt
type GenerateTestResponse struct {
	Challenge      string                  `json:"challenge"`
	Timeout        int64                   `json:"timeout"`
	Input          *database.AttemptInput  `json:"input"`
	ExpectedOutput *database.AttemptOutput `json:"expected_output"`
}

// GenerateLiveResponse is the response for a live challenge response attempt
type GenerateLiveResponse struct {
	ID        string                 `json:"id"`
	Challenge string                 `json:"challenge"`
	RespondBy int64                  `json:"respond_by"`
	Input     *database.AttemptInput `json:"input"`
}

// get a challenge
func generateTest(w http.ResponseWriter, r *http.Request) {
	token, ok := context.Get(r, "token").(*auth.Token)
	if ok == false {
		errors.InternalServerError(w)
		return
	}

	vars := mux.Vars(r)

	c, err := database.GetChallengeByID(vars["id"])
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			errors.NotFound(w)
			return
		}
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	if time.Now().Before(c.PublishDate) && !token.HasScope("admin.challenges") {
		errors.NotFound(w)
		return
	}

	input, expectedOutput, err := challenge.Generate(c)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	resp := &GenerateTestResponse{
		Challenge:      c.ID,
		Timeout:        int64(c.Timeout.Seconds() * 1000),
		Input:          input,
		ExpectedOutput: expectedOutput,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}

// get a challenge
func generateLive(w http.ResponseWriter, r *http.Request) {
	token, ok := context.Get(r, "token").(*auth.Token)
	if ok == false {
		errors.InternalServerError(w)
		return
	}

	if !token.HasScope("challenge.attempt") {
		errors.MissingPermission(w, "challenge.attempt")
		return
	}

	vars := mux.Vars(r)

	c, err := database.GetChallengeByID(vars["id"])
	if err != nil {
		if strings.Contains(err.Error(), "no documents") {
			errors.NotFound(w)
			return
		}
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	if time.Now().Before(c.PublishDate) && !token.HasScope("admin.challenges") {
		errors.NotFound(w)
		return
	}

	count, err := database.GetAttemptCount(c, token.ID)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	if c.MaxLiveAttempts <= count {
		errors.MaxAttemptsReached(w)
		return
	}

	input, expectedOutput, err := challenge.Generate(c)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	now := time.Now()

	attempt := &database.Attempt{
		User:      token.Subject,
		Challenge: c.ID,

		Started: now,
		Timeout: now.Add(c.Timeout),

		Input:          input,
		ExpectedOutput: expectedOutput,
	}

	attempt.Save()

	resp := &GenerateLiveResponse{
		ID:        attempt.ID,
		Challenge: c.ID,
		RespondBy: attempt.Timeout.UnixNano(),
		Input:     input,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}
