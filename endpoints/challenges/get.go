package challenges

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

// GetChallengeResponse is the response for a single challenge
type GetChallengeResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	PublishDate int64  `json:"publish_date"`
	ResultsDate int64  `json:"results_date"`
}

// get a challenge
func get(w http.ResponseWriter, r *http.Request) {
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

	resp := &GetChallengeResponse{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		PublishDate: c.PublishDate.Unix(),
		ResultsDate: c.ResultsDate.Unix(),
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}
