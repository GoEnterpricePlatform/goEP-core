package service_test

import (
	"testing"
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/claim"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/port"
	"github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/service"
	"github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestService_ParseRefreshToken(t *testing.T) {
	sharedSrv := service.NewTokenSrv("a_secret", "r_secret", time.Minute, time.Hour, 24*time.Minute, "my-app")

	testTable := map[string]struct {
		setup         func() (port.TokenSrv, string)
		assertionFunc func(subTest *testing.T, gotClaims *claim.RefreshTokenClaims, gotErr error)
	}{
		"success": {
			setup: func() (port.TokenSrv, string) {
				_, token, _, err := sharedSrv.CreateRefreshToken("123", false)
				require.NoError(t, err)
				return sharedSrv, token
			},
			assertionFunc: func(subTest *testing.T, gotClaims *claim.RefreshTokenClaims, gotErr error) {
				require.NoError(subTest, gotErr)
				require.Equal(subTest, "123", gotClaims.UserID)
				require.Equal(subTest, "123", gotClaims.Subject)
				require.NotEmpty(subTest, gotClaims.ID)
			},
		},
		"expired token": {
			setup: func() (port.TokenSrv, string) {
				shortS := service.NewTokenSrv("a_secret", "r_secret", 1*time.Second, 1*time.Second, 1*time.Second, "my-app")
				_, token, _, err := shortS.CreateRefreshToken("123", false)
				require.NoError(t, err)

				time.Sleep(2 * time.Second)
				return shortS, token
			},
			assertionFunc: func(subTest *testing.T, gotClaims *claim.RefreshTokenClaims, gotErr error) {
				require.Nil(subTest, gotClaims)
				require.Error(subTest, gotErr)
				require.ErrorIs(subTest, domain.ErrTokenExpired, gotErr)
			},
		},
		"invalid signature": {
			setup: func() (port.TokenSrv, string) {
				badSrv := service.NewTokenSrv("a_secret", "r_wrong", time.Minute, time.Hour, 24*time.Minute, "my-app")
				_, token, _, err := badSrv.CreateRefreshToken("123", false)
				require.NoError(t, err)
				return sharedSrv, token
			},
			assertionFunc: func(subTest *testing.T, gotClaims *claim.RefreshTokenClaims, gotErr error) {
				require.Nil(subTest, gotClaims)
				require.Error(subTest, gotErr)
				require.ErrorIs(subTest, domain.ErrTokenSignature, gotErr)

			},
		},
		"malformed token or empty": {
			setup: func() (port.TokenSrv, string) {
				return sharedSrv, "a.a.a"
			},
			assertionFunc: func(subTest *testing.T, gotClaims *claim.RefreshTokenClaims, gotErr error) {
				require.Nil(subTest, gotClaims)
				require.Error(subTest, gotErr)
				require.ErrorIs(subTest, domain.ErrTokenMalformed, gotErr)
			},
		},
		"token not valid (nbf)": {
			setup: func() (port.TokenSrv, string) {
				c := claim.NewRefreshTokenClaim("123", time.Minute)
				// Force the NotBefore (nbf) field to one hour in the future.
				// The token should not be considered valid until that time.
				c.NotBefore = jwt.NewNumericDate(time.Now().Add(time.Hour))
				token, _ := c.GetToken("r_secret")
				return sharedSrv, token
			},
			assertionFunc: func(subTest *testing.T, gotClaims *claim.RefreshTokenClaims, gotErr error) {
				require.Nil(subTest, gotClaims)
				require.Contains(subTest, gotErr.Error(), "invalid refresh token")
			},
		},
	}

	for testName, test := range testTable {
		t.Run(testName, func(subTest *testing.T) {
			srv, token := test.setup()

			gotClaims, gotErr := srv.ParseRefreshToken(token)
			test.assertionFunc(subTest, gotClaims, gotErr)
		})
	}
}
