package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/amorindev/go-cms-tmpl/pkg/posts/core"
	"github.com/amorindev/go-cms-tmpl/pkg/posts/domain"
	sharedC "github.com/amorindev/go-cms-tmpl/pkg/shared/api/core"
	sharedD "github.com/amorindev/go-cms-tmpl/pkg/shared/domain"
)

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		sharedC.RespondError(w, sharedD.NewAppError(sharedD.ErrCodeInvalidParams, "missing post id"))
		return
	}

	var req core.UpdatePostReq
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
		Title: req.Title,
		Desc: req.Desc,
	}

	if err := h.PostSrv.Update(context.Background(), id, post); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(post)
}