package navigation

import sharedNavigation "github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"

// AdminMenu contains the dashboard sections owned by Catalog.
var AdminMenu = []sharedNavigation.Item{
	{
		Title:      "General",
		Href:       "/v1/catalog/",
		ActivePage: "general",
		Permission: "view.catalog",
	},
	{
		Title:      "Variations",
		Href:       "/v1/catalog/variations",
		ActivePage: "variations",
		Permission: "view.catalog",
	},
	{
		Title:      "Paddle plans",
		Href:       "/v1/catalog/paddle-plans",
		ActivePage: "paddle-plans",
		Permission: "view.catalog",
	},
	{
		Title:      "Polar plans",
		Href:       "/v1/catalog/polar-plans",
		ActivePage: "polar-plans",
		Permission: "view.catalog",
	},
}
