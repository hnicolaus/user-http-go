package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
)

func Test_validatePhoneNumber(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name  string
		input *string

		wantValidPhoneNumber string
		wantErrorList        []string
		wantError            bool
	}{
		{
			name:                 "success",
			input:                stringPtr("+628123456789"),
			wantValidPhoneNumber: "+628123456789",
			wantErrorList:        nil,
		},
		{
			name:                 "nil",
			input:                nil,
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62 (rule 2)",
				"phone_number should be 10 to 13 digits (rule 1)",
			},
		},
		{
			name:                 "empty",
			input:                stringPtr(""),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62 (rule 2)",
				"phone_number should be 10 to 13 digits (rule 1)",
			},
		},
		{
			name:                 "fail-rule-1-length-min",
			input:                stringPtr("+62812345  "),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should be 10 to 13 digits (rule 1)",
			},
		},
		{
			name:                 "fail-rule-1-length-max",
			input:                stringPtr("  +6281234567890"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should be 10 to 13 digits (rule 1)",
			},
		},
		{
			name:                 "fail-rule-1-non-numbers",
			input:                stringPtr("+628123456a"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should only contain numbers (rule 1)",
			},
		},
		{
			name:                 "fail-rule-2",
			input:                stringPtr("08123456789"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62 (rule 2)",
			},
		},
		{
			name:                 "fail-all-rules",
			input:                stringPtr("0345abc"),
			wantValidPhoneNumber: "",
			wantErrorList: []string{
				"phone_number should start with +62 (rule 2)",
				"phone_number should be 10 to 13 digits (rule 1)",
				"phone_number should only contain numbers (rule 1)",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotValidPhoneNumber, gotErrorList := validatePhoneNumber(test.input)
			if !reflect.DeepEqual(gotValidPhoneNumber, test.wantValidPhoneNumber) {
				t.Errorf("util.validatePhoneNumber() gotValidPhoneNumber = %v, wantValidPhoneNumber %v", gotValidPhoneNumber, test.wantValidPhoneNumber)
			}

			if !reflect.DeepEqual(gotErrorList, test.wantErrorList) {
				t.Errorf("util.validatePhoneNumber() gotErrorList = %v, wantErrorList %v", gotErrorList, test.wantErrorList)
			}
		})
	}
}

func Test_validateFullName(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name  string
		input *string

		wantValidFullName string
		wantErrorList     []string
	}{
		{
			name:              "success",
			input:             stringPtr("abcd"),
			wantValidFullName: "abcd",
			wantErrorList:     nil,
		},
		{
			name:              "nil",
			input:             nil,
			wantValidFullName: "",
			wantErrorList: []string{
				"full_name should be 3 to 60 characters (rule 3)",
			},
		},
		{
			name:              "fail-length-min",
			input:             stringPtr(" ab"),
			wantValidFullName: "",
			wantErrorList: []string{
				"full_name should be 3 to 60 characters (rule 3)",
			},
		},
		{
			name:              "fail-length-max",
			input:             stringPtr("    S2EeAKi6fze0JVsVbo6OR9uxmzdy89Kiy59z4Wzi2jTdomVUSUIh8G1GmHpJF "),
			wantValidFullName: "",
			wantErrorList: []string{
				"full_name should be 3 to 60 characters (rule 3)",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotValidFullName, gotErrorList := validateFullName(test.input)
			if !reflect.DeepEqual(gotValidFullName, test.wantValidFullName) {
				t.Errorf("util.validateFullName() gotValidFullName = %v, wantValidFullName %v", gotValidFullName, test.wantValidFullName)
			}

			if !reflect.DeepEqual(gotErrorList, test.wantErrorList) {
				t.Errorf("util.validateFullName() gotErrorList = %v, wantErrorList %v", gotErrorList, test.wantErrorList)
			}
		})
	}
}

func Test_validatePassword(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name  string
		input *string

		wantValidPassword string
		wantErrorList     []string
	}{
		{
			name:              "success",
			input:             stringPtr("A1.123"),
			wantValidPassword: "A1.123",
			wantErrorList:     nil,
		},
		{
			name:              "fail-all-rules",
			input:             stringPtr(""),
			wantValidPassword: "",
			wantErrorList: []string{
				"password should be 6 to 64 characters (rule 4)",
				"password should contain a capital letter (rule 4)",
				"password should contain a number (rule 4)",
				"password should contain a special alphanumeric character (rule 4)",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotValidPassword, gotErrorList := validatePassword(test.input)
			if !reflect.DeepEqual(gotValidPassword, test.wantValidPassword) {
				t.Errorf("util.validatePassword() gotValidPassword = %v, wantValidPassword %v", gotValidPassword, test.wantValidPassword)
			}

			if !reflect.DeepEqual(gotErrorList, test.wantErrorList) {
				t.Errorf("util.validatePassword() gotErrorList = %v, wantErrorList %v", gotErrorList, test.wantErrorList)
			}
		})
	}
}

