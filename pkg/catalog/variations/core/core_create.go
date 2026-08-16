package core

import (
	"strings"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/shared/domain"
)

type CreateVariationReq struct {
	Name string `json:"name"`
}

func (v CreateVariationReq) Validate() error {
	if strings.TrimSpace(v.Name) == "" {
		return domain.NewAppError(domain.ErrCodeInvalidParams, "name is required")
	}
	return nil
}

type CreateVarOptionReq struct {
	Label string  `json:"label"`
	Value *string `json:"value"`
}

func (vO CreateVarOptionReq) Validate() error {
	return validateVarOptionFields(vO.Label, vO.Value)
}