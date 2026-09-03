package middlewares

import (
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/claim"
)

func (m *MdwSrvTmpl) RequirePermission(permission string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := r.Context().Value(AccessTokenClaimsTmplIDKey).(*claim.AccessTokenClaims)
		if !ok || claims == nil {
			http.Redirect(
				w,
				r,
				"/v1/admin/auth/sign-in",
				http.StatusFound,
			)
			return
		}

		// ROLE_SYSTEM_ADMIN has full access.
		if HasRole(claims.Roles, "ROLE_SYSTEM_ADMIN") {
			next.ServeHTTP(w, r)
			return
		}

		if !HasPermission(claims.Permissions, permission) {
			http.Redirect(
				w,
				r,
				"/v1/goep-admin/access-denied",
				http.StatusFound,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
