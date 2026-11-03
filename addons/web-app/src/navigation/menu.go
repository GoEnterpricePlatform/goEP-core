package navigation

type MenuItem struct {
	Title      string
	Href       string
	ActivePage string
	Permission string
}

var DashboardMenu = []MenuItem{
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