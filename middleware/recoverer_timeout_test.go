package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRecovererTimeout(t *testing.T) {
	h := Timeout(10 * time.Millisecond)(Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		panic("test panic post-timeout")
	})))

	r, _ := http.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()

	h.ServeHTTP(w, r)

	// Wait for the orphaned goroutine to finish and panic
	time.Sleep(100 * time.Millisecond)

	if w.Code != http.StatusGatewayTimeout {
		t.Errorf("expected status 504, got %d", w.Code)
	}
}
