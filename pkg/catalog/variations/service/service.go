package service

import "github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/variations/port"

var _ port.VariationSrv = &Service{}

type Service struct {
	VariationRepo port.VariationRepo
	VarOptionRepo port.VarOptionRepo
	OptionUsage   port.OptionUsageChecker
}

func NewVariationSrv(variationRepo port.VariationRepo, varOptionRepo port.VarOptionRepo, optionUsage port.OptionUsageChecker) *Service {
	return &Service{
		VariationRepo: variationRepo,
		VarOptionRepo: varOptionRepo,
		OptionUsage:   optionUsage,
	}
}
