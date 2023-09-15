package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/labstack/echo/v4"

	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
)

const (
	duplicatePhoneNumberErrMsg = "phone number is already registered to an existing user"
	successMsg                 = "user successfully created"
)

func (s *Server) RegisterUser(ctx echo.Context) error {
	var (
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

	userID, err := s.Repository.InsertUser(context.Background(), user)
	if err != nil {
		if utils.IsUniqueConstraintViolation(err) {
			response.Header.Messages = []string{duplicatePhoneNumberErrMsg}
			return ctx.JSON(http.StatusConflict, response)
		}

		response.Header.Messages = []string{err.Error()}
		return ctx.JSON(http.StatusInternalServerError, response)
	}

	response.Header.Success = true
	response.Header.Messages = []string{successMsg}
	response.User.Id = userID
	return ctx.JSON(http.StatusOK, response)
}

func convertRegisterUserRequestToUser(request generated.User) (user repository.User, errorMsgs []string) {
	validPhoneNumber, phoneNumberErrorMsgs := validatePhoneNumber(request.PhoneNumber)
	validFullName, fullNameErrorMsgs := validateFullName(request.FullName)
	validPassword, passwordErrorMsgs := validatePassword(request.Password)

	errorList := append(phoneNumberErrorMsgs, fullNameErrorMsgs...)
	errorList = append(errorList, passwordErrorMsgs...)

	if len(errorList) > 0 {
		return repository.User{}, errorList
	}

	return repository.User{
		FullName:    validFullName,
		PhoneNumber: validPhoneNumber,
		Password:    validPassword,
	}, nil
}

func validatePhoneNumber(input *string) (validPhoneNumber string, errorList []string) {
	phoneNumber := ""

	if input != nil {
		phoneNumber = strings.TrimSpace(*input)
	}

	// Verify "+62" prefix
	if !strings.HasPrefix(phoneNumber, "+62") {
		errorList = append(errorList, "phone_number should start with +62 (rule 2)")
	}

	// Check the length of the phone number
	if len(phoneNumber) < 10 || len(phoneNumber) > 13 {
		errorList = append(errorList, "phone_number should be 10 to 13 digits (rule 1)")
	}

	// Check if all remaining characters are digits
	for i := 3; i < len(phoneNumber); i++ {
		c := phoneNumber[i]
		if !unicode.IsDigit(rune(c)) {
			errorList = append(errorList, "phone_number should only contain numbers (rule 1)")
		}
	}

	if len(errorList) == 0 {
		validPhoneNumber = phoneNumber
	}

	return validPhoneNumber, errorList
}

func validateFullName(input *string) (validFullName string, errorList []string) {
	fullName := ""

	if input != nil {
		fullName = strings.TrimSpace(*input)
	}

	if len(fullName) < 3 || len(fullName) > 60 {
		errorList = append(errorList, "full_name should be 3 to 60 characters (rule 3)")
	}

	if len(errorList) == 0 {
		validFullName = fullName
	}

	return validFullName, errorList
}

func validatePassword(input *string) (validPassword string, errorList []string) {
	password := ""

	if input != nil {
		password = *input
	}

	if len(password) < 6 || len(password) > 64 {
		errorList = append(errorList, "password should be 6 to 64 characters (rule 4)")
	}

	containsCapital, containsNumber, containsSpecialAlphaNumeric := false, false, false
	for _, c := range password {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
			containsSpecialAlphaNumeric = true
		}
		if unicode.IsUpper(c) {
			containsCapital = true
		}
		if unicode.IsNumber(c) {
			containsNumber = true
		}

		if containsSpecialAlphaNumeric && containsCapital && containsNumber {
			break
		}
	}

	if !containsCapital {
		errorList = append(errorList, "password should contain a capital letter (rule 4)")
	}
	if !containsNumber {
		errorList = append(errorList, "password should contain a number (rule 4)")
	}
	if !containsSpecialAlphaNumeric {
		errorList = append(errorList, "password should contain an alphanumeric character (rule 4)")

	}

	if len(errorList) == 0 {
		validPassword = password
	}

	return validPassword, errorList
}
