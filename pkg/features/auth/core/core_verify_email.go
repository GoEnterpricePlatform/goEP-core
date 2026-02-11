package core

import sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"

// VerifyEmailReq represents the request payload for verifying a user's email address.
// It is used after registration when the user receives an OTP code via email.
// This struct is handled by the endpoint that validates the OTP, whether the email
// was sent during the signUp flow or through the resend verification feature.
type VerifyEmailReq struct {
	OtpID   string `json:"otp_id"`
	OtpCode string `json:"otp_code"`
	UserID  string `json:"user_id"`
}

func (req *VerifyEmailReq) IsVerifyEmailOTPValid() error {
	if req.OtpID == "" {
		return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "otp_id field is required")
	}
	if req.OtpCode == "" {
		return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "otp_code field is required")
	}

	if req.UserID == "" {
		return sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "user_id is required")
	}
	return nil
}
