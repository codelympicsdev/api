package auth

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

// Init the auth stuff
func init() {
	pemFile := os.Getenv("PRIVATE_KEY")
	var err error
	var keyFile []byte
	if strings.HasPrefix(pemFile, "-----") {
		keyFile = []byte(pemFile)
	} else {
		keyFile, err = ioutil.ReadFile(pemFile)
		if err != nil {
			log.Fatalln("failed to open private key " + err.Error())
		}
	}

	key, _ := pem.Decode(keyFile)
	privateKey, err = x509.ParsePKCS1PrivateKey(key.Bytes)
	if err != nil {
		log.Fatalln("failed to decode private key " + err.Error())
	}

	PublicKey = &privateKey.PublicKey

	secret = jwt.NewRSA(jwt.SHA512, privateKey, PublicKey)
}

// Token is the structure for the JWT token
type Token struct {
	jwt.Payload

	ID        string `json:"id,omitempty"`
	FullName  string `json:"full_name,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`

	Scopes []string `json:"scopes"`
}

// NewToken from user and scopes
func NewToken(user *database.User, client *database.APIClient, requestedScopes []string) *Token {
	now := time.Now()

	scopes := []string{}

	for _, scope := range requestedScopes {
		if (hasScope(user.Scopes, scope) || hasScope(user.Scopes, "admin")) && hasScope(client.Scopes, scope) {
			scopes = append(scopes, scope)
		}
	}

	return &Token{
		Payload: jwt.Payload{
			Issuer:         Issuer,
			Subject:        user.ID,
			Audience:       jwt.Audience{client.ID},
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
func Validate(token string) (*Token, error) {
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

	if err := t.Validate(iatValidator, expValidator); err != nil {
		return nil, err
	}

	return t, nil
}

// HasScope checks if a token has this scope
func (t *Token) HasScope(requiredScope string) bool {
	if t == nil || len(t.Scopes) == 0 {
		return false
	}
	return hasScope(t.Scopes, requiredScope)
}

func hasScope(scopes []string, requiredScope string) bool {
	for _, scope := range scopes {
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

	token, err := Validate(authorizationParts[1])
	if err != nil {
		return nil, err
	}

	return token, nil
}
