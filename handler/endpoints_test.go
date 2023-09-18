package handler

import (
	"bytes"
	"encoding/json"
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
	"github.com/lib/pq"
)

func TestRegisterUser(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	int64Ptr := func(in int64) *int64 {
		return &in
	}

	errorConflictUserPhoneNumber := pq.Error{
		Code: "23505",
	}

	tests := []struct {
		name                               string
		mockRepository                     func(controller *gomock.Controller) *repository.MockRepositoryInterface
		requestBody                        generated.User
		fnConvertRegisterUserRequestToUser func(generated.User) (repository.User, []string)
		wantResponse                       generated.RegisterUserResponse
		wantHttpStatusCode                 int
	}{
		{
			name: "success",
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (repository.User, []string) {
				user := repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().InsertUser(gomock.Any(), repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(123), nil)

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
				User: generated.User{
					Id: int64Ptr(123),
				},
			},
			wantHttpStatusCode: http.StatusOK,
		},
		{
			name: "fail-insert-user",
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (repository.User, []string) {
				user := repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().InsertUser(gomock.Any(), repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(0), errors.New("error-insert-user"))

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"error-insert-user"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fail-insert-user-conflict-phone-number",
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (repository.User, []string) {
				user := repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().InsertUser(gomock.Any(), repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(0), &errorConflictUserPhoneNumber)

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"phone number is already registered to an existing user"},
				},
			},
			wantHttpStatusCode: http.StatusConflict,
		},
		{
			name: "fail-insert-user-conflict-phone-number",
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("P455w0rd!."),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (repository.User, []string) {
				user := repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().InsertUser(gomock.Any(), repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
					Password:    "P455w0rd!.",
				}).Return(int64(0), &errorConflictUserPhoneNumber)

				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{duplicatePhoneNumberErrorMsg},
				},
			},
			wantHttpStatusCode: http.StatusConflict,
		},
		{
			name: "fail-invalid-input",
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro"),
				PhoneNumber: stringPtr("+62812"),
				Password:    stringPtr("P455w"),
			},
			fnConvertRegisterUserRequestToUser: func(generated.User) (repository.User, []string) {
				return repository.User{}, []string{"invalid full name", "invalid phone number", "invalid password"}
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)
				return mock
			},
			wantResponse: generated.RegisterUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"invalid full name", "invalid phone number", "invalid password"},
				},
			},
			wantHttpStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Repository: test.mockRepository(controller),
			}

			requestBodyJSON, _ := json.Marshal(test.requestBody)
			requestBody := []byte(requestBodyJSON)
			requestBodyBuffer := bytes.NewBuffer(requestBody)

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/v1/user", requestBodyBuffer)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)

			fnConvertRegisterUserRequestToUser = test.fnConvertRegisterUserRequestToUser

			gotHttpStatusCode, gotResponse := handler.registerUser(ctx)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.RegisterUser() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.RegisterUser() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}

		})
	}
}

