package middleware

import (
	"net/http"
	"time"
)

// Timeout is a middleware that cancels ctx after dt duration and writes
// a 504 Gateway Timeout response.
func Timeout(dt time.Duration) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.TimeoutHandler(next, dt, "Gateway Timeout")
	}
}
