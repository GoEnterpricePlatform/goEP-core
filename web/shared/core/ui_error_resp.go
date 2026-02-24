package core

import (
	"errors"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

const InternalServerErrorMessage = "Ooops! Something went wrong. Please help us by reporting this issue."

// Alternatively, you could send only the appErr.msg
func UiErrorResp(err error) string {
	var appErr domain.AppError
	if errors.As(err, &appErr) {
		return appErr.Error()
	}
	return InternalServerErrorMessage
}