func TestUserLogin(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	int64Ptr := func(in int64) *int64 {
		return &in
	}

	tests := []struct {
		name               string
		mockRepository     func(controller *gomock.Controller) *repository.MockRepositoryInterface
		requestBody        generated.User
		wantResponse       generated.UserLoginResponse
		wantCtxUserID      int64
		wantHttpStatusCode int
	}{
		{
			name: "success",
			requestBody: generated.User{
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("Password123!."),
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					PhoneNumber: "+628123456789",
				}).Return([]repository.User{
					{
						ID:          123,
						FullName:    "SawitPro User",
						PhoneNumber: "+628123456789",
						Password:    "$2a$12$41bm0d9VyLDKALovox4S9.FoNezvO9tB8ck94/0fEyKcYIFmV8guq",
					},
				}, nil)

				mock.EXPECT().IncrementSuccessfulLoginCount(gomock.Any(), int64(123)).Return(nil)

				return mock
			},
			wantResponse: generated.UserLoginResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
				User: generated.User{
					Id: int64Ptr(123),
				},
			},
			wantCtxUserID:      123,
			wantHttpStatusCode: http.StatusOK,
		},
		{
			name: "fail-invalid-password",
			requestBody: generated.User{
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("Password123.!"),
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
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
			wantResponse: generated.UserLoginResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"invalid password"},
				},
			},
			wantHttpStatusCode: http.StatusBadRequest,
		},
		{
			name: "fail-get-user",
			requestBody: generated.User{
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("Password123.!"),
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					PhoneNumber: "+628123456789",
				}).Return([]repository.User{}, errors.New("error-get-users"))

				return mock
			},
			wantResponse: generated.UserLoginResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"error-get-users"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fail-user-does-not-exist",
			requestBody: generated.User{
				PhoneNumber: stringPtr("+628123456789"),
				Password:    stringPtr("Password123!."),
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					PhoneNumber: "+628123456789",
				}).Return([]repository.User{}, nil)

				return mock
			},
			wantResponse: generated.UserLoginResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"user not found"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Repository: test.mockRepository(controller),
			}

			requestBodyJSON, _ := json.Marshal(test.requestBody)
			requestBody := []byte(requestBodyJSON)
			requestBodyBuffer := bytes.NewBuffer(requestBody)

			e := echo.New()
			request := httptest.NewRequest(http.MethodPost, "/v1/user/login", requestBodyBuffer)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)

			gotHttpStatusCode, gotResponse := handler.userLogin(ctx)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.UserLogin() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.UserLogin() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}

			if test.wantResponse.Header.Success {
				gotCtxUserID, _ := ctx.Get(string(utils.JWTClaimUserID)).(int64)
				if gotCtxUserID != test.wantCtxUserID {
					t.Errorf("handler.UserLogin() gotCtxUserID = %v, wantCtxUserID %v", gotCtxUserID, test.wantCtxUserID)
				}

				gotCtxPermissions, _ := ctx.Get(string(utils.JWTClaimPermissions)).([]utils.JWTPermission)
				wantCtxPermissions := []utils.JWTPermission{utils.JWTPermissionGetUser, utils.JWTPermissionUpdateUser}
				if !reflect.DeepEqual(gotCtxPermissions, wantCtxPermissions) {
					t.Errorf("handler.UserLogin() gotCtxPermissions = %v, wantCtxPermissions %v", gotCtxPermissions, wantCtxPermissions)
				}
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	tests := []struct {
		name           string
		mockRepository func(controller *gomock.Controller) *repository.MockRepositoryInterface
		ctxPermissions []utils.JWTPermission
		ctxUserID      int64

		wantResponse       generated.GetUserResponse
		wantHttpStatusCode int
	}{
		{
			name: "success",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			ctxUserID: 123,
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					UserID: 123,
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
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
				User: generated.User{
					FullName:    stringPtr("SawitPro User"),
					PhoneNumber: stringPtr("+628123456789"),
				},
			},
			wantHttpStatusCode: http.StatusOK,
		},
		{
			name:           "fail-not-authorized-no-permission",
			ctxPermissions: []utils.JWTPermission{},
			ctxUserID:      123,
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)
				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"not authorized: missing required permission"},
				},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name:           "fail-not-authorized-wrong-permission",
			ctxPermissions: []utils.JWTPermission{utils.JWTPermissionUpdateUser},
			ctxUserID:      123,
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)
				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"not authorized: missing required permission"},
				},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name:           "fail-not-authorized-no-user-id",
			ctxPermissions: []utils.JWTPermission{utils.JWTPermissionGetUser},
			ctxUserID:      0,
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)
				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"missing user_id"},
				},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name: "fail-get-user",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			ctxUserID: 123,
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					UserID: 123,
				}).Return([]repository.User{}, errors.New("error-get-user"))

				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"error-get-user"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fail-get-user",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			ctxUserID: 123,
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().GetUsers(gomock.Any(), repository.UserFilter{
					UserID: 123,
				}).Return([]repository.User{}, nil)

				return mock
			},
			wantResponse: generated.GetUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"user not found"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Repository: test.mockRepository(controller),
			}

			e := echo.New()
			request := httptest.NewRequest(http.MethodGet, "/v1/user", nil)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)

			if test.ctxUserID != 0 {
				ctx.Set(string(utils.JWTClaimUserID), test.ctxUserID)
			}
			if len(test.ctxPermissions) > 0 {
				ctx.Set(string(utils.JWTClaimPermissions), test.ctxPermissions)
			}

			gotHttpStatusCode, gotResponse := handler.getUser(ctx)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.GetUser() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.GetUser() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	stringPtr := func(in string) *string {
		return &in
	}

	errorConflictUserPhoneNumber := pq.Error{
		Code: "23505",
	}

	tests := []struct {
		name                             string
		mockRepository                   func(controller *gomock.Controller) *repository.MockRepositoryInterface
		ctxPermissions                   []utils.JWTPermission
		ctxUserID                        int64
		requestBody                      generated.User
		fnConvertUpdateUserRequestToUser func(int64, generated.User) (repository.User, []string)

		wantResponse       generated.UpdateUserResponse
		wantHttpStatusCode int
	}{
		{
			name: "success",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionUpdateUser,
			},
			ctxUserID: 123,
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
			},
			fnConvertUpdateUserRequestToUser: func(int64, generated.User) (repository.User, []string) {
				user := repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().UpdateUser(gomock.Any(), repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
				}).Return(nil)

				return mock
			},
			wantResponse: generated.UpdateUserResponse{
				Header: generated.ResponseHeader{
					Success:  true,
					Messages: []string{successMsg},
				},
			},
			wantHttpStatusCode: http.StatusOK,
		},
		{
			name: "fail-not-authorized-permission",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionGetUser,
			},
			ctxUserID: 123,
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)
				return mock
			},
			wantResponse: generated.UpdateUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"not authorized: missing required permission"},
				},
			},
			wantHttpStatusCode: http.StatusForbidden,
		},
		{
			name: "fail-update-user",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionUpdateUser,
			},
			ctxUserID: 123,
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
			},
			fnConvertUpdateUserRequestToUser: func(int64, generated.User) (repository.User, []string) {
				user := repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().UpdateUser(gomock.Any(), repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
				}).Return(errors.New("error-update-user"))

				return mock
			},
			wantResponse: generated.UpdateUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{"error-update-user"},
				},
			},
			wantHttpStatusCode: http.StatusInternalServerError,
		},
		{
			name: "fail-update-user-conflict-phone-number",
			ctxPermissions: []utils.JWTPermission{
				utils.JWTPermissionUpdateUser,
			},
			ctxUserID: 123,
			requestBody: generated.User{
				FullName:    stringPtr("SawitPro User"),
				PhoneNumber: stringPtr("+628123456789"),
			},
			fnConvertUpdateUserRequestToUser: func(int64, generated.User) (repository.User, []string) {
				user := repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
				}

				errorMsgs := []string{}

				return user, errorMsgs
			},
			mockRepository: func(controller *gomock.Controller) *repository.MockRepositoryInterface {
				mock := repository.NewMockRepositoryInterface(controller)

				mock.EXPECT().UpdateUser(gomock.Any(), repository.User{
					FullName:    "SawitPro User",
					PhoneNumber: "+628123456789",
				}).Return(&errorConflictUserPhoneNumber)

				return mock
			},
			wantResponse: generated.UpdateUserResponse{
				Header: generated.ResponseHeader{
					Success:  false,
					Messages: []string{duplicatePhoneNumberErrorMsg},
				},
			},
			wantHttpStatusCode: http.StatusConflict,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			controller := gomock.NewController(t)
			handler := &Server{
				Repository: test.mockRepository(controller),
			}

			requestBodyJSON, _ := json.Marshal(test.requestBody)
			requestBody := []byte(requestBodyJSON)
			requestBodyBuffer := bytes.NewBuffer(requestBody)

			e := echo.New()
			request := httptest.NewRequest(http.MethodPut, "/v1/user", requestBodyBuffer)
			recorder := httptest.NewRecorder()
			ctx := e.NewContext(request, recorder)
			if test.ctxUserID != 0 {
				ctx.Set(string(utils.JWTClaimUserID), test.ctxUserID)
			}
			if len(test.ctxPermissions) > 0 {
				ctx.Set(string(utils.JWTClaimPermissions), test.ctxPermissions)
			}

			fnConvertUpdateUserRequestToUser = test.fnConvertUpdateUserRequestToUser

			gotHttpStatusCode, gotResponse := handler.updateUser(ctx)

			if gotHttpStatusCode != test.wantHttpStatusCode {
				t.Errorf("handler.UpdateUser() httpStatusCode = %v, wantHttpStatusCode %v", gotHttpStatusCode, test.wantHttpStatusCode)
			}

			if !reflect.DeepEqual(test.wantResponse, gotResponse) {
				t.Errorf("handler.UpdateUser() response = %v, wantResponse %v", gotResponse, test.wantResponse)
			}
		})
	}
}
