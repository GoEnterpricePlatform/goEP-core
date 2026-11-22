package navigation

import sharedNavigation "github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"

// AdminMenu contains the pages owned by the billing module.
// Add section entries with Children and register a route and page for each child.
var AdminMenu = []sharedNavigation.Item{
	{
		Title:      "Overview",
		Href:       "/v1/billing/",
		ActivePage: "overview",
	},
	{
		Title: "Example group",
		Children: []sharedNavigation.Item{
			{Title: "Example child", Href: "/v1/billing/example", ActivePage: "example"},
		},
	},
}
