package middleware

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aige/requestlog"
)

type responseWriter struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func RequestLogger(repo requestlog.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			var reqBody []byte
			if r.Body != nil {
				reqBody, _ = io.ReadAll(r.Body)
				r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
			}

			rw := &responseWriter{ResponseWriter: w, body: &bytes.Buffer{}, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)
			respBody := truncate(rw.body.String(), 500)
			reqStr := truncate(string(reqBody), 500)

			log.Printf("[API] %s %s %d %v | req=%s | resp=%s",
				r.Method, r.URL.Path, rw.statusCode, duration,
				reqStr, respBody,
			)

			// Persist request log to database (fire and forget)
			go func() {
				userID := ""
				adminID := ""
				if uid := GetUserID(r.Context()); uid != "" {
					userID = uid
				}
				if aid := GetAdminID(r.Context()); aid != "" {
					adminID = aid
				}
				if err := repo.Insert(context.Background(), &requestlog.RequestLog{
					Method:       r.Method,
					Path:         r.URL.Path,
					StatusCode:   rw.statusCode,
					DurationMs:   int(duration.Milliseconds()),
					RequestBody:  reqStr,
					ResponseBody: respBody,
					UserID:       userID,
					AdminID:      adminID,
					IPAddress:    r.RemoteAddr,
				}); err != nil {
					log.Printf("WARN failed to persist request log: %v", err)
				}
			}()
		})
	}
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}
