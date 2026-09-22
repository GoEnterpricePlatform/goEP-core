package helpers

import (
	"sort"

	"github.com/GoEnterpricePlatform/goEP-core/pkg/catalog/plans/domain"
)

// CalculateVariations builds and assigns the list of unique plan variations
// (e.g., Billing Cycle) based on the options defined across all PlanItems.
//
// Process:
//  1. Iterate over each PlanItem and its options.
//  2. Group option values by variation name using a map (acting as a set to avoid duplicates).
//  3. Convert each set into a slice and sort values alphabetically for consistency.
//  4. Assign the resulting list to p.Variations.
//
// Example:
//  - PlanItems with options:
//      Plan=Basic, Billing Cycle=Monthly
//      Plan=Basic, Billing Cycle=Yearly
//      Plan=Pro, Billing Cycle=Monthly
//      Plan=Pro, Billing Cycle=Yearly
//      Plan=Custom
//  - Result in p.Variations:
//      [
//        {Name: "Plan", Values: ["Monthly", "Yearly"]}
//      ]
func CalculateVariations(p *domain.Plan) {
	variationMap := make(map[string]map[string]struct{})

	for _, item := range p.Items {
		for _, opt := range item.Options {
			if _, exists := variationMap[opt.Name]; !exists {
				variationMap[opt.Name] = make(map[string]struct{})
			}
			variationMap[opt.Name][opt.VarOptName] = struct{}{}
		}
	}

	var variations []*domain.Variation
	for name, valuesSet := range variationMap {
		var values []string
		for val := range valuesSet {
			values = append(values, val)
		}
		sort.Strings(values)
		variations = append(variations, &domain.Variation{
			Name:   name,
			Values: values,
		})
	}

	p.Variations = variations
}