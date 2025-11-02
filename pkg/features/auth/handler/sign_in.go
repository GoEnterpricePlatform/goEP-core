package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/features/auth/core"
	cShared "github.com/amorindev/go-tmpl/pkg/shared/api/core"
	dShared "github.com/amorindev/go-tmpl/pkg/shared/domain"
)

// SignIn handles user authentication requests
func (h Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req core.SignInReq

	// Decode JSON request body into SignInReq struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	// Validate the sign in request
	err = req.IsSignInValid()
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	user, session, err := h.AuthSrv.SignIn(r.Context(), req.Email, req.Password, req.RememberMe)
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	cookie := &http.Cookie{
		Name:     "refreshToken",
		Value:    session.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(session.RefreshTokenExpIn),
	}

	if h.AppEnv == "prod" {
		cookie.Secure = true
		cookie.SameSite = http.SameSiteNoneMode
	} else {
		cookie.Secure = false
		cookie.SameSite = http.SameSiteLaxMode
	}

	session.RefreshToken = ""

	// create response
	resp := core.NewSignInResp(user, session)

	http.SetCookie(w, cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
