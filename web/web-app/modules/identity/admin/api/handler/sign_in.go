package handler

import (
	"net/http"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
	"github.com/starfederation/datastar-go/datastar"
)

// SignIn handles the admin sign-in request from the form,
// validates credentials, creates the session cookie (refresh token),
// and redirects the user to the goep-admin .
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseMultipartForm(10 << 20) // 10MB maximum
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"error": sharedC.UiErrorResp(sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid form")),
		})
		return
	}

	rememberMe := r.FormValue("remember_me")
	// Convert rememberMe to boolean
	rememberMeBool := false
	if rememberMe == "true" || rememberMe == "on" {
		rememberMeBool = true
	}

	req := &core.SignInReq{
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		RememberMe: rememberMeBool,
	}

	err = req.IsSignInValid()
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"error": sharedC.UiErrorResp(err),
		})
		return
	}

	_, session, err := h.AdminSrv.SignIn(r.Context(), req.Email, req.Password, rememberMeBool)
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"error": sharedC.UiErrorResp(err),
		})
		return
	}

	// In the admin panel, email verification is not required,
	// so a valid session should always exist after sign-in.
	// If email verification is added later and the session is missing,
	// redirect the user to the email verification flow.
	h.CookieSrv.SetAccessToken(w, session.AccessToken)

	h.CookieSrv.SetRefreshToken(w, session.RefreshToken, time.Duration(session.RefreshTokenExpIn))

	// redirect to goep-admin
	sse := datastar.NewSSE(w, r)
	sse.Redirect("/v1/goep-admin")
}
