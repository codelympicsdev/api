package auth

import (
	"github.com/codelympicsdev/api/database"
	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateOTP creates an OTP secret for the user
func GenerateOTP(user *database.User) (*otp.Key, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "codelympics.dev",
		AccountName: user.Email,
	})
	if err != nil {
		return nil, err
	}

	return key, nil
}

// IsOTPValid checks if OTP is valid
func IsOTPValid(user *database.User, otp string) bool {
	return totp.Validate(otp, user.OTPSecret)
}
