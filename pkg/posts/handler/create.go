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

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req core.CreatePostReq

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

	if err := h.PostSrv.Create(context.Background(), post); err != nil {
		sharedC.RespondError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}