func Test_convertRegisterUserRequestToUser(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name                  string
		input                 generated.User
		fnValidatePhoneNumber func(*string) (string, []string)
		fnValidatePassword    func(*string) (string, []string)
		fnValidateFullName    func(*string) (string, []string)

		wantUser      repository.User
		wantErrorMsgs []string
	}{
		{
			name: "success",
			input: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnValidatePhoneNumber: func(*string) (string, []string) {
				return "+628123456789", []string{}
			},
			fnValidateFullName: func(*string) (string, []string) {
				return "SawitPro User", []string{}
			},
			fnValidatePassword: func(*string) (string, []string) {
				return "P455w0rd!.", []string{}
			},
			wantUser: repository.User{
				FullName:    "SawitPro User",
				PhoneNumber: "+628123456789",
				Password:    "P455w0rd!.",
			},
			wantErrorMsgs: nil,
		},
		{
			name: "fail-all-validations",
			input: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnValidatePhoneNumber: func(*string) (string, []string) {
				return "", []string{"invalid-phone-number-1", "invalid-phone-number-2", "invalid-phone-number-3"}
			},
			fnValidateFullName: func(*string) (string, []string) {
				return "", []string{"invalid-full-name-1", "invalid-full-name-2"}
			},
			fnValidatePassword: func(*string) (string, []string) {
				return "", []string{"invalid-password"}
			},
			wantUser: repository.User{},
			wantErrorMsgs: []string{
				"invalid-phone-number-1", "invalid-phone-number-2", "invalid-phone-number-3",
				"invalid-full-name-1", "invalid-full-name-2",
				"invalid-password",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fnValidateFullName = test.fnValidateFullName
			fnValidatePassword = test.fnValidatePassword
			fnValidatePhoneNumber = test.fnValidatePhoneNumber

			gotUser, gotErrorMsgs := convertRegisterUserRequestToUser(test.input)
			if !reflect.DeepEqual(gotUser, test.wantUser) {
				t.Errorf("util.convertRegisterUserRequestToUser() gotUser = %v, wantUser %v", gotUser, test.wantUser)
			}

			if !reflect.DeepEqual(gotErrorMsgs, test.wantErrorMsgs) {
				t.Errorf("util.convertRegisterUserRequestToUser() gotErrorMsgs = %v, wantErrorMsgs %v", gotErrorMsgs, test.wantErrorMsgs)
			}
		})
	}
}

func Test_convertUpdateUserRequestToUser(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name                  string
		inputUserID           int64
		input                 generated.User
		fnValidatePhoneNumber func(*string) (string, []string)
		fnValidateFullName    func(*string) (string, []string)

		wantUser      repository.User
		wantErrorMsgs []string
	}{
		{
			name:        "success",
			inputUserID: 123,
			input: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnValidatePhoneNumber: func(*string) (string, []string) {
				return "+628123456789", []string{}
			},
			fnValidateFullName: func(*string) (string, []string) {
				return "SawitPro User", []string{}
			},
			wantUser: repository.User{
				ID:          123,
				FullName:    "SawitPro User",
				PhoneNumber: "+628123456789",
			},
			wantErrorMsgs: nil,
		},
		{
			name: "fail-all-validations",
			input: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnValidatePhoneNumber: func(*string) (string, []string) {
				return "", []string{"invalid-phone-number-1", "invalid-phone-number-2", "invalid-phone-number-3"}
			},
			fnValidateFullName: func(*string) (string, []string) {
				return "", []string{"invalid-full-name-1", "invalid-full-name-2"}
			},
			wantUser: repository.User{},
			wantErrorMsgs: []string{
				"invalid-phone-number-1", "invalid-phone-number-2", "invalid-phone-number-3",
				"invalid-full-name-1", "invalid-full-name-2",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fnValidateFullName = test.fnValidateFullName
			fnValidatePhoneNumber = test.fnValidatePhoneNumber

			gotUser, gotErrorMsgs := convertUpdateUserRequestToUser(test.inputUserID, test.input)
			if !reflect.DeepEqual(gotUser, test.wantUser) {
				t.Errorf("util.convertUpdateUserRequestToUser() gotUser = %v, wantUser %v", gotUser, test.wantUser)
			}

			if !reflect.DeepEqual(gotErrorMsgs, test.wantErrorMsgs) {
				t.Errorf("util.convertUpdateUserRequestToUser() gotErrorMsgs = %v, wantErrorMsgs %v", gotErrorMsgs, test.wantErrorMsgs)
			}
		})
	}
}

