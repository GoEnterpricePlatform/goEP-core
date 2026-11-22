package navigation

import sharedNavigation "github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"

// AdminMenu contains the pages owned by the billing module.
// Items with Children link to a section page and its child pages.
var AdminMenu = []sharedNavigation.Item{
	{
		Title:      "Overview",
		Href:       "/v1/billing/",
		ActivePage: "overview",
	},
	{
		Title:      "Example group",
		Href:       "/v1/billing/example-group",
		ActivePage: "example-group",
		Children: []sharedNavigation.Item{
			{Title: "Example child", Href: "/v1/billing/example", ActivePage: "example"},
		},
	},
}
