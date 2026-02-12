package handler

import (
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/auth/core"
	sharedC "github.com/amorindev/go-cms-tmpl/web/shared/core"
	"github.com/amorindev/go-cms-tmpl/web/shared/templates"
)

// SignIn handles the admin sign-in request from the form,
// validates credentials, creates the session cookie (refresh token),
// and redirects the user to the admin home.
func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	remember := r.FormValue("remember") == "true"

	req := &core.SignInReq{
		Email:      r.FormValue("email"),
		Password:   r.FormValue("password"),
		RememberMe: remember,
	}

	err := req.IsSignInValid()
	if err != nil {
		h.AdminRenderer.Render(w, "sign-in", templates.ErrorData{
			ErrorMsg: sharedC.UiErrorResp(err),
		})
		return
	}

	_, session, err := h.AdminSrv.SignIn(r.Context(), req.Email, req.Password, remember)
	if err != nil {
		h.AdminRenderer.Render(w, "sign-in", templates.ErrorData{
			ErrorMsg: sharedC.UiErrorResp(err),
		})
		return
	}

	// In the admin panel, email verification is not required,
	// so a valid session should always exist after sign-in.
	// If email verification is added later and the session is missing,
	// redirect the user to the email verification flow.
	accessCookie := h.CookieSrv.CreateForTemplate("accessToken", session.AccessToken, int(session.AccessTokenExpIn))
	http.SetCookie(w, accessCookie)

	refreshCookie := h.CookieSrv.CreateForTemplate("refreshToken", session.RefreshToken, int(session.RefreshTokenExpIn))
	http.SetCookie(w, refreshCookie)

	http.Redirect(w, r, "/v1/admin/home", http.StatusFound)
}