func Test_authorize(t *testing.T) {
	tests := []struct {
		name               string
		ctxUserID          int64
		ctxPermissions     []utils.JWTPermission
		requiredPermission utils.JWTPermission

		wantUserID int64
		wantErr    error
	}{
		{
			name:      "success",
			ctxUserID: 123,
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
				utils.JWTPermissionUpdateUser,
			},
			requiredPermission: utils.JWTPermissionGetUser,
			wantUserID:         123,
			wantErr:            nil,
		},
		{
			name:      "fail-not-authorized-permission",
			ctxUserID: 123,
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionUpdateUser,
			},
			requiredPermission: utils.JWTPermissionGetUser,
			wantUserID:         0,
			wantErr:            errors.New("not authorized: missing required permission"),
		},
		{
			name:      "fail-not-authorized-user-id",
			ctxUserID: 0,
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			requiredPermission: utils.JWTPermissionGetUser,
			wantUserID:         0,
			wantErr:            errors.New("missing user_id"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/v1/unit-test", nil)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)
			if test.ctxUserID != 0 {
				ctx.Set(string(utils.JWTClaimUserID), test.ctxUserID)
			}
			if len(test.ctxPermissions) > 0 {
				ctx.Set(string(utils.JWTClaimPermissions), test.ctxPermissions)
			}

			gotUserID, gotErr := authorize(ctx, test.requiredPermission)
			if gotUserID != test.wantUserID {
				t.Errorf("util.authorize() gotUserID = %v, wantUserID %v", gotUserID, test.wantUserID)
			}

			if !reflect.DeepEqual(gotErr, test.wantErr) {
				t.Errorf("util.authorize() gotErr = %v, wantErr %v", gotErr, test.wantErr)
			}
		})
	}
}

func Test_getSingleUser(t *testing.T) {
	tests := []struct {
		name           string
		userFilter     repository.UserFilter
		mockRepository func(controller *gomock.Controller) *repository.MockRepositoryInterface

		wantUser repository.User
		wantErr  error
	}{
		{
			name: "success",
			userFilter: repository.UserFilter{
				UserID:      123,
				PhoneNumber: "+628123456789",
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					UserID:      123,
					PhoneNumber: "+628123456789",
				}).Return([]repository.User{
					{
						ID:          123,
						FullName:    "SawitPro User",
						PhoneNumber: "+628123456789",
						Password:    "$2a$12$41bm0d9VyLDKALovox4S9.FoNezvO9tB8ck94/0fEyKcYIFmV8guq",
					},
				}, nil)

				return mock
			},
			wantUser: repository.User{
				ID:          123,
				FullName:    "SawitPro User",
				PhoneNumber: "+628123456789",
				Password:    "$2a$12$41bm0d9VyLDKALovox4S9.FoNezvO9tB8ck94/0fEyKcYIFmV8guq",
			},
			wantErr: nil,
		},
		{
			name: "fail-get-user",
			userFilter: repository.UserFilter{
				UserID:      123,
				PhoneNumber: "+628123456789",
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					UserID:      123,
					PhoneNumber: "+628123456789",
				}).Return([]repository.User{}, errors.New("error-get-user"))

				return mock
			},
			wantUser: repository.User{},
			wantErr:  errors.New("error-get-user"),
		},
		{
			name: "fail-user-not-found",
			userFilter: repository.UserFilter{
				UserID:      123,
				PhoneNumber: "+628123456789",
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					UserID:      123,
					PhoneNumber: "+628123456789",
				}).Return([]repository.User{}, nil)

				return mock
			},
			wantUser: repository.User{},
			wantErr:  errors.New("user not found"),
		},
		{
			name:       "fail-empty-filter",
			userFilter: repository.UserFilter{},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)
				return mock
			},
			wantUser: repository.User{},
			wantErr:  errors.New("userFilter cannot be empty to get a single user"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Repository: test.mockRepository(controller),
			}

			gotUser, gotErr := handler.getSingleUser(context.Background(), test.userFilter)

			if !reflect.DeepEqual(test.wantUser, gotUser) {
				t.Errorf("handler.getSingleUser() gotUser = %v, wantUser %v", gotUser, test.wantUser)
			}

			if !reflect.DeepEqual(test.wantErr, gotErr) {
				t.Errorf("handler.getSingleUser() gotErr = %v, wantErr %v", gotErr, test.wantErr)
			}
		})
	}
}
