package handler

import (
	"errors"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// SignOut handles user logout by removing the refresh token both
// from the client (cookie) and from the database.
func (h Handler) SignOut(w http.ResponseWriter, r *http.Request) {
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

	// Remove the refresh token from the database to invalidate the session
	err = h.AuthSrv.SignOut(r.Context(), c.ID)
	if err != nil {
		handler.ClearCookie(w, "refreshToken", h.AppEnv)
		sharedC.RespondError(w, err)
		return
	}

	// Clear the cookie and return success response
	handler.ClearCookie(w, "refreshToken", h.AppEnv)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
