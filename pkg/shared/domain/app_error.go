package domain

import (
	"errors"
	"fmt"
	"log"
)

// AppError is a custom error type that implements the error interface
type AppError struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// NewAppError creates a new AppError with the given code and message.
func NewAppError(code string, msg string) AppError {
	return AppError{
		Code: code,
		Msg:  msg,
	}
}

// Error returns a string representation of the error. It is part of the error interface.
func (e AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Msg)
}

func ManageError(err error, msg string) error {
	var appErr AppError

	switch {
	case errors.Is(err, ErrDuplicateKey):
		log.Println("duplicate key")
		appErr = AppError{
			Code: ErrCodeDuplicateKey,
			Msg:  "Duplicate key",
		}
	case errors.Is(err, ErrIncorrectID):
		log.Println("incorrect id error")
		appErr = AppError{
			Code: ErrCodeInvalidParams,
			Msg:  "Incorrect id",
		}
	case errors.Is(err, ErrNotFound):
		log.Println("not found error")
		appErr = AppError{
			Code: ErrCodeNotFound,
			Msg:  "Not found",
		}
	case errors.Is(err, ErrTimeout):
		log.Println("timeout error")
		appErr = AppError{
			Code: ErrCodeTimeout,
			Msg:  "Timeout",
		}
	case errors.Is(err, ErrTokenExpired):
		log.Println("token expired")
		appErr = AppError{
			Code: ErrCodeUnauthorized,
			Msg:  "Token expired",
		}
	case errors.Is(err, ErrTokenSignature):
		log.Println("invalid token signature")
		appErr = AppError{
			Code: ErrCodeUnauthorized,
			Msg:  "Invalid token signature",
		}
	case errors.Is(err, ErrTokenMalformed):
		log.Println("malformed token")
		appErr = AppError{
			Code: ErrCodeInvalidParams,
			Msg:  "Malformed token",
		}
	case errors.Is(err, ErrTokenInvalid):
		log.Println("invalid token")
		appErr = AppError{
			Code: ErrCodeUnauthorized,
			Msg:  "Invalid token",
		}
	case errors.Is(err, ErrTokenInvalidClaim):
		log.Println("invalid token claims")
		appErr = AppError{
			Code: ErrCodeInvalidParams,
			Msg:  "Invalid token claims",
		}
	case errors.Is(err, ErrAuthHeaderMissing):
		log.Println("missing authorization header")
		appErr = AppError{
			Code: ErrCodeUnauthorized,
			Msg:  "Authorization header is required",
		}
	case errors.Is(err, ErrAuthHeaderInvalid):
		log.Println("invalid authorization header format")
		appErr = AppError{
			Code: ErrCodeUnauthorized,
			Msg:  "invalid authorization header format",
		}
	case errors.Is(err, ErrAccountInactive):
		log.Println("account is inactive")
		appErr = AppError{
			Code: ErrCodeUnauthorized,
			Msg:  "account is inactive",
		}
	case errors.Is(err, ErrPassDoNotMatch):
		log.Println("password does not match")
		appErr = AppError{
			Code: ErrCodeUnauthorized,
			Msg:  "password does not match",
		}
	case errors.Is(err, ErrInvalidOtpCode):
		log.Println("invalid otp code")
		appErr = AppError{
			Code: ErrCodeInvalidParams,
			Msg:  "Invalid OTP code",
		}

	case errors.Is(err, ErrOtpPurposeNotAllowed):
		log.Println("otp purpose not allowed")
		appErr = AppError{
			Code: ErrCodeInvalidParams,
			Msg:  "OTP not valid for this purpose",
		}
	case errors.Is(err, ErrOtpExpired):
		log.Println("otp expired")
		appErr = AppError{
			Code: ErrCodeInvalidParams,
			Msg:  "OTP code has expired",
		}

	default:
		log.Println(err.Error())
		appErr = AppError{
			Code: ErrCodeInternalServerError,
			Msg:  "Server Error",
		}
	}

	// We only add the custom message if the error is not an internal server error
	if msg != "" && appErr.Code != ErrCodeInternalServerError {
		appErr.Msg = fmt.Sprintf("%s: %s", appErr.Msg, msg)
	}
	return appErr
}
