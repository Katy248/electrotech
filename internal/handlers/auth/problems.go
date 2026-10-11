package auth

import (
	"electrotech/pkg/problem"
	"net/http"
)

var (
	errAuth = problem.New(
		"err:electrotech/auth", http.StatusUnauthorized, "auth error")
	errNoHeader = errAuth.Extend(
		"/no-header", http.StatusUnauthorized, "no authorization header provided")
	errInvalidToken = errAuth.Extend(
		"/invalid-token", http.StatusUnauthorized, "invalid JWT token")
	errInvalidCredentials = errAuth.Extend(
		"invalid-credentials",
		http.StatusUnauthorized,
		"invalid credentials",
	)
	errBadRequest = errAuth.Extend(
		"/bad-request",
		http.StatusBadRequest,
		"bad request",
	)
	errCreateAuthResponse = errAuth.Extend(
		"/internal",
		http.StatusInternalServerError,
		"internal server error").WithDetails("creating auth response")
)
