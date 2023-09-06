// Package v1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.14.0 DO NOT EDIT.
package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// Defines values for ServiceStatusStatus.
const (
	Fail ServiceStatusStatus = "fail"
	Pass ServiceStatusStatus = "pass"
	Warn ServiceStatusStatus = "warn"
)

// ErrorResponse The ErrorResponse contains the Problem Details for HTTP APIs as specified in [RFC7807](https://tools.ietf.org/html/rfc7807).
//
// It provides more details about problems occurred in the storage server.
//
// Return values contain the following members:
// - **title** (string) - A short, human-readable summary of the problem type.
// - **status** (number) - The HTTP status code generated by the origin server for this occurrence of the problem.
// - **backend** (string) The name of the storage backend. This can provide context to the error.
// - **detail** (string) - A human-readable explanation specific to this occurrence of the problem.
type ErrorResponse struct {
	// Backend The name of the storage backend. This can provide context to the error.
	Backend string `json:"backend"`

	// Detail A human-readable explanation specific to this occurrence of the problem.
	Detail string `json:"detail"`

	// Status HTTP status-code
	Status int `json:"status"`

	// Title A short, human-readable summary of the problem type.
	Title string `json:"title"`
}

// Key The key under which secrets can be stored or retrieved.
//
// The key should be considered opaque and no assumptions should be made about its value or format.
// Note: When the key is used in the URL path, symbols such as slashes and hash symbols must be escaped.
type Key = string

// KeyList List of keys currently stored in the store.
// Note: Keys will be in unescaped form. No assumptions should be made about the order of the keys.
type KeyList = []Key

// Secret The secret value stored under the provided key.
type Secret = string

// SecretResponse Response object containing the secret value.
type SecretResponse struct {
	// Secret The secret value stored under the provided key.
	Secret Secret `json:"secret"`
}

// ServiceStatus Response for the health check endpoint.
type ServiceStatus struct {
	// Details Additional details about the service status.
	Details *string `json:"details,omitempty"`

	// Status Indicates whether the service status is acceptable. Possible values are:
	// * **pass**: healthy.
	// * **fail**: unhealthy.
	// * **warn**: healthy, with some concerns.
	Status ServiceStatusStatus `json:"status"`
}

// ServiceStatusStatus Indicates whether the service status is acceptable. Possible values are:
// * **pass**: healthy.
// * **fail**: unhealthy.
// * **warn**: healthy, with some concerns.
type ServiceStatusStatus string

// StoreSecretRequest Request body to store a secret value. The secret value must not be empty.
type StoreSecretRequest struct {
	// Secret The secret value stored under the provided key.
	Secret Secret `json:"secret"`
}

// StoreSecretJSONRequestBody defines body for StoreSecret for application/json ContentType.
type StoreSecretJSONRequestBody = StoreSecretRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Health check
	// (GET /health)
	HealthCheck(ctx echo.Context) error
	// List all keys in the store
	// (GET /secrets)
	ListKeys(ctx echo.Context) error
	// Delete the secret for the provided key.
	// (DELETE /secrets/{key})
	DeleteSecret(ctx echo.Context, key Key) error
	// Lookup the secret for the provided key
	// (GET /secrets/{key})
	LookupSecret(ctx echo.Context, key Key) error
	// Store a new secret under the provided key
	// (POST /secrets/{key})
	StoreSecret(ctx echo.Context, key Key) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// HealthCheck converts echo context to params.
func (w *ServerInterfaceWrapper) HealthCheck(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.HealthCheck(ctx)
	return err
}

// ListKeys converts echo context to params.
func (w *ServerInterfaceWrapper) ListKeys(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.ListKeys(ctx)
	return err
}

// DeleteSecret converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteSecret(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "key" -------------
	var key Key

	err = runtime.BindStyledParameterWithLocation("simple", false, "key", runtime.ParamLocationPath, ctx.Param("key"), &key)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter key: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeleteSecret(ctx, key)
	return err
}

// LookupSecret converts echo context to params.
func (w *ServerInterfaceWrapper) LookupSecret(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "key" -------------
	var key Key

	err = runtime.BindStyledParameterWithLocation("simple", false, "key", runtime.ParamLocationPath, ctx.Param("key"), &key)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter key: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.LookupSecret(ctx, key)
	return err
}

// StoreSecret converts echo context to params.
func (w *ServerInterfaceWrapper) StoreSecret(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "key" -------------
	var key Key

	err = runtime.BindStyledParameterWithLocation("simple", false, "key", runtime.ParamLocationPath, ctx.Param("key"), &key)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter key: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.StoreSecret(ctx, key)
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

	router.GET(baseURL+"/health", wrapper.HealthCheck)
	router.GET(baseURL+"/secrets", wrapper.ListKeys)
	router.DELETE(baseURL+"/secrets/:key", wrapper.DeleteSecret)
	router.GET(baseURL+"/secrets/:key", wrapper.LookupSecret)
	router.POST(baseURL+"/secrets/:key", wrapper.StoreSecret)

}

type HealthCheckRequestObject struct {
}

type HealthCheckResponseObject interface {
	VisitHealthCheckResponse(w http.ResponseWriter) error
}

