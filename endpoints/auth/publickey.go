package auth

import (
	"crypto/x509"
	"encoding/pem"
	"net/http"

	"github.com/codelympicsdev/api/auth"
)

// publickey for jwt signature check
func publickey(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.Write(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(auth.PublicKey),
	}))
}
