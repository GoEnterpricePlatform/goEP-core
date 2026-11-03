package port

import "net/http"

// Authenticator is the authentication contract consumed by backend modules.
type Authenticator interface {
	AccessTokenMdw(next http.HandlerFunc) http.HandlerFunc
}
