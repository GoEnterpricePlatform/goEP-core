package domain

import "time"

// Variation represents a variation for a catalog item, such as
// a one-time product or a subscription plan.
//
// Examples:
//
// One-time product:
//   Name: "Color"
//   Options:
//     - Black
//     - White
//
// Subscription plan:
//   Name: "Billing Interval"
//   Options:
//     - Monthly
//     - Yearly
type Variation struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Options   []*VarOption `json:"options,omitempty"`
	CreatedAt *time.Time   `json:"created_at"`
	UpdatedAt *time.Time   `json:"updated_at"`
}