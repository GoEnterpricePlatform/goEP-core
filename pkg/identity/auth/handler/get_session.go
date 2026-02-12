package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/auth/core"
	sharedC "github.com/amorindev/go-cms-tmpl/pkg/shared/api/core"
	"github.com/amorindev/go-cms-tmpl/pkg/shared/api/handler"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

// Restore the user session using the refresh token, validate it, create a new session, and avoid returning the refresh token in the response body.
func (h Handler) GetSession(w http.ResponseWriter, r *http.Request) {
    // Try to read the refresh token cookie
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		w.Header().Set("Content-Type","application/json")
        w.WriteHeader(http.StatusOK)
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
	user,session, err := h.AuthSrv.GetSession(context.Background(), c.ID, c.UserID)
	if err != nil {
		handler.ClearCookie(w, "refreshToken", h.AppEnv)
		sharedC.RespondError(w, err)
		return
	}

    // Never return a refresh token to the client
	session.RefreshToken = ""

    // create response
	resp := core.NewSignInResp(user, session,"")


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
