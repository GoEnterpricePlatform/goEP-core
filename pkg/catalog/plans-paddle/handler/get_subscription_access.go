package handler

import (
	"encoding/json"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/claim"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/core"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/api/middlewares"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/domain"
	"net/http"
)

// GetSubscriptionAccess handles GET /paddle/subscription/access when the app needs to decide whether the signed-in user may use subscription-gated features, such as after login or before opening a protected screen.
// It takes the user ID from the access token and checks the local subscription cache; only active or trialing subscriptions grant access. It does not call Paddle: webhook processing keeps this local state current. If no eligible subscription exists, the response has has_access=false.
func (h *Handler) GetSubscriptionAccess(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.AccessTokenClaimsIDKey).(*claim.AccessTokenClaims)
	if !ok || claims.UserID == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeUnauthorized, "authentication is required"))
		return
	}
	access, err := h.PaddlePlanSrv.GetSubscriptionAccess(r.Context(), claims.UserID)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(access)
}
