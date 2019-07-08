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

// SubmitRequest is what is used to sign in
type SubmitRequest struct {
	Output *database.AttemptOutput `json:"output"`
}

// submit a response for a solve attempt
func submit(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

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

	if attempt.User != token.ID {
		errors.InvalidCredentials(w)
		return
	}

	if !token.HasScope("challenge.attempt.write") {
		errors.MissingPermission(w, "challenge.attempt.write")
		return
	}

	if attempt.Completed != (time.Time{}) || attempt.RecievedOutput != nil {
		errors.Duplicate(w)
		return
	}

	if attempt.Timeout.Before(now) {
		errors.AttemptTimedOut(w)
		return
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errors.WrongContentType(w)
		return
	}

	var req SubmitRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.MalformedBody(w)
		return
	}

	if req.Output == nil {
		errors.MissingField(w, "output")
		return
	}

	attempt.Completed = time.Now()
	attempt.RecievedOutput = req.Output

	err = attempt.Save()
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	w.Write([]byte(`{}`))
}
