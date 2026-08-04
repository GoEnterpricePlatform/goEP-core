package domain

import "time"

// VarOption represents a selectable value for a variation.
//
// Product examples:
//   Variation: "Color"
//     - Label: "Black"
//       Value: "#000000"
//
//   Variation: "Size"
//     - Label: "XL"
//       Value: nil
//
// Subscription plan examples:
//   Variation: "Billing Interval"
//     - Label: "Monthly"
//       Value: nil
type VarOption struct {
	ID          string     `json:"id"`
	VariationID string     `json:"variation_id"`
	Label       string     `json:"label"`
	Value       *string    `json:"value"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
