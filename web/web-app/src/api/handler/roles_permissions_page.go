package handler

import (
	"context"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/claim"
	"github.com/GoEnterpricePlatform/goEP-core/web/shared/api/middlewares"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/src/ui/pages"
	"github.com/starfederation/datastar-go/datastar"
)

func (h *Handler) RolesPermissionsPage(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(
		middlewares.AccessTokenClaimsTmplIDKey,
	).(*claim.AccessTokenClaims)

	if !ok || claims == nil {
		http.Redirect(
			w,
			r,
			"/v1/goep-admin/auth/sign-in",
			http.StatusFound,
		)
		return
	}

	roles, err := h.RoleSrv.GetAll(context.Background())
	if err != nil {
		sse := datastar.NewSSE(w, r)
		sse.MarshalAndPatchSignals(map[string]any{
			"error": sharedC.UiErrorResp(err),
		})
		return
	}

	err = pages.RolesPermissionsPage(claims, roles).Render(r.Context(), w)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
