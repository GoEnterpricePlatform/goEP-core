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
		Permission: "view.general",
	},
	{
		Title:      "Users",
		Href:       "/v1/goep-admin/users",
		ActivePage: "users",
		Permission: "view.users",
	},
	{
		Title:      "Roles and permissions",
		Href:       "/v1/goep-admin/roles-permissions",
		ActivePage: "roles-permissions",
		Permission: "view.roles.permissions",
	},
	{
		Title:      "Settings",
		Href:       "/v1/goep-admin/settings",
		ActivePage: "settings",
		Permission: "view.settings",
	},
}