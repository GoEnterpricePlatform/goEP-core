package service

import planDomain "github.com/GoEnterpricePlatform/goEP-core/addons/catalog/plans/domain"

func isSupportedPlanItemStatus(status planDomain.PlanStatus) bool {
	switch status {
	case planDomain.PlanStatusActive, planDomain.PlanStatusInactive, planDomain.PlanStatusDraft:
		return true
	default:
		return false
	}
}
