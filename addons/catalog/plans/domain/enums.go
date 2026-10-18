package domain

// PlanStatus defines the possible states of a plan
type PlanStatus string

const (
	PlanStatusActive   PlanStatus = "active"   // The plan is available for sale
	PlanStatusInactive PlanStatus = "inactive" // The plan is not available for customers
	PlanStatusDraft    PlanStatus = "draft"    // The plan is not payment configurated o sync
)
