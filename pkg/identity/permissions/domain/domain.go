package domain

type Permission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewPermission(name PermissionName) *Permission {
	return &Permission{
		Name: string(name),
	}
}

// system permissions should be saved in the db using initializer
type PermissionName string

const (
	// The user is authorized to enter the administrative domain of the system.
	// What it is for: Separate normal users from administrative users, prevent any authenticated account from operating internal functions
	// When used: access to the administrative panel, starting any internal system operations
	PAdminAccess PermissionName = "admin.access"

	// The user can manage other administrative accounts.
	// What it's for: creating admins, editing other admins' permissions, revoking access
	// When to use it: adding/removing administrators, changing roles or permissions
	PAdminManage PermissionName = "admin.manage"

	// The user can view system users and their basic information.
	// What it is for: Allow support, audit, or review of user accounts without granting modification rights.
	// When used: Listing users, viewing user profiles, support or auditing operations.
	PUserRead PermissionName = "user.read"

	// The user can change the active status of a user account.
	// What it is for: Enable administrative control over user availability without modifying user data.
	// When used: Activating or deactivating user accounts, suspending access, enforcing moderation decisions.
	PUserChangeActive PermissionName = "user.change.active"

	// The user can view global system configuration.
	// What it is for: Provide visibility into system-wide settings without allowing changes.
	// When used: Viewing system configuration, diagnostics, audits, and support-related checks.
	PSettingsRead PermissionName = "settings.read"

	// The user can modify global system configuration.
	// What it is for: Allow controlled changes to system-wide behavior and critical settings.
	// When used: Updating system settings, enabling or disabling features, changing global operational parameters.
	PSettingsUpdate PermissionName = "settings.update"
)
