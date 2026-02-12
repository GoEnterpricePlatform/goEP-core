package service_test

import (
	"testing"
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/identity/tokens/service"
	"github.com/stretchr/testify/require"
)


func TestService_CreateAccessToken(t *testing.T) {

	s := service.NewTokenSrv("a_secret", "r_secret", time.Minute, time.Hour, 24*time.Minute, "my-app")

	testTable := map[string]struct {
		userID string
		email  string
		roles  []string
	}{
		"success": {
			userID: "123",
			email:  "user@example.com",
			roles:  []string{"User"},
		},
		"success without roles": {
			userID: "123",
			email:  "user@example.com",
			roles:  nil,
		},
	}

	for testName, test := range testTable {
		t.Run(testName, func(subTest *testing.T) {
			gotToken, gotExp, gotErr := s.CreateAccessToken(test.userID, test.email, test.roles)
			require.NoError(subTest, gotErr)
			require.NotEmpty(subTest, gotToken)
			require.Equal(subTest, int64(s.AccessExpiresIn.Seconds()), gotExp)
		})
	}
}
