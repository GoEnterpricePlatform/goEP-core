package logger

import (
	"fmt"
	"net/http"
	"time"
)

// ANSI colors
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
)

// colors by status
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return Green
	case code >= 300 && code < 400:
		return Blue
	case code >= 400 && code < 500:
		return Yellow
	default:
		return Red
	}
}

// colors by method
func colorForMethod(method string) string {
	switch method {
	case "GET":
		return Cyan
	case "POST":
		return Yellow
	case "PUT":
		return Blue
	case "DELETE":
		return Red
	default:
		return Reset
	}
}

// Middleware
func (l *HttpLogger) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(rw, r)

		// Ignore OPTIONS requests (preflight CORS)
		if r.Method == http.MethodOptions {
			return
		}

		latency := time.Since(start)
		status := rw.statusCode
		statusColor := colorForStatus(status)
		methodColor := colorForMethod(r.Method)

		msg := fmt.Sprintf("%s| %s%d%s | %s%s%s %s | %s | %v",
			Reset, statusColor, status, Reset,
			methodColor, r.Method, Reset,
			r.URL.Path,
			latency,
			r.RemoteAddr,
		)

		l.log.Info(msg)
	})
}
