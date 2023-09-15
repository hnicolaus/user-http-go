package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/lib/pq"
)

// IsUniqueConstraintViolation determines if an error is a Postgres UNIQUE constraint error.
func IsUniqueConstraintViolation(err error) bool {
	// http://godoc.org/github.com/lib/pq#Error
	switch e := err.(type) {
	case *pq.Error:
		if e.Code == "23505" {
			return true
		}
	}

	return false
}

func GenerateJWT(phoneNumber string) (string, error) {
	// Get the private key file path and passphrase from environment variables
	privateKeyFile := os.Getenv("PRIVATE_KEY_FILE")
	privateKeyPassphrase := os.Getenv("PRIVATE_KEY_PASSPHRASE")

	// Load the RSA private key
	privateKeyData, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEMWithPassword(privateKeyData, privateKeyPassphrase)
	if err != nil {
		return "", err
	}

	// Create a new token object
	token := jwt.New(jwt.SigningMethodRS256)

	// Set claims (payload data)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_phone_number"] = phoneNumber               // Subject
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix() // Token expiration time

	// Sign the token with the RSA private key
	return token.SignedString(privateKey)
}
