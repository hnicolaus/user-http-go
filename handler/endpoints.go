package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	var (
		context = context.Background()

		response = generated.RegisterUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	request := generated.User{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	user, errorList := convertRegisterUserRequestToUser(request)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return ctx.JSON(http.StatusBadRequest, response)
	}

	userID, err := s.Repository.InsertUser(context, user)
	if err != nil {
		if utils.IsUniqueConstraintViolation(err) {
			response.Header.Messages = []string{"phone number is already registered to an existing user"}
			return ctx.JSON(http.StatusConflict, response)
		}

		response.Header.Messages = []string{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User.Id = &userID
	return ctx.JSON(http.StatusOK, response)
}

// NOTE: Check Authenticated cmd/main.go that returns JWT token after successful login
func (s *Server) UserLogin(ctx echo.Context) error {
	var (
		context = context.Background()

		response = generated.UserLoginResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	request := generated.User{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	// Get user's phone number from request body
	validPhoneNumber, errorList := validatePhoneNumber(request.PhoneNumber)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Get user data
	user, err := s.getSingleUser(context, repository.UserFilter{PhoneNumber: validPhoneNumber})
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	// Validate password format is valid
	inputPassword, errorList := validatePassword(request.Password)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Validate input password (plain) matches user's password (hashed and salted)
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)) != nil {
		response.Header.Messages = []string{"invalid password"}
		return ctx.JSON(http.StatusBadRequest, response)
	}

	// Increment successful login count for the users
	s.Repository.IncrementSuccessfulLoginCount(context, user.ID)

	// Set data to Echo context so AuthenticatedMiddleware can generate and return JWT in the Authorization header
	ctx.Set(string(utils.JWTClaimUserID), user.ID)
	ctx.Set(string(utils.JWTClaimPermissions), []utils.JWTPermission{utils.JWTPermissionGetUser, utils.JWTPermissionUpdateUser})

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.Data.Id = &user.ID

	return ctx.JSON(http.StatusOK, response)
}

// NOTE: Check AuthenticationMiddleware cmd/main.go that authenticates the JWT token
func (s *Server) GetUser(ctx echo.Context) error {
	var (
		context = context.Background()

		response = generated.GetUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	// Authenticate and get userID of the requester
	userID, err := authenticate(ctx, utils.JWTPermissionGetUser)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return ctx.JSON(http.StatusForbidden, response)
	}

	// Get data for the userID
	user, err := s.getSingleUser(context, repository.UserFilter{UserID: userID})
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.Data = generated.User{
		FullName:    &user.FullName,
		PhoneNumber: &user.PhoneNumber,
	}

	return ctx.JSON(http.StatusOK, response)
}

// NOTE: Check AuthenticationMiddleware cmd/main.go that authenticates the JWT token
func (s *Server) UpdateUser(ctx echo.Context) error {
	var (
		context = context.Background()

		response = generated.UpdateUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	// Authenticate and get userID of the requester
	userID, err := authenticate(ctx, utils.JWTPermissionUpdateUser)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return ctx.JSON(http.StatusForbidden, response)
	}

	// Update user data
	request := generated.User{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, nil)
	}

	updateRequest, errorList := convertUpdateUserRequestToUser(userID, request)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return ctx.JSON(http.StatusBadRequest, response)
	}

	if err := s.Repository.UpdateUser(context, updateRequest); err != nil {
		if utils.IsUniqueConstraintViolation(err) {
			response.Header.Messages = []string{"phone number is already registered to an existing user"}
			return ctx.JSON(http.StatusConflict, response)
		}

		response.Header.Messages = []string{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	return ctx.JSON(http.StatusOK, response)
}
