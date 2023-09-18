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

var (
	//define function wrappers so we can inject dummy function in UT
	fnConvertRegisterUserRequestToUser func(generated.User) (repository.User, []string)        = convertRegisterUserRequestToUser
	fnConvertUpdateUserRequestToUser   func(int64, generated.User) (repository.User, []string) = convertUpdateUserRequestToUser
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	return ctx.JSON(s.registerUser(ctx))
}
func (s *Server) registerUser(ctx echo.Context) (int, generated.RegisterUserResponse) {
	var (
		context = context.Background()

		response = generated.RegisterUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	request := generated.User{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusBadRequest, response
	}

	user, errorList := fnConvertRegisterUserRequestToUser(request)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	userID, err := s.Repository.InsertUser(context, user)
	if err != nil {
		if utils.IsUniqueConstraintViolation(err) {
			response.Header.Messages = []string{duplicatePhoneNumberErrorMsg}
			return http.StatusConflict, response
		}

		response.Header.Messages = []string{err.Error()}
		return http.StatusInternalServerError, response
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User.Id = &userID
	return http.StatusCreated, response
}

// NOTE: Check Authenticated cmd/main.go that returns JWT token after successful login
func (s *Server) UserLogin(ctx echo.Context) error {
	return ctx.JSON(s.userLogin(ctx))
}
func (s *Server) userLogin(ctx echo.Context) (int, generated.UserLoginResponse) {
	var (
		context = context.Background()

		response = generated.UserLoginResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	request := generated.User{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusBadRequest, response
	}

	// Get user's phone number from request body
	validPhoneNumber, errorList := validatePhoneNumber(request.PhoneNumber)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	// Get user data
	user, err := s.getSingleUser(context, repository.UserFilter{PhoneNumber: validPhoneNumber})
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusInternalServerError, response
	}

	// Validate password format is valid
	inputPassword, errorList := validatePassword(request.Password)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	// Validate input password (plain) matches user's password (hashed and salted)
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(inputPassword)) != nil {
		response.Header.Messages = []string{"invalid password"}
		return http.StatusBadRequest, response
	}

	// Increment successful login count for the users
	s.Repository.IncrementSuccessfulLoginCount(context, user.ID)

	// Set data to Echo context so we can rely on AuthenticatedMiddleware to generate and return JWT in the Authorization header
	ctx.Set(string(utils.JWTClaimUserID), user.ID)
	ctx.Set(string(utils.JWTClaimPermissions), []utils.JWTPermission{utils.JWTPermissionGetUser, utils.JWTPermissionUpdateUser})

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User.Id = &user.ID
	return http.StatusOK, response
}

// NOTE: Check AuthenticationMiddleware cmd/main.go that authenticates the JWT token
func (s *Server) GetUser(ctx echo.Context) error {
	return ctx.JSON(s.getUser(ctx))
}
func (s *Server) getUser(ctx echo.Context) (int, generated.GetUserResponse) {
	var (
		context = context.Background()

		response = generated.GetUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	// Authorize and get userID of the requester
	userID, err := authorize(ctx, utils.JWTPermissionGetUser)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusForbidden, response
	}

	// Get data for the userID
	user, err := s.getSingleUser(context, repository.UserFilter{UserID: userID})
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusInternalServerError, response
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User = generated.User{
		FullName:    &user.FullName,
		PhoneNumber: &user.PhoneNumber,
	}

	return http.StatusOK, response
}

// NOTE: Check AuthenticationMiddleware cmd/main.go that authenticates the JWT token
func (s *Server) UpdateUser(ctx echo.Context) error {
	return ctx.JSON(s.updateUser(ctx))
}
func (s *Server) updateUser(ctx echo.Context) (int, generated.UpdateUserResponse) {
	var (
		context = context.Background()

		response = generated.UpdateUserResponse{
			Header: generated.ResponseHeader{}, //success is false by default
		}
	)

	// Authenticate and get userID of the requester
	userID, err := authorize(ctx, utils.JWTPermissionUpdateUser)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusForbidden, response
	}

	// Update user data
	request := generated.User{}
	err = json.NewDecoder(ctx.Request().Body).Decode(&request)
	if err != nil {
		response.Header.Messages = []string{err.Error()}
		return http.StatusBadRequest, response
	}

	updateRequest, errorList := fnConvertUpdateUserRequestToUser(userID, request)
	if len(errorList) > 0 {
		response.Header.Messages = errorList
		return http.StatusBadRequest, response
	}

	if err := s.Repository.UpdateUser(context, updateRequest); err != nil {
		if utils.IsUniqueConstraintViolation(err) {
			response.Header.Messages = []string{duplicatePhoneNumberErrorMsg}
			return http.StatusConflict, response
		}

		response.Header.Messages = []string{err.Error()}
		return http.StatusInternalServerError, response
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	return http.StatusOK, response
}
