package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/features/auth/core"
	sharedC "github.com/amorindev/go-tmpl/pkg/shared/api/core"
	sharedD "github.com/amorindev/go-tmpl/pkg/shared/domain"
)

func (h Handler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req core.VerifyEmailReq

	// Decode JSON request body into SignUpReq struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	err = req.IsVerifyEmailOTPValid()
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}

	user, session, err := h.AuthSrv.VerifyEmail(r.Context(), req.OtpID, req.OtpCode, req.UserID)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}

	resp := core.NewSignInResp(user, session)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
