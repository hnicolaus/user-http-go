package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const (
	jwtExpiryDuration = time.Minute * 30 // Token expires in 30 minutes
)

func GenerateJWT(phoneNumber string, permissions []JWTPermission) (string, error) {
	privateKeyData, err := os.ReadFile("../rsa")
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	claims := CustomClaims{
		PhoneNumber: phoneNumber,
		Permissions: permissions,
		ExpiresAt:   time.Now().Add(jwtExpiryDuration).Unix(),
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the RSA private key
	return token.SignedString(privateKey)
}
