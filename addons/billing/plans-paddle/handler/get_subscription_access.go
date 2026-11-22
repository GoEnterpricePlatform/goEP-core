package handler

import (
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/tokens/claim"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/core"
	"github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/api/middlewares"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/internal/shared/back/domain"
)

func (h *Handler) GetSubscriptionAccess(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(middlewares.AccessTokenClaimsIDKey).(*claim.AccessTokenClaims)
	if !ok || claims.UserID == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeUnauthorized, "authentication is required"))
		return
	}
	access, err := h.PaddleBillingSrv.GetSubscriptionAccess(r.Context(), claims.UserID)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(access)
}
