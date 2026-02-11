package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/features/auth/core"
	sharedC "github.com/amorindev/go-cms-tmpl/pkg/shared/api/core"
	"github.com/amorindev/go-cms-tmpl/pkg/shared/api/handler"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

// RefreshToken renews the user's session.
// Reads and validates the "refreshToken" cookie.
// Calls AuthSrv.RefreshToken to verify the session, remove the old one, and create a new one.
// Updates the cookie with the new refresh token.
// Clears the refresh token from the response and returns the new session.
func (h Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	// Try to read the refresh token cookie
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		handler.ClearCookie(w, "refreshToken", h.AppEnv)
		if errors.Is(err, http.ErrNoCookie) {
			sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "refresh token cookie not found"))
			return
		}
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInternalServerError, "failed to read refresh token cookie"))
		return
	}

	// Validate and parse the refresh token
	c, err := h.TokenSrv.ParseRefreshToken(cookie.Value)
	if err != nil {
		handler.ClearCookie(w, "refreshToken", h.AppEnv)
		sharedC.RespondError(w, sharedD.ManageError(err, ""))
		return
	}

	// Call business logic
	session, err := h.AuthSrv.RefreshToken(context.Background(), c.ID, c.UserID)
	if err != nil {
		handler.ClearCookie(w, "refreshToken", h.AppEnv)
		sharedC.RespondError(w, err)
		return
	}

	// Update cookie if a new refresh token was issued
	newCookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    session.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(session.RefreshTokenExpIn),
	}

	if h.AppEnv == "prod" {
		newCookie.Secure = true
		newCookie.SameSite = http.SameSiteNoneMode
	} else {
		newCookie.Secure = false
		newCookie.SameSite = http.SameSiteLaxMode
	}

	session.RefreshToken = ""

	// create response
	resp := core.NewRefreshTokenResp(session)

	http.SetCookie(w, newCookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
