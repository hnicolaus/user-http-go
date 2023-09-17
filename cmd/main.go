package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/dgrijalva/jwt-go"

	"github.com/labstack/echo/v4"
)

const (
	jwtExpiryDuration = time.Minute * 30 // Token expires in 30 minutes by default
)

func main() {
	e := echo.New()
	e.Pre(AuthenticationMiddleware) // register pre-handler middleware
	e.Use(AuthenticatedMiddleware)  // register post-handler middleware

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
	}
	return handler.NewServer(opts)
}

func AuthenticationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		whitelistedEndpoints := map[string]bool{
			"GET - /v1/user": true,
			"PUT - /v1/user": true,
		}

		endpoint := fmt.Sprintf("%s - %s", ctx.Request().Method, ctx.Request().URL.Path)
		if whitelistedEndpoints[endpoint] {
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
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			// Parse & validate token using publicKey
			tok, err := jwt.ParseWithClaims(token, &utils.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				// Use publicKey for token verification
				return publicKey, nil
			})
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, err.Error())
			}

			// Check if token is valid (has not expired)
			if claims, ok := tok.Claims.(*utils.CustomClaims); ok && tok.Valid {
				if time.Now().Unix() >= claims.ExpiresAt {
					return ctx.JSON(http.StatusInternalServerError, errors.New("JWT has expired"))
				}

				// Set custom claims to context so handler can use the values, i.e. authorization
				ctx.Set(string(utils.JWTClaimUserID), claims.UserID)
				ctx.Set(string(utils.JWTClaimPermissions), claims.Permissions)
			}
		}

		return next(ctx)
	}
}

func AuthenticatedMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		whitelistedEndpoints := map[string]bool{
			"POST - /v1/user/login": true,
		}

		endpoint := fmt.Sprintf("%s - %s", ctx.Request().Method, ctx.Request().URL.Path)
		if whitelistedEndpoints[endpoint] {
			// use `Before` hook so middleware can write token to the response header right before handler writes to response body
			ctx.Response().Before(func() {
				privateKeyData, err := os.ReadFile("../rsa")
				if err != nil {
					return
				}

				privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
				if err != nil {
					return
				}

				// Authorize JWT
				permissions, ok := ctx.Get(string(utils.JWTClaimPermissions)).([]utils.JWTPermission)
				if !ok {
					return
				}

				userID, ok := ctx.Get(string(utils.JWTClaimUserID)).(int64)
				if !ok {
					return
				}

				claims := utils.CustomClaims{
					UserID:      userID,
					Permissions: permissions,
					ExpiresAt:   time.Now().Add(jwtExpiryDuration).Unix(),
				}

				token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
				jwtToken, err := token.SignedString(privateKey) // sign the token with the RSA private key
				if err != nil {
					return
				}

				bearerToken := fmt.Sprintf("Bearer %s", jwtToken)

				if ctx.Response().Status == http.StatusOK {
					ctx.Response().Header().Set(echo.HeaderAuthorization, bearerToken)
				}
			})
		}

		return next(ctx)
	}
}
