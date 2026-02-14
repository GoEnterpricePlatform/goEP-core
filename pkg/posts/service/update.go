package service

import (
	"context"
	"time"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
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