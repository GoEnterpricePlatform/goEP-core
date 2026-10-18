package handler

import (
	"net/http"
)

// requirePaddleEnabled prevents Paddle routes from running when Paddle billing is disabled in the application configuration.
func (h *Handler) requirePaddleEnabled(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !h.IsEnabled {
			http.Error(w, "Paddle is not enabled", http.StatusServiceUnavailable)
			return
		}
		next(w, r)
	}
}
