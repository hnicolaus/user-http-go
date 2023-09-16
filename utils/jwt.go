package utils

import (
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
)

const (
	jwtExpiryDuration = time.Minute * 30 // Token expires in 30 minutes by default
)

func GenerateJWT(userID int64, permissions []JWTPermission) (string, error) {
	privateKeyData, err := os.ReadFile("../rsa")
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return "", err
	}

	claims := CustomClaims{
		UserID:      userID,
		Permissions: permissions,
		ExpiresAt:   time.Now().Add(jwtExpiryDuration).Unix(),
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// Sign the token with the RSA private key
	return token.SignedString(privateKey)
}

func AuthenticateJWT(ctx echo.Context) error {
	authHeader := ctx.Request().Header.Get("Authorization")
	if authHeader == "" {
		return ctx.String(http.StatusUnauthorized, "missing authorization header")
	}

	if len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return ctx.String(http.StatusUnauthorized, "token format is invalid")
	}

	token := authHeader[7:]

	publicKeyData, err := os.ReadFile("../rsa.pub")
	if err != nil {
		return err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return err
	}

	// Parse & validate token using publicKey
	tok, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Use publicKey for token verification
		return publicKey, nil
	})
	if err != nil {
		return err
	}

	// Check if token is valid (has not expired)
	if claims, ok := tok.Claims.(*CustomClaims); ok && tok.Valid {
		if time.Now().Unix() >= claims.ExpiresAt {
			return errors.New("JWT has expired")
		}

		// Set custom claims to context so handler can use the values, i.e. authorization
		ctx.Set(string(JWTClaimUserID), claims.UserID)
		ctx.Set(string(JWTClaimPermissions), claims.Permissions)

		return nil
	}

	return nil
}
