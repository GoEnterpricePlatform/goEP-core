package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

type UpdateVariationReq struct {
	Name string `json:"name"`
}

func (v UpdateVariationReq) Validate() error {
	if strings.TrimSpace(v.Name) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "name is required")
	}
	return nil
}

type UpdateVarOptionReq struct {
	Label string  `json:"label"`
	Value *string `json:"value"`
}

func (vO UpdateVarOptionReq) Validate() error{
    return validateVarOptionFields(vO.Label,vO.Value)
}