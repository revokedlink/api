package hooks

import (
	"fmt"
	"revoked/util"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func BindWorkspacesHooks(app core.App) {
	app.OnRecordCreateRequest(util.Coll.Workspaces).BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Auth == nil {
			return apis.NewForbiddenError(util.Errors.NotAuthorized.ErrorText, nil)
		}

		workspaceType := e.Record.GetString(util.Fields.Workspace.Type)

		// Check limits
		memberships, err := e.App.FindRecordsByFilter(
			util.Coll.WorkspaceMembers,
			fmt.Sprintf("%s = {:user}", util.Fields.WorkspaceMember.User),
			"",
			0,
			0,
			map[string]any{"user": e.Auth.Id},
		)
		if err != nil {
			return err
		}

		count := 0
		for _, m := range memberships {
			ws, _ := e.App.FindRecordById(util.Coll.Workspaces, m.GetString(util.Fields.WorkspaceMember.Workspace))
			if ws != nil && ws.GetString(util.Fields.Workspace.Type) == workspaceType {
				count++
			}
		}

		if workspaceType == util.TypePersonal && count >= 1 {
			return apis.NewBadRequestError(util.Errors.PersonalWorkspaceLimitReached.ErrorText, nil)
		}

		if workspaceType == util.TypeBusiness && count >= 3 {
			return apis.NewBadRequestError(util.Errors.BusinessWorkspaceLimitReached.ErrorText, nil)
		}

		// Save the workspace first
		nextErr := e.Next()
		if nextErr != nil {
			return nextErr
		}

		// Add creator as admin member
		workspaceMemberships, err := e.App.FindCollectionByNameOrId(util.Coll.WorkspaceMembers)
		if err != nil {
			return nil
		}

		membership := core.NewRecord(workspaceMemberships)
		membership.Set(util.Fields.WorkspaceMember.User, e.Auth.Id)
		membership.Set(util.Fields.WorkspaceMember.Workspace, e.Record.Id)
		membership.Set(util.Fields.WorkspaceMember.Role, util.RoleAdmin)

		if err := e.App.Save(membership); err != nil {
			return err
		}

		return nil
	})

	app.OnRecordAfterDeleteSuccess(util.Coll.Workspaces).BindFunc(func(e *core.RecordEvent) error {
		users, err := e.App.FindRecordsByFilter(
			util.Coll.Users,
			fmt.Sprintf("%s = {:workspace}", util.Fields.User.ActiveWorkspace),
			"",
			0,
			0,
			map[string]any{"workspace": e.Record.Id},
		)
		if err != nil {
			return e.Next()
		}

		for _, user := range users {
			user.Set(util.Fields.User.ActiveWorkspace, "")
			user.Set(util.Fields.User.ActiveRole, "")
			_ = e.App.Save(user)
		}

		return e.Next()
	})
}
