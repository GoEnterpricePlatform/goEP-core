package handler

import (
	"context"
	"net/http"

	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/core"
	"github.com/GoEnterpricePlatform/goEP-core/addons/identity/roles/domain"
	sharedC "github.com/GoEnterpricePlatform/goEP-core/web/shared/core"
	"github.com/GoEnterpricePlatform/goEP-core/web/web-app/modules/identity/ui/components"
	"github.com/starfederation/datastar-go/datastar"
)

func (h Handler) Update(w http.ResponseWriter, r *http.Request) {

	// It is used to test the state of charge
	//time.Sleep(3 * time.Second)

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

	var signals struct {
		RoleForm struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"role_form"`
	}

	if err := datastar.ReadSignals(r, &signals); err != nil {
		sse := datastar.NewSSE(w, r)

		_ = sse.MarshalAndPatchSignals(map[string]any{
			"error":      "Invalid request",
			"is_loading": false,
		})

		return
	}

	if id == "" {
		sse := datastar.NewSSE(w, r)

		_ = sse.MarshalAndPatchSignals(map[string]any{
			"error":      "Role ID is required",
			"is_loading": false,
		})

		return
	}

	req := &core.UpdateRoleReq{
		Name: signals.RoleForm.Name,
	}

	if err := req.Validate(); err != nil {
		sse := datastar.NewSSE(w, r)

		_ = sse.MarshalAndPatchSignals(map[string]any{
			"error":      sharedC.UiErrorResp(err),
			"is_loading": false,
		})

		return
	}

	role := &domain.Role{
		ID:   id,
		Name: req.Name,
	}

	if err := h.RoleSrv.Update(context.Background(), id, role); err != nil {
		sse := datastar.NewSSE(w, r)

		_ = sse.MarshalAndPatchSignals(map[string]any{
			"error":      sharedC.UiErrorResp(err),
			"is_loading": false,
		})

		return
	}

	sse := datastar.NewSSE(w, r)

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

	_ = sse.PatchElementTempl(
		components.RoleRow(role),
		datastar.WithSelector("#role-row-"+id),
		datastar.WithModeOuter(),
	)
}
