package handler

import (
	"context"
	"net/http"

	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
	"github.com/starfederation/datastar-go/datastar"
)

func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {

	// It is used to test the state of charge
	// time.Sleep(3 * time.Second)

	// It serves to validate that the error is displayed correctly in the UI
	/* if true {
		sse := datastar.NewSSE(w, r)

		_ = sse.MarshalAndPatchSignals(map[string]any{
			"error":      "Invalid request",
			"is_loading": false,
		})
		return
	} */

	id := r.PathValue("id")

	sse := datastar.NewSSE(w, r)

	if id == "" {
		_ = sse.MarshalAndPatchSignals(map[string]any{
			"error":      "Role ID is required",
			"is_loading": false,
		})

		return
	}

	if err := h.RoleSrv.Delete(context.Background(), id); err != nil {
		_ = sse.MarshalAndPatchSignals(map[string]any{
			"error":      sharedC.UiErrorResp(err),
			"is_loading": false,
		})

		return
	}

	_ = sse.MarshalAndPatchSignals(map[string]any{
		"error":      "",
		"is_loading": false,

		"modal": map[string]any{
			"domain": "",
			"action": "",
		},

		"role_form": map[string]any{
			"id":   "",
			"name": "",
		},
	})

	_ = sse.RemoveElement(
		"#role-row-" + id,
	)
}
