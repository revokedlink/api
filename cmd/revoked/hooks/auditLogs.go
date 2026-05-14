package hooks

import (
	"encoding/json"
	"revoked/util"

	"github.com/pocketbase/pocketbase/core"
)

func BindAuditLogHooks(app core.App) {
	app.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Collection.Name == util.Coll.AuditLogs {
			return e.Next()
		}

		if err := e.Next(); err != nil {
			return err
		}

		return logAuditAction(app, e, "create", nil, e.Record.PublicExport())
	})

	app.OnRecordUpdateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Collection.Name == util.Coll.AuditLogs {
			return e.Next()
		}

		oldData := e.Record.PublicExport()

		if err := e.Next(); err != nil {
			return err
		}

		return logAuditAction(app, e, "update", oldData, e.Record.PublicExport())
	})

	app.OnRecordDeleteRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Collection.Name == util.Coll.AuditLogs {
			return e.Next()
		}

		// Capture data before deletion
		oldData := e.Record.PublicExport()

		if err := e.Next(); err != nil {
			return err
		}

		return logAuditAction(app, e, "delete", oldData, nil)
	})
}

func logAuditAction(app core.App, e *core.RecordRequestEvent, action string, oldData any, newData any) error {
	auditCollection, err := app.FindCollectionByNameOrId(util.Coll.AuditLogs)
	if err != nil {
		return e.Next() // Don't block the main operation if audit logging fails
	}

	auditRecord := core.NewRecord(auditCollection)

	// Set User
	if e.Auth != nil {
		auditRecord.Set(util.Fields.AuditLog.User, e.Auth.Id)
	}

	// Set Action and Metadata
	auditRecord.Set(util.Fields.AuditLog.Action, action)
	auditRecord.Set(util.Fields.AuditLog.Collection, e.Collection.Name)
	auditRecord.Set(util.Fields.AuditLog.RecordId, e.Record.Id)

	if oldData != nil {
		oldDataJSON, _ := json.Marshal(oldData)
		auditRecord.Set(util.Fields.AuditLog.OldData, string(oldDataJSON))
	}

	if newData != nil {
		newDataJSON, _ := json.Marshal(newData)
		auditRecord.Set(util.Fields.AuditLog.NewData, string(newDataJSON))
	}

	// Set Network Info
	auditRecord.Set(util.Fields.AuditLog.Ip, e.Request.RemoteAddr)
	auditRecord.Set(util.Fields.AuditLog.UserAgent, e.Request.UserAgent())

	// Set Workspace context
	workspaceId := ""

	// 1. Check if the record has a workspace field (most business collections do)
	if ws := e.Record.GetString("workspace"); ws != "" {
		workspaceId = ws
	}

	// 2. If the record IS a workspace itself, log that workspace ID
	if workspaceId == "" && e.Collection.Name == util.Coll.Workspaces {
		workspaceId = e.Record.Id
	}

	// 3. If the record is a Workspace Member, use its workspace field
	if workspaceId == "" && e.Collection.Name == util.Coll.WorkspaceMembers {
		workspaceId = e.Record.GetString(util.Fields.WorkspaceMember.Workspace)
	}

	// 4. Fallback: use the authenticated user's active workspace context
	if workspaceId == "" && e.Auth != nil {
		workspaceId = e.Auth.GetString(util.Fields.User.ActiveWorkspace)
	}

	if workspaceId != "" {
		auditRecord.Set(util.Fields.AuditLog.Workspace, workspaceId)
	}

	// Save audit log in a separate goroutine or just save it.
	// Saving it synchronously ensures it's logged, but might slow down the request slightly.
	if err := app.Save(auditRecord); err != nil {
		app.Logger().Error("Failed to save audit log", "error", err)
	}

	return e.Next()
}
