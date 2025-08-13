package middleware

import (
	"fmt"
	"lotest/internal/logger"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func LogHandler(next http.HandlerFunc, logCh chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: w, Status: http.StatusOK}

		next.ServeHTTP(recorder, r)

		logCh <- logger.FormatLog(logger.Info,
			fmt.Sprintf("Method: %s, Path: %s, Status: %d, Duration: %s",
				r.Method,
				r.URL.Path,
				recorder.Status,
				time.Since(start)),
		)
	}
}
