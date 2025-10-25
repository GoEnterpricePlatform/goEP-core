package resend

import (
	"github.com/resend/resend-go/v2"
)

func NewResendClient(apiKey string) *resend.Client {
	return resend.NewClient(apiKey)
}
