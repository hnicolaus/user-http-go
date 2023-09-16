package utils

import "github.com/dgrijalva/jwt-go"

type JWTPermission string

const (
	ProfileGet    JWTPermission = "get_profile"
	ProfileUpdate JWTPermission = "update_profile"
)

// CustomClaims represents the claims you want to include in your JWT.
type CustomClaims struct {
	UserID      int64           `json:"user_id"`
	PhoneNumber string          `json:"phone_number"`
	Permissions []JWTPermission `json:"permissions"`
	ExpiresAt   int64           `json:"exp"`
	jwt.StandardClaims
}
