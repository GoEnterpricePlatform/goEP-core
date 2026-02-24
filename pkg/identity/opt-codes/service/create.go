package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/identity/opt-codes/domain"
)

// Create generates a 6-digit one-time password (OTP) as a zero-padded string.
// The code is cryptographically secure, randomly generated in the range 000000–999999.
func (s *Service) Create(ctx context.Context, otp *domain.OtpCode) error {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return err
	}

	otp.OptCode = fmt.Sprintf("%06d", nBig.Int64())

	err = s.OtpCodeRepo.Insert(ctx, otp)
	if err != nil {
		return err
	}

	return nil
}
