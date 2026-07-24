package middleware

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime/debug"
)

// Recoverer is a middleware that recovers from panics, logs the panic (and a
// backtrace), and returns a HTTP 500 (Internal Server Error) status if
// possible. Recoverer prints a pretty backtrace (using colors) to stderr.
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil && rvr != http.ErrAbortHandler {
				logEntry := GetLogEntry(r)

				var written bool
				if ww, ok := w.(WrapResponseWriter); ok {
					written = ww.Written()
				}

				if r.Context().Err() != nil || written {
					return
				}

				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					PrintPrettyStack(rvr)
				}

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
