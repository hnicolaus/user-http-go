package handler

import (
	"context"
	"errors"
	"strings"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/labstack/echo/v4"
)

const (
	successMsg = "request successful"
)

func authenticate(ctx echo.Context, requiredPermission utils.JWTPermission) (userID int64, err error) {
	permissions, ok := ctx.Get(string(utils.JWTClaimPermissions)).([]utils.JWTPermission)
	if !ok {
		return userID, errors.New("not authorized: roles are missing")
	}

	hasRole := false
	for _, permission := range permissions {
		if permission == requiredPermission {
			hasRole = true
		}
	}

	if !hasRole {
		return userID, errors.New("not authorized: no accepted roles")
	}

	userID, ok = ctx.Get(string(utils.JWTClaimUserID)).(int64)
	if !ok {
		return userID, errors.New("missing user_id")
	}

	return userID, nil
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

func convertUpdateUserRequestToUser(userID int64, request generated.User) (user repository.User, errorMsgs []string) {
	if request.PhoneNumber != nil {
		validPhoneNumber, phoneNumberErrorMsgs := validatePhoneNumber(request.PhoneNumber)
		user.PhoneNumber = validPhoneNumber

		errorMsgs = append(errorMsgs, phoneNumberErrorMsgs...)
	}

	if request.FullName != nil {
		validFullName, fullNameErrorMsgs := validateFullName(request.FullName)
		user.FullName = validFullName

		errorMsgs = append(errorMsgs, fullNameErrorMsgs...)

	}

	if len(errorMsgs) > 0 {
		return repository.User{}, errorMsgs
	}

	user.ID = userID
	return user, nil
}

func (s *Server) getSingleUser(ctx context.Context, userFilter repository.UserFilter) (user repository.User, err error) {
	if userFilter == (repository.UserFilter{}) {
		return user, errors.New("userFilter cannot be empty to get a single user")
	}

	users, err := s.Repository.GetUsers(ctx, userFilter)
	if err != nil {
		return user, err
	}

	if len(users) == 0 {
		return user, errors.New("user not found")
	}

	return users[0], nil
}
