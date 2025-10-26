package domain

import "errors"

const (
	ErrCodeDuplicateKey        = "duplicate_key"
	ErrCodeInternalServerError = "internal_server_error"
	ErrCodeInvalidParams       = "invalid_params"
	ErrCodeNotFound            = "not_found"
	ErrCodeTimeout             = "timeout"
	ErrCodeUnauthorized        = "unauthorized"
	ErrCodeForbidden           = "forbidden"
)

var (
	ErrDuplicateKey = errors.New("duplicate key error")
	ErrIncorrectID  = errors.New("incorrect id error")
	ErrNotFound     = errors.New("record not found error")
	ErrTimeout      = errors.New("timeout error")
)

var (
	ErrTokenExpired         = errors.New("token expired")
	ErrTokenSignature       = errors.New("invalid signature")
	ErrTokenMalformed       = errors.New("malformed token")
	ErrTokenInvalid         = errors.New("invalid token")
	ErrTokenInvalidClaim    = errors.New("invalid claims")
	ErrAuthHeaderMissing    = errors.New("authorization header is required")
	ErrAuthHeaderInvalid    = errors.New("invalid authorization format")
	ErrRefreshTokenRequired = errors.New("refresh token is required")
	ErrInvalidRequestBody   = errors.New("invalid request body")
)

var (
	ErrAccountInactive = errors.New("account inactive")
	ErrPassDoNotMatch  = errors.New("passwords do not match")
)

var (
	ErrInvalidOtpCode       = errors.New("invalid otp code")
	ErrOtpPurposeNotAllowed = errors.New("otp purpose not allowed")
	ErrOtpExpired           = errors.New("otp expired")
)
