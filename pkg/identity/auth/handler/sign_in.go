package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/core"
	sharedH "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/handler"
	dShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	cShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
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

	user, session, otpID, err := h.AuthSrv.SignIn(r.Context(), req.Email, req.Password, req.RememberMe)
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	if session != nil {
		cookie := sharedH.CreateCookie(session.RefreshToken, int(session.RefreshTokenExpIn), h.AppEnv)
		session.RefreshToken = ""
		http.SetCookie(w, cookie)
	}

	// create response
	resp := core.NewSignInResp(user, session, otpID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
