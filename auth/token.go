package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
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
	keyFile, err := ioutil.ReadFile("test.pem")
	if err != nil {
		log.Fatal("failed to open private key ", err.Error())
	}

	key, _ := pem.Decode(keyFile)
	privateKey, err = x509.ParsePKCS1PrivateKey(key.Bytes)
	if err != nil {
		log.Fatal("failed to decode private key ", err.Error())
	}

	PublicKey = &privateKey.PublicKey

	secret = jwt.NewRSA(jwt.SHA512, privateKey, PublicKey)

	
}

// Token is the structure for the JWT token
type Token struct {
	jwt.Payload

	RequiresUpgrade bool `json:"requires_upgrade"`

	ID        string `json:"id,omitempty"`
	FullName  string `json:"full_name,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`

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

// NewUnverifiedToken is a token used when the authentication is not complete
func NewUnverifiedToken(user *database.User, scopes []string, audience string) *Token {
	now := time.Now()

	return &Token{
		Payload: jwt.Payload{
			Issuer:         Issuer,
			Subject:        user.ID,
			Audience:       jwt.Audience{audience},
			ExpirationTime: now.Add(24 * time.Hour).Unix(),
			IssuedAt:       now.Unix(),
		},

		RequiresUpgrade: true,

		Scopes: scopes,
	}
}

// Upgrade an unverified token to a verified one
func (t *Token) Upgrade(user *database.User) {
	t.RequiresUpgrade = false

	t.ID = user.ID
	t.FullName = user.FullName
	t.Email = user.Email
	t.AvatarURL = user.AvatarURL
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

// HasScope checks if a token has this scope
func (t *Token) HasScope(requiredScope string) bool {
	for _, scope := range t.Scopes {
		if strings.HasPrefix(requiredScope, scope) {
			return true
		}
	}

	return false
}

// ErrHeaderStructureIncorrect means an incorrectly formatter error
var ErrHeaderStructureIncorrect = errors.New("header structure incorrect")

// TokenFromHeader gets a token from the HTTP Authorization header
func TokenFromHeader(r *http.Request) (*Token, error) {
	authorizationHeader := r.Header.Get("Authorization")

	authorizationParts := strings.Split(authorizationHeader, " ")

	if len(authorizationParts) != 2 || authorizationParts[0] != "Bearer" || authorizationParts[1] == "" {
		return nil, ErrHeaderStructureIncorrect
	}

	token, err := Validate(authorizationParts[1], Issuer)
	if err != nil {
		return nil, err
	}

	return token, nil
}
