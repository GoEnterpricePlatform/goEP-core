package handler

import (
	"context"
	"io"
	"net/http"

	cShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	dShared "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h Handler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("userId")

	if userID == "" {
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInvalidParams, "missing user id"))
		return
	}

	err := r.ParseMultipartForm(20 << 20) // 20MB
	if err != nil {
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInvalidParams, "invalid form"))
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInvalidParams, "error retrieving file"))
		return
	}
	defer file.Close()

	// Proper Content-Type detection
	fileBuffer := make([]byte, 512)
	_, err = file.Read(fileBuffer)
	if err != nil {
		msg := "failed to read uploaded file for content-type detection"
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInternalServerError, msg))
		return
	}
	fileType := http.DetectContentType(fileBuffer)

	// Reset file reader position
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		msg := "failed to reset file pointer after reading"
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInternalServerError, msg))
		return
	}

	allowedTypes := map[string]bool{
		"image/jpeg":    true,
		"image/png":     true,
		"image/gif":     true,
		"image/webp":    true,
		"image/avif":    true,
		"image/svg+xml": true,
	}

	if !allowedTypes[fileType] {
		cShared.RespondError(w, dShared.NewAppError(dShared.ErrCodeInternalServerError, "invalid image format"))
		return
	}

	err = h.UserSrv.UploadAvatar(context.Background(), header.Filename, file, userID, fileType)
	if err != nil {
		cShared.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
