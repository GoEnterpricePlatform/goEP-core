package core

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

const InternalServerErrorMessage = "Ooops! Something went wrong. Please help us by reporting this issue."

var ErrCodeMapping = map[string]int{
	domain.ErrCodeDuplicateKey:  http.StatusConflict,
	domain.ErrCodeNotFound:      http.StatusNotFound,
	domain.ErrCodeInvalidParams: http.StatusBadRequest,
	domain.ErrCodeUnauthorized:  http.StatusUnauthorized,
	domain.ErrCodeForbidden:     http.StatusForbidden,
}

func RespondError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	var appErr domain.AppError
	if errors.As(err, &appErr) {
		if status, ok := ErrCodeMapping[appErr.Code]; ok {
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(appErr)
			return
		}
	}

	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(domain.AppError{
		Code: domain.ErrCodeInternalServerError,
		Msg:  InternalServerErrorMessage,
	})
}
