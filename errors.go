package smarthome_sdk

import "errors"

var (
	ErrNotInitialized      = errors.New("action failed: initialize connection first")
	ErrInvalidURL          = errors.New("invalid url: the url could not be parsed")
	ErrConnFailed          = errors.New("connection failed: request failed due to network issues")
	ErrServiceUnavailable  = errors.New("request failed: smarthome is currently unavailable")
	ErrInternalServerError = errors.New("request failed: smarthome failed internally")
	ErrInvalidCredentials  = errors.New("authentication failed: invalid credentials")
	ErrNoCookiesSent       = errors.New("login request did not respond with an expected cookie")
	ErrUnknown             = errors.New("request failed with unknown error")
)
