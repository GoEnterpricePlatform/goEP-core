package logger

import "net/http"

// custom response writer to capture status
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Preserve http.Flusher for SSE (Datastar) 
func (rw *responseWriter) Flush() { 
	if f, ok := rw.ResponseWriter.(http.Flusher); ok { 
		f.Flush() 
	} 
}


