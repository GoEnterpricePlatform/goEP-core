package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/core"
	"github.com/GoEnterpricePlatform/goEP-core/pkg/posts/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/api/core"
	sharedD "github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

func (h Handler) Patch(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing post id"))
		return
	}

	var req core.PatchPostReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "invalid request body"))
		return
	}
	defer r.Body.Close()

	if err := req.Validate(); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	post := &domain.Post{
		Desc: req.Desc,
	}

	if req.Title != nil {
		post.Title = *req.Title
	}

	updated, err := h.PostSrv.Patch(context.Background(), id, post)
	if err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updated)
}
