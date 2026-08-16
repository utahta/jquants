package client

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError is returned when the J-Quants API responds with a status other than
// 200 OK. It carries the raw status code and response body so that callers can
// branch on the status instead of parsing the error message.
//
//	var apiErr *client.APIError
//	if errors.As(err, &apiErr) {
//		log.Printf("status=%d body=%s", apiErr.StatusCode, apiErr.Body)
//	}
//
// Service methods wrap the error with %w, so errors.As works on the error
// returned by any service method.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int
	// Body is the raw response body. The API returns a JSON object such as
	// {"message": "..."} for most errors, but this is not guaranteed.
	Body string
}

// Error implements the error interface. The format is fixed for backward
// compatibility with callers that match on the message.
func (e *APIError) Error() string {
	return fmt.Sprintf("API error: status=%d, body=%s", e.StatusCode, e.Body)
}

// StatusCode returns the HTTP status code carried by err. ok is false when err
// is not (and does not wrap) an *APIError, which is the case for transport,
// context and decoding errors.
func StatusCode(err error) (code int, ok bool) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.StatusCode, true
	}
	return 0, false
}

// IsRateLimitExceeded reports whether err is an API error caused by exceeding
// the rate limit (429). Wait before retrying: repeating requests while over the
// limit can block access for about five minutes.
func IsRateLimitExceeded(err error) bool {
	code, ok := StatusCode(err)
	return ok && code == http.StatusTooManyRequests
}

// IsAuthError reports whether err is an API error caused by authentication or
// authorization (401, 403). The API returns 403 for an invalid API key as well
// as for data that is not included in the contracted plan, so retrying does not
// help.
func IsAuthError(err error) bool {
	code, ok := StatusCode(err)
	return ok && (code == http.StatusUnauthorized || code == http.StatusForbidden)
}

// IsServerError reports whether err is an API error with a 5xx status. These
// are usually temporary and can be retried after a delay.
func IsServerError(err error) bool {
	code, ok := StatusCode(err)
	return ok && code >= 500 && code <= 599
}
