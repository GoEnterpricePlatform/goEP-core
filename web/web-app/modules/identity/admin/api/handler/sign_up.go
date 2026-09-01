package handler

import (
	"context"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/auth/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/users/domain"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
	"github.com/starfederation/datastar-go/datastar"
)

// SignUp handles the admin signup form submission.
// It validates the input, creates a new admin user, and redirects to the sign-in page on success.
// If any validation or business error occurs, the signup page is re-rendered with an error message.
func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	err := r.ParseMultipartForm(10 << 20) // 10MB maximum
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"error": sharedC.UiErrorResp(sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid form")),
		})
		return
	}

	req := &core.SignUpReq{
		Email:           r.FormValue("email"),
		Password:        r.FormValue("password"),
		ConfirmPassword: r.FormValue("confirm_password"),
	}

	err = req.IsSignUpValid()
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"error": sharedC.UiErrorResp(err),
		})
		return
	}

	// Create a new user domain object
	user := domain.NewUser(req.Email, req.Password)

	err = h.AdminSrv.SignUp(context.Background(), user)
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"error": sharedC.UiErrorResp(err),
		})
		return
	}

	// redirect to SignIn
	sse := datastar.NewSSE(w, r)
	sse.Redirect("/v1/admin/auth/sign-in")
}
