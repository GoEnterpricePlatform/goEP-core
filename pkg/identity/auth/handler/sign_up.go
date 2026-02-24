package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
	cShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	dShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// Signup handles user registration, validates input, creates a new user, and returns a JSON response.
func (h Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req core.SignUpReq

	// Decode JSON request body into SignUpReq struct
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInvalidParams, "invalid request body"))
		return
	}

	defer r.Body.Close()

	// Validate the sign-up request
	err = req.IsSignUpValid()
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	// Create a new user domain object
	user := domain.NewUser(req.Email, req.Password)

	otpID, err := h.AuthSrv.SignUp(context.Background(), user)
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	// Create response from the created user domain
	resp := core.NewSignUpResp(user, otpID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
