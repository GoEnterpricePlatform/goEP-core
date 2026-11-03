package handler

import (
	"context"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/claim"
	identityPages "github.com/GoEnterpricePlatform/goEP-core/addons/identity/ui/pages"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/api/middlewares"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/internal/shared/front/core"
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

	h.renderDashboard(w, r, claims, "roles-permissions", identityPages.RolesPermissionsContent(roles))
}