type HealthCheck200JSONResponse ServiceStatus

func (response HealthCheck200JSONResponse) VisitHealthCheckResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type HealthCheck503JSONResponse ServiceStatus

func (response HealthCheck503JSONResponse) VisitHealthCheckResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

type ListKeysRequestObject struct {
}

type ListKeysResponseObject interface {
	VisitListKeysResponse(w http.ResponseWriter) error
}

type ListKeys200JSONResponse KeyList

func (response ListKeys200JSONResponse) VisitListKeysResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type ListKeys500JSONResponse ErrorResponse

func (response ListKeys500JSONResponse) VisitListKeysResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type DeleteSecretRequestObject struct {
	Key Key `json:"key"`
}

type DeleteSecretResponseObject interface {
	VisitDeleteSecretResponse(w http.ResponseWriter) error
}

type DeleteSecret204Response struct {
}

func (response DeleteSecret204Response) VisitDeleteSecretResponse(w http.ResponseWriter) error {
	w.WriteHeader(204)
	return nil
}

type DeleteSecret404JSONResponse ErrorResponse

func (response DeleteSecret404JSONResponse) VisitDeleteSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type DeleteSecret500JSONResponse ErrorResponse

func (response DeleteSecret500JSONResponse) VisitDeleteSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type LookupSecretRequestObject struct {
	Key Key `json:"key"`
}

type LookupSecretResponseObject interface {
	VisitLookupSecretResponse(w http.ResponseWriter) error
}

type LookupSecret200JSONResponse SecretResponse

func (response LookupSecret200JSONResponse) VisitLookupSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type LookupSecret404JSONResponse ErrorResponse

func (response LookupSecret404JSONResponse) VisitLookupSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type LookupSecret500JSONResponse ErrorResponse

func (response LookupSecret500JSONResponse) VisitLookupSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type StoreSecretRequestObject struct {
	Key  Key `json:"key"`
	Body *StoreSecretJSONRequestBody
}

type StoreSecretResponseObject interface {
	VisitStoreSecretResponse(w http.ResponseWriter) error
}

type StoreSecret200JSONResponse SecretResponse

func (response StoreSecret200JSONResponse) VisitStoreSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type StoreSecret400JSONResponse ErrorResponse

func (response StoreSecret400JSONResponse) VisitStoreSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type StoreSecret409JSONResponse ErrorResponse

func (response StoreSecret409JSONResponse) VisitStoreSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type StoreSecret500JSONResponse ErrorResponse

func (response StoreSecret500JSONResponse) VisitStoreSecretResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Health check
	// (GET /health)
	HealthCheck(ctx context.Context, request HealthCheckRequestObject) (HealthCheckResponseObject, error)
	// List all keys in the store
	// (GET /secrets)
	ListKeys(ctx context.Context, request ListKeysRequestObject) (ListKeysResponseObject, error)
	// Delete the secret for the provided key.
	// (DELETE /secrets/{key})
	DeleteSecret(ctx context.Context, request DeleteSecretRequestObject) (DeleteSecretResponseObject, error)
	// Lookup the secret for the provided key
	// (GET /secrets/{key})
	LookupSecret(ctx context.Context, request LookupSecretRequestObject) (LookupSecretResponseObject, error)
	// Store a new secret under the provided key
	// (POST /secrets/{key})
	StoreSecret(ctx context.Context, request StoreSecretRequestObject) (StoreSecretResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// HealthCheck operation middleware
func (sh *strictHandler) HealthCheck(ctx echo.Context) error {
	var request HealthCheckRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.HealthCheck(ctx.Request().Context(), request.(HealthCheckRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "HealthCheck")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(HealthCheckResponseObject); ok {
		return validResponse.VisitHealthCheckResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// ListKeys operation middleware
func (sh *strictHandler) ListKeys(ctx echo.Context) error {
	var request ListKeysRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.ListKeys(ctx.Request().Context(), request.(ListKeysRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "ListKeys")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(ListKeysResponseObject); ok {
		return validResponse.VisitListKeysResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// DeleteSecret operation middleware
func (sh *strictHandler) DeleteSecret(ctx echo.Context, key Key) error {
	var request DeleteSecretRequestObject

	request.Key = key

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.DeleteSecret(ctx.Request().Context(), request.(DeleteSecretRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "DeleteSecret")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(DeleteSecretResponseObject); ok {
		return validResponse.VisitDeleteSecretResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// LookupSecret operation middleware
func (sh *strictHandler) LookupSecret(ctx echo.Context, key Key) error {
	var request LookupSecretRequestObject

	request.Key = key

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.LookupSecret(ctx.Request().Context(), request.(LookupSecretRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "LookupSecret")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(LookupSecretResponseObject); ok {
		return validResponse.VisitLookupSecretResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// StoreSecret operation middleware
func (sh *strictHandler) StoreSecret(ctx echo.Context, key Key) error {
	var request StoreSecretRequestObject

	request.Key = key

	var body StoreSecretJSONRequestBody
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	request.Body = &body

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.StoreSecret(ctx.Request().Context(), request.(StoreSecretRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "StoreSecret")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(StoreSecretResponseObject); ok {
		return validResponse.VisitStoreSecretResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}
