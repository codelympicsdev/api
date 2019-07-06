package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"strings"
	"time"

	"github.com/codelympicsdev/api/database"
	"github.com/gbrlsnchs/jwt/v3"
)

// Issuer is the issuer of the JWT token
const Issuer = "codelympics.dev"

// Validity is how long a key is valid
const Validity = 24 * time.Hour

// PublicKey used to verify tokens
var PublicKey *rsa.PublicKey
var privateKey *rsa.PrivateKey

var secret *jwt.RSA

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("failed to generate private key ", err.Error())
	}

	PublicKey = &privateKey.PublicKey

	secret = jwt.NewRSA(jwt.SHA512, privateKey, PublicKey)
}

// Token is the structure for the JWT token
type Token struct {
	jwt.Payload

	ID        string `json:"id"`
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`

	Scopes []string `json:"scopes"`
}

// NewToken from user and scopes
func NewToken(user *database.User, scopes []string, audience string) *Token {
	now := time.Now()

	return &Token{
		Payload: jwt.Payload{
			Issuer:         Issuer,
			Subject:        user.ID,
			Audience:       jwt.Audience{audience},
			ExpirationTime: now.Add(24 * time.Hour).Unix(),
			IssuedAt:       now.Unix(),
		},

		ID:        user.ID,
		FullName:  user.FullName,
		Email:     user.Email,
		AvatarURL: user.AvatarURL,

		Scopes: scopes,
	}
}

// Sign a token and return a JWT
func (t *Token) Sign() (string, error) {
	header := jwt.Header{}

	token, err := jwt.Sign(header, t, secret)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

// Validate a token
func Validate(token string, audience string) (*Token, error) {
	raw, err := jwt.Parse([]byte(token))
	if err != nil {
		return nil, err
	}
	if err = raw.Verify(secret); err != nil {
		return nil, err
	}

	var (
		t *Token
	)

	if _, err = raw.Decode(&t); err != nil {
		return nil, err
	}

	now := time.Now()
	iatValidator := jwt.IssuedAtValidator(now)
	expValidator := jwt.ExpirationTimeValidator(now, true)
	audValidator := jwt.AudienceValidator(jwt.Audience{audience})

	if err := t.Validate(iatValidator, expValidator, audValidator); err != nil {
		return nil, err
	}

	return t, nil
}

// User gets this tokens' user from the database
func (t *Token) User() *database.User {
	return &database.User{}
}

// HasScope checks if a token has this scope
func (t *Token) HasScope(requiredScope string) bool {
	for _, scope := range t.Scopes {
		if strings.HasPrefix(requiredScope, scope) {
			return true
		}
	}

	return false
}
