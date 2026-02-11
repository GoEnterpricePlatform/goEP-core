package service_test

import (
	"testing"
	"time"

	"github.com/amorindev/go-cms-tmpl/internal/tokens/service"
	"github.com/stretchr/testify/require"
)

func TestService_CreateRefreshToken(t *testing.T) {
	s := service.NewTokenSrv("a_secret", "r_secret", time.Minute, time.Hour, 24*time.Minute, "my-app")

	testTable := map[string]struct {
		userID        string
		RememberMe    bool
		assertionFunc func(subTest *testing.T, gotClaimID string, gotToken string, gotExpIn int64, gotErr error)
	}{
		"success": {
			userID:     "123",
			RememberMe: false,
			assertionFunc: func(subTest *testing.T, gotClaimID string, gotToken string, gotExpIn int64, gotErr error) {
				require.NoError(subTest, gotErr)
				require.NotEmpty(subTest, gotClaimID)
				require.NotEmpty(subTest, gotToken)
				require.Equal(subTest, int64(s.RefreshExpiresIn.Seconds()), gotExpIn)
			},
		},
		"success with rememberMe": {
			userID:     "123",
			RememberMe: true,
			assertionFunc: func(subTest *testing.T, gotClaimID, gotToken string, gotExpIn int64, gotErr error) {
				require.NoError(subTest, gotErr)
				require.NotEmpty(subTest, gotClaimID)
				require.NotEmpty(subTest, gotToken)
				require.Equal(subTest, int64(s.RefreshExpiresInRemember.Seconds()), gotExpIn)
			},
		},
	}

	for testName, test := range testTable {
		t.Run(testName, func(subTest *testing.T) {
			gotClaimID, gotToken, gotExpIn, gotErr := s.CreateRefreshToken(test.userID, test.RememberMe)
			test.assertionFunc(subTest, gotClaimID, gotToken, gotExpIn, gotErr)
		})
	}
}
