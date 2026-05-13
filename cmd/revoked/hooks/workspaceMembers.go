package hooks

import (
	"fmt"
	"revoked/util"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

func BindWorkspaceMembersHooks(app core.App) {
	// Prevent duplicate memberships and enforce member limits per workspace type
	app.OnRecordCreate(util.Coll.WorkspaceMembers).BindFunc(func(e *core.RecordEvent) error {
		userId := e.Record.GetString(util.Fields.WorkspaceMember.User)
		workspaceId := e.Record.GetString(util.Fields.WorkspaceMember.Workspace)

		if userId == "" || workspaceId == "" {
			return e.Next()
		}

		existing, err := e.App.FindFirstRecordByFilter(
			util.Coll.WorkspaceMembers,
			fmt.Sprintf(
				"%s = {:user} && %s = {:workspace}",
				util.Fields.WorkspaceMember.User,
				util.Fields.WorkspaceMember.Workspace,
			),
			map[string]any{
				"user":      userId,
				"workspace": workspaceId,
			},
		)
		if err == nil && existing != nil {
			return validation.Errors{
				util.Fields.WorkspaceMember.User: util.AsValidationError(util.Errors.DuplicateWorkspaceMember),
			}
		}

		workspace, err := e.App.FindRecordById(util.Coll.Workspaces, workspaceId)
		if err != nil {
			return err
		}

		workspaceType := workspace.GetString(util.Fields.Workspace.Type)

		var maxMembers int
		switch workspaceType {
		case util.TypePersonal:
			maxMembers = util.MaximumPersonalMembers
		case util.TypeBusiness:
			maxMembers = util.MaximumBusinessMembers
		default:
			return e.Next()
		}

		count, err := e.App.CountRecords(
			util.Coll.WorkspaceMembers,
			dbx.HashExp{util.Fields.WorkspaceMember.Workspace: workspaceId},
		)
		if err != nil {
			return err
		}

		if int(count) >= maxMembers {
			return validation.Errors{
				util.Fields.WorkspaceMember.Workspace: util.AsValidationError(util.Errors.WorkspaceMemberLimitReached),
			}
		}

		return e.Next()
	})
}
