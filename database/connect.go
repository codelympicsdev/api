package database

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"os"
	"time"

	"github.com/lucacasonato/wrap"
)

var db *wrap.Database

func init() {
	var databaseURI = os.Getenv("MONGO_URI")

	client, err := wrap.Connect(databaseURI, 5*time.Second)
	if err != nil {
		log.Fatalln(err.Error())
	}

	db = client.Database("production")
}

func randomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
