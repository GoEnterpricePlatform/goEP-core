package service

import (
	"context"
	"time"

	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
)

func (h *Service) Update(ctx context.Context, id string, post *domain.Post) error {
	now := time.Now().UTC()
	post.UpdatedAt = &now

	err := h.PostRepo.Update(ctx, id, post)
	if err != nil {
		return sharedD.ManageError(err, "")
	}
	return nil
}
