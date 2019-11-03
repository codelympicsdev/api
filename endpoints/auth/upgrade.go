package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
)

// OPTUpgradeRequest is what is used to upgrade a token with a OTP
type OPTUpgradeRequest struct {
	Token string `json:"token"`
	OTP   string `json:"otp"`
}

func upgradeWithOTP(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		errors.WrongContentType(w)
		return
	}

	var req OPTUpgradeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errors.MalformedBody(w)
		return
	}

	if req.OTP == "" {
		errors.MissingField(w, "otp")
		return
	}

	if req.Token == "" {
		errors.MissingField(w, "token")
		return
	}

	token, err := auth.Validate(req.Token)
	if err != nil {
		errors.InvalidCredentials(w)
		return
	}

	if !token.RequiresUpgrade {
		errors.TokenAlreadyUpgraded(w)
		return
	}

	user, err := database.GetUserByID(token.Subject)
	if err != nil {
		errors.UserDoesntExist(w)
		return
	}

	if !auth.IsOTPValid(user, req.OTP) {
		errors.InvalidCredentials(w)
		return
	}

	token.Upgrade(user)

	t, err := token.Sign()
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}

	var resp = AuthResponse{
		Token: t,
	}
	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Println(err.Error())
		errors.InternalServerError(w)
		return
	}
}
