package migrations

import (
	"revoked/util"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	migrations.Register(func(app core.App) error {

		workspaces, err := app.FindCollectionByNameOrId(util.Coll.Workspaces)
		if err != nil {
			return err
		}
		workspaces.ListRule = types.Pointer(util.WorkspaceAnyAdmin(util.ScopeWorkspacesRead, "id"))
		workspaces.ViewRule = types.Pointer(util.WorkspaceAnyAdmin(util.ScopeWorkspacesRead, "id"))
		workspaces.CreateRule = types.Pointer("@request.auth.id != ''")
		workspaces.UpdateRule = types.Pointer(util.WorkspaceAnyAdmin(util.ScopeWorkspacesUpdate, "id") + " && @request.body.type:isset = false")
		workspaces.DeleteRule = types.Pointer(util.WorkspaceAnyAdmin(util.ScopeWorkspacesDelete, "id"))

		if err := app.Save(workspaces); err != nil {
			return err
		}

		workspaceMembers, err := app.FindCollectionByNameOrId(util.Coll.WorkspaceMembers)
		if err != nil {
			return err
		}

		// Users can see their own membership OR admins can see all memberships
		listRule := "user = @request.auth.id || " + util.WorkspaceAnyAdmin(util.ScopeWorkspaceMembersRead, "workspace")
		workspaceMembers.ListRule = &listRule
		workspaceMembers.ViewRule = &listRule

		workspaceMembers.CreateRule = types.Pointer(util.WorkspaceAnyAdmin(util.ScopeWorkspaceMembersCreate, "@request.body.workspace"))
		workspaceMembers.UpdateRule = types.Pointer(util.WorkspaceAnyAdmin(util.ScopeWorkspaceMembersUpdate, "workspace"))
		workspaceMembers.DeleteRule = types.Pointer(util.WorkspaceAnyAdmin(util.ScopeWorkspaceMembersDelete, "workspace"))

		if err := app.Save(workspaceMembers); err != nil {
			return err
		}

		users, err := app.FindCollectionByNameOrId(util.Coll.Users)
		if err != nil {
			return err
		}
		users.CreateRule = types.Pointer("@request.auth.id = ''")
		users.ViewRule = types.Pointer("id = @request.auth.id")
		users.UpdateRule = types.Pointer("@request.auth.id != ''")

		if err := app.Save(users); err != nil {
			return err
		}

		return nil

	}, func(app core.App) error {
		return nil
	})
}
