package hooks

import (
	"revoked/cmd/revoked/routes"

	"github.com/pocketbase/pocketbase/core"
)

func BindHooksAndRoutes(app core.App) {
	BindUsersHooks(app)
	BindWorkspacesHooks(app)
	BindWorkspaceMembersHooks(app)
	BindApiKeyHooks(app)
	RegisterTenancyHooks(app)
	BindAuditLogHooks(app)
	routes.HealthzRoute(app)
}
