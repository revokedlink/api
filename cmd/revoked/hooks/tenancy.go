package hooks

import (
	"revoked/util"

	"github.com/pocketbase/pocketbase/core"
)

// RegisterTenancyHooks registers hooks that automatically fill the workspace and user fields
// for business collections based on the authenticated context.
func RegisterTenancyHooks(app core.App) {
	collections := []string{
		util.Coll.Documents,
	}

	for _, collName := range collections {
		app.OnRecordCreateRequest(collName).BindFunc(func(e *core.RecordRequestEvent) error {
			if e.Auth == nil {
				return e.Next()
			}

			e.Record.Set(util.Fields.Document.User, e.Auth.Id)

			if e.Auth.Collection().Name == util.Coll.Users {
				e.Record.Set(util.Fields.Document.Workspace, e.Auth.GetString(util.Fields.User.ActiveWorkspace))
			} else if e.Auth.Collection().Name == util.Coll.ApiKeys {
				e.Record.Set(util.Fields.Document.Workspace, e.Auth.GetString(util.Fields.ApiKey.Workspace))
			}

			return e.Next()
		})
	}
}
