package util

const (
	SlugPattern                = "^[a-z0-9_-]+$"
	MaximumBusinessMemberships = 3
	MaximumPersonalMemberships = 1
	MaximumBusinessMembers     = 50
	MaximumPersonalMembers     = 1
	RoleAdmin                  = "admin"
	RoleMember                 = "member"

	TypePersonal = "personal"
	TypeBusiness = "business"

	ADMIN    = "dev@dev.com"
	PASSWORD = "dev@dev.com"
)

var WorkspaceRoles = []string{RoleAdmin, RoleMember}
var WorkspaceTypes = []string{TypePersonal, TypeBusiness}
