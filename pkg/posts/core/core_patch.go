package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

// Note: We use pointers for both fields to distinguish between "not provided" and "explicitly empty".
//   - If the pointer is nil → the field was not sent, so we don't update it.
//   - If the pointer is non-nil but empty (e.g. ""), it means the client wants to clear the value.
//
// This allows partial updates without accidentally overwriting fields with zero values.
type PatchPostReq struct {
	Title *string `json:"title"`
	Desc  *string `json:"desc"`
}

func (req PatchPostReq) Validate() error {
	if req.Title != nil {
		name := strings.TrimSpace(*req.Title)
		if name == "" {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "title is required")
		}

		if len(name) > 100 {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "title cannot exceed 100 characters")
		}
	}

	if req.Desc != nil {
		desc := strings.TrimSpace(*req.Desc)
		if len(desc) > 255 {
			return domain.NewAppError(domain.ErrCodeInvalidParams, "desc cannot exceed 255 characters")
		}
	}
	return nil
}
