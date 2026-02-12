package handler

import (
	"context"
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/auth/core"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/users/domain"
	sharedC "github.com/amorindev/go-cms-tmpl/web/shared/core"
	"github.com/amorindev/go-cms-tmpl/web/shared/templates"
)

// SignUp handles the admin signup form submission.
// It validates the input, creates a new admin user, and redirects to the sign-in page on success.
// If any validation or business error occurs, the signup page is re-rendered with an error message.
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	req := &core.SignUpReq{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	err := req.IsSignUpValid()
	if err != nil {
		h.AdminRenderer.Render(w, "sign-up", templates.ErrorData{
			ErrorMsg: sharedC.UiErrorResp(err),
		})
		return
	}

	// Create a new user domain object
	user := domain.NewUser(req.Email, req.Password)

	err = h.AdminSrv.SignUp(context.Background(), user)
	if err != nil {
		h.AdminRenderer.Render(w, "sign-up", templates.ErrorData{
			ErrorMsg: sharedC.UiErrorResp(err),
		})
		return
	}

	// redirect to SignIn
	http.Redirect(w, r, "/v1/admin/sign-in", http.StatusSeeOther)
}
