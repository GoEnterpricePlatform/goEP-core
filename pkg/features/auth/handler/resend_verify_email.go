package handler

import (
	"encoding/json"
	"net/http"

	"github.com/amorindev/go-tmpl/pkg/features/auth/core"
	cShared "github.com/amorindev/go-tmpl/pkg/shared/api/core"
	dShared "github.com/amorindev/go-tmpl/pkg/shared/domain"
)

// ResendVerifyEmail handles the HTTP request to resend a verification email to the user.
// It expects a JSON body containing the user's email address.
func (h Handler) ResendVerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req core.EmailReq

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	err = req.IsEmailValid()
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	otpID, err := h.AuthSrv.ResendVerifyEmail(r.Context(), req.Email)
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	resp := core.OtpIDResp{
		OtpID: otpID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
