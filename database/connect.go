package database

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"github.com/lucacasonato/wrap"
)

var databaseURI = os.Getenv("MONGO_URI")

var db *wrap.Database

// Connect to the db
func Connect() error {
	client, err := wrap.Connect(databaseURI, 5*time.Second)
	if err != nil {
		return err
	}

	db = client.Database("production")

	return nil
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
