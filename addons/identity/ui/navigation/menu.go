package navigation

import sharedNavigation "github.com/GoEnterpricePlatform/goEP-core/internal/modules/navigation"

// AdminMenu contains the sections owned by Identity's admin UI.
var AdminMenu = []sharedNavigation.Item{
	{
		Title:      "General",
		Href:       "/v1/goep-admin/general",
		ActivePage: "general",
		Permission: "VIEW_GENERAL",
	},
	{
		Title:      "Users",
		Href:       "/v1/goep-admin/users",
		ActivePage: "users",
		Permission: "VIEW_USERS",
	},
	{
		Title:      "Roles and permissions",
		Href:       "/v1/goep-admin/roles-permissions",
		ActivePage: "roles-permissions",
		Permission: "VIEW_ROLES_PERMISSIONS",
	},
	{
		Title:      "Settings",
		Href:       "/v1/goep-admin/settings",
		ActivePage: "settings",
		Permission: "VIEW_SETTINGS",
	},
}
