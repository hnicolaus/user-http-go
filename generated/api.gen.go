// Package generated provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package generated

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

// GetUserResponse defines model for GetUserResponse.
type GetUserResponse struct {
	Header ResponseHeader `json:"header"`
	User   User           `json:"user"`
}

// RegisterUserResponse defines model for RegisterUserResponse.
type RegisterUserResponse struct {
	Header ResponseHeader `json:"header"`
	User   User           `json:"user"`
}

// ResponseHeader defines model for ResponseHeader.
type ResponseHeader struct {
	// Messages Array of error message(s).
	Messages []string `json:"messages"`

	// Success Boolean to denote whether response is OK or not.
	Success bool `json:"success"`
}

// UpdateUserResponse defines model for UpdateUserResponse.
type UpdateUserResponse struct {
	Header ResponseHeader `json:"header"`
}

// User defines model for User.
type User struct {
	// FullName User's full name.
	FullName *string `json:"full_name,omitempty"`
	Id       *int64  `json:"id,omitempty"`

	// Password User's password.
	Password *string `json:"password,omitempty"`

	// PhoneNumber User's phone number.
	PhoneNumber *string `json:"phone_number,omitempty"`
}

// UserLoginResponse defines model for UserLoginResponse.
type UserLoginResponse struct {
	Header ResponseHeader `json:"header"`
	User   User           `json:"user"`
}

// RegisterUserJSONRequestBody defines body for RegisterUser for application/json ContentType.
type RegisterUserJSONRequestBody = User

// UpdateUserJSONRequestBody defines body for UpdateUser for application/json ContentType.
type UpdateUserJSONRequestBody = User

// UserLoginJSONRequestBody defines body for UserLogin for application/json ContentType.
type UserLoginJSONRequestBody = User

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get an existing new user
	// (GET /v1/user)
	GetUser(ctx echo.Context) error
	// Create a new user
	// (POST /v1/user)
	RegisterUser(ctx echo.Context) error
	// Update an existing user
	// (PUT /v1/user)
	UpdateUser(ctx echo.Context) error
	// Existing user login
	// (POST /v1/user/login)
	UserLogin(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetUser(ctx)
	return err
}

// RegisterUser converts echo context to params.
func (w *ServerInterfaceWrapper) RegisterUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.RegisterUser(ctx)
	return err
}

// UpdateUser converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UpdateUser(ctx)
	return err
}

// UserLogin converts echo context to params.
func (w *ServerInterfaceWrapper) UserLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.UserLogin(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/v1/user", wrapper.GetUser)
	router.POST(baseURL+"/v1/user", wrapper.RegisterUser)
	router.PUT(baseURL+"/v1/user", wrapper.UpdateUser)
	router.POST(baseURL+"/v1/user/login", wrapper.UserLogin)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xVXWsbOxD9K2LuhXsv+GadJi1035rQD9OWQto8hRDk3bFXYVdSR7NOTdj/XiTt1nZX",
	"IdQklPTJa3HOfJyjGd1CYRprNGp2kN+CKypsZPh8i3zukM7QWaMd+iNLxiKxwgCoUJZI/utvwgXk8Fe2",
	"CZb1kbKB/y6iuwm07n6WzwxdNwHCr60iLCG/GBL2ES4nwGuLkIOZX2PBPvQZLpVjpCdY+E6yUckNOieX",
	"8btEV5CyrIyGHF4RybUwC4FEhkQP/Nf9dwATUIxN4PQJHZPSS5+wP5Ce7f+7tijQJeKfGFOj1IKNKFEb",
	"RnFTIVdIgvqahXLi03thSGjDPmsfex6ZIzWGVJNNVylJzm0pGR/DybQ7yRpcyoxFW9dXWjY4VssT/nHC",
	"I4RHbMmx0V6VIYqhRjLkoDS/ON7glGZcxgtnpXM3hso78wyAZBpbGY1Xum3msYl0BA8SEZSI0t0hygez",
	"VPoJTZjHK70wPnitCuyrjibCx9mXMBOKa+ylEZ+RVqpAmMAKyUXNDg+mB1OPNBa1tApyOApH3iquQuvZ",
	"6jAb2lgi+x+vjPSyz0rIh70KvoHYf+A9m079T2E0ow40aW2tikDMrp0vYNjP90n08+oO7Y/tFwWhZCxF",
	"P5H+1oZlcDw9Gl+YN4bmqixRe8TzWO0uYqYZSctaOKQVUtxIwSrXNo2kdexeSC3wm3Ks9FJovBFBLn9h",
	"jUvotb3TIbqOjk9MuX4wvfortXOnmFrsRh4dPljO5FP160YlbDiRpehVEv+LmV7JWpVCadty5Lwcc06N",
	"XtQqzvfe3p6GKoXcNbVNeLrZ7b/d0YebusSDdZefbYDuM3iP514sf2c422HVDkstq/3eD/s+Oas/noY/",
	"ydbRc5dwNQC27Nx3OPd27/W2ZyLaFOqMHAf5xS20VEMOFbPNs6w2hawrb2N32X0PAAD//0qC0fUIDAAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
