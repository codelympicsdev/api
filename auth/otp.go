package auth

import (
	"github.com/codelympicsdev/api/database"
	"github.com/pquerna/otp/totp"
)

// GenerateOTPSecret creates an OTP secret for the user
func GenerateOTPSecret(user *database.User) (string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "codelympics.dev",
		AccountName: user.Email,
	})
	if err != nil {
		return "", err
	}

	return key.Secret(), nil
}

// IsOTPValid checks if OTP is valid
func IsOTPValid(user *database.User, otp string) bool {
	return totp.Validate(otp, user.OTPSecret)
}
