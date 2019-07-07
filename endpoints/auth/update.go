package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codelympicsdev/api/auth"
	"github.com/codelympicsdev/api/database"
	"github.com/codelympicsdev/api/endpoints/errors"
	"github.com/gorilla/context"
)

// UpdatePasswordRequest is the request for updating a password
type UpdatePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// EnableOTPRequest is the request for enabling otp
type EnableOTPRequest struct {
	CurrentPassword string `json:"current_password"`
}

// EnableOTPResponse is the response for enabling otp
type EnableOTPResponse struct {
	URL string `json:"url"`
}

// VerifyOTPRequest is the request for verifying otp
type VerifyOTPRequest struct {
	OTP string `json:"otp"`
}

// update the password
func updatePassword(w http.ResponseWriter, r *http.Request) {
	token, ok := context.Get(r, "token").(*auth.Token)
	if ok == false {
		errors.InternalServerError(w)
		return
	}

	if token.HasScope("auth") {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			errors.WrongContentType(w)
			return
		}

		var req UpdatePasswordRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errors.MalformedBody(w)
			return
		}

		if req.CurrentPassword == "" {
			errors.MissingField(w, "current_password")
			return
		}

		if req.NewPassword == "" {
			errors.MissingField(w, "new_password")
			return
		}

		u, err := database.GetUserByID(token.ID)
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		if !auth.DoesPasswordMatch(u, req.CurrentPassword) {
			errors.InvalidCredentials(w)
			return
		}

		pw, err := auth.HashPassword(req.NewPassword)
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		u.Password = pw

		err = u.Save()
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	} else {
		errors.MissingPermission(w, "auth")
		return
	}
}

// enable otp for a user
func enableOTP(w http.ResponseWriter, r *http.Request) {
	token, ok := context.Get(r, "token").(*auth.Token)
	if ok == false {
		errors.InternalServerError(w)
		return
	}

	if token.HasScope("auth") {
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			errors.WrongContentType(w)
			return
		}

		var req EnableOTPRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errors.MalformedBody(w)
			return
		}

		if req.CurrentPassword == "" {
			errors.MissingField(w, "current_password")
			return
		}

		u, err := database.GetUserByID(token.ID)
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		if !auth.DoesPasswordMatch(u, req.CurrentPassword) {
			errors.InvalidCredentials(w)
			return
		}

		key, err := auth.GenerateOTP(u)
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		u.OTPEnabled = false
		u.OTPSecret = key.Secret()

		err = u.Save()
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		var resp = EnableOTPResponse{
			URL: key.URL(),
		}
		w.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}
	} else {
		errors.MissingPermission(w, "auth")
		return
	}
}

// verify otp after enabling
func verifyOTP(w http.ResponseWriter, r *http.Request) {
	token, ok := context.Get(r, "token").(*auth.Token)
	if ok == false {
		errors.InternalServerError(w)
		return
	}

	if token.HasScope("auth") {
		u, err := database.GetUserByID(token.ID)
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		if u.OTPEnabled {
			errors.OTPAlreadyEnabled(w)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			errors.WrongContentType(w)
			return
		}

		var req VerifyOTPRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			errors.MalformedBody(w)
			return
		}

		if req.OTP == "" {
			errors.MissingField(w, "otp")
			return
		}

		if !auth.IsOTPValid(u, req.OTP) {
			errors.InvalidCredentials(w)
			return
		}
		u.OTPEnabled = true

		err = u.Save()
		if err != nil {
			log.Println(err.Error())
			errors.InternalServerError(w)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	} else {
		errors.MissingPermission(w, "auth")
		return
	}
}
