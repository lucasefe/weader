package util

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"
)

// TimerHandler is the handler type for the time recorder
type TimerHandler struct {
	handler http.Handler
}

// NewTimer builds a new Timer
func NewTimer(handler http.Handler) *TimerHandler {
	return &TimerHandler{handler: handler}
}

func (s *TimerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rec := httptest.NewRecorder()

	defer func(start time.Time) {
		for k, v := range rec.Header() {
			w.Header()[k] = v
		}

		w.Header().Set("X-Runtime", fmt.Sprintf("%v", time.Since(start)))
		w.Write(rec.Body.Bytes())
	}(time.Now())

	s.handler.ServeHTTP(rec, r)
}
