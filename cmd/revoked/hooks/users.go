package hooks

import (
	"fmt"
	"revoked/util"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func BindUsersHooks(app core.App) {
	// Restrict users from updating sensitive fields of their own account
	// and prevent them from updating other users accounts
	app.OnRecordUpdateRequest(util.Coll.Users).BindFunc(func(e *core.RecordRequestEvent) error {
		// If e.Auth is nil, this is a system-level update (e.g. from another hook),
		// so we bypass the restrictions.
		if e.Auth == nil {
			return e.Next()
		}

		// Check if activeWorkspace or activeRole are being updated
		info, _ := e.RequestInfo()
		requestedWS := info.Body[util.Fields.User.ActiveWorkspace]
		requestedRole := info.Body[util.Fields.User.ActiveRole]

		if requestedWS != nil || requestedRole != nil {
			// Determine the target state
			targetWS := e.Record.GetString(util.Fields.User.ActiveWorkspace)
			if requestedWS != nil {
				targetWS = fmt.Sprintf("%v", requestedWS)
			}

			targetRole := e.Record.GetString(util.Fields.User.ActiveRole)
			if requestedRole != nil {
				targetRole = fmt.Sprintf("%v", requestedRole)
			}

			// Allow clearing the context
			if targetWS == "" && targetRole == "" {
				return e.Next()
			}

			// Validate that the user has a membership matching the TARGET state
			filter := fmt.Sprintf("%s = {:workspace} && %s = {:user} && %s = {:role}",
				util.Fields.WorkspaceMember.Workspace,
				util.Fields.WorkspaceMember.User,
				util.Fields.WorkspaceMember.Role,
			)
			params := map[string]any{
				"workspace": targetWS,
				"user":      e.Record.Id,
				"role":      targetRole,
			}

			member, err := e.App.FindFirstRecordByFilter(util.Coll.WorkspaceMembers, filter, params)

			if err != nil || member == nil {
				return apis.NewNotFoundError("", nil)
			}
		}

		if err := util.RestrictFields(e,
			util.Fields.User.Email,
			util.Fields.User.Verified,
			util.Fields.User.Avatar,
			util.Fields.User.Active,
		); err != nil {
			return err
		}

		return e.Next()
	})

	// Automatically create a personal workspace when a users account is created.
	// We NO LONGER set it as active automatically to avoid race conditions during onboarding.
	app.OnRecordAfterCreateSuccess(util.Coll.Users).BindFunc(func(e *core.RecordEvent) error {
		workspaces, err := e.App.FindCollectionByNameOrId(util.Coll.Workspaces)
		if err != nil {
			return err
		}

		workspace := core.NewRecord(workspaces)
		workspace.Set(util.Fields.Workspace.Name, "Personal")
		workspace.Set(util.Fields.Workspace.Type, util.TypePersonal)
		workspace.Set(util.Fields.Workspace.Slug, "personal-"+e.Record.Id)
		if err := e.App.Save(workspace); err != nil {
			return err
		}

		workspaceMemberships, err := e.App.FindCollectionByNameOrId(util.Coll.WorkspaceMembers)
		if err != nil {
			return err
		}

		membership := core.NewRecord(workspaceMemberships)
		membership.Set(util.Fields.WorkspaceMember.User, e.Record.Id)
		membership.Set(util.Fields.WorkspaceMember.Workspace, workspace.Id)
		membership.Set(util.Fields.WorkspaceMember.Role, util.RoleAdmin)

		if err := e.App.Save(membership); err != nil {
			return err
		}

		return e.Next()
	})

	// Automatically delete the user's personal workspace after the user's account is deleted
	app.OnRecordAfterDeleteSuccess(util.Coll.Users).BindFunc(func(e *core.RecordEvent) error {
		workspace, err := e.App.FindFirstRecordByFilter(
			util.Coll.Workspaces,
			fmt.Sprintf("%s = {:slug}", util.Fields.Workspace.Slug),
			map[string]any{
				"slug": "personal-" + e.Record.Id,
			},
		)

		if err != nil {
			return e.Next()
		}

		if workspace.GetString(util.Fields.Workspace.Type) == util.TypePersonal {
			if err := e.App.Delete(workspace); err != nil {
				return err
			}
		}

		return e.Next()
	})
}
