package service_test

import (
	"testing"
	"time"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/permissions/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/tokens/service"
	"github.com/stretchr/testify/require"
)

func TestService_CreateAccessToken(t *testing.T) {

	s := service.NewTokenSrv("a_secret", "r_secret", time.Minute, time.Hour, 24*time.Minute, "my-app")

	testTable := map[string]struct {
		userID      string
		email       string
		roles       []string
		permissions []string
	}{
		"success": {
			userID:      "123",
			email:       "user@example.com",
			roles:       []string{"User"},
			permissions: []string{string(domain.PAdminAccess)},
		},
		"success without roles and permissions": {
			userID:      "123",
			email:       "user@example.com",
			roles:       nil,
			permissions: nil,
		},
	}

	for testName, test := range testTable {
		t.Run(testName, func(subTest *testing.T) {
			gotToken, gotExp, gotErr := s.CreateAccessToken(test.userID, test.email, test.permissions, test.roles)
			require.NoError(subTest, gotErr)
			require.NotEmpty(subTest, gotToken)
			require.Equal(subTest, int64(s.AccessExpiresIn.Seconds()), gotExp)
		})
	}
}
