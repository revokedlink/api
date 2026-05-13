package hooks

import (
	"revoked/util"

	"github.com/pocketbase/pocketbase/core"
)

func BindApiKeyHooks(app core.App) {
	app.OnRecordCreateRequest(util.Coll.ApiKeys).BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Auth == nil {
			return e.RequestEvent.UnauthorizedError(util.Errors.NotAuthenticated.ErrorCode, nil)
		}

		user, err := e.App.FindRecordById(util.Coll.Users, e.Auth.Id)
		if err != nil {
			return e.RequestEvent.InternalServerError("", nil)
		}

		requestedWS := e.Record.GetString(util.Fields.ApiKey.Workspace)
		activeWS := user.GetString(util.Fields.User.ActiveWorkspace)

		if requestedWS != activeWS {
			return e.RequestEvent.BadRequestError(util.Errors.ForbiddenWorkspaceAccess.ErrorCode, nil)
		}

		var raw struct {
			Scopes []string `json:"scopes"`
		}
		if err := e.RequestEvent.BindBody(&raw); err == nil {
			seen := make(map[string]bool)
			for _, s := range raw.Scopes {
				if seen[s] {
					return e.RequestEvent.BadRequestError(util.Errors.DuplicateValues.ErrorCode, nil)
				}
				seen[s] = true
			}
		}

		scopes := e.Record.GetStringSlice(util.Fields.ApiKey.Scopes)
		validScopes := make(map[string]bool)
		for _, s := range util.AllScopes {
			validScopes[s] = true
		}

		for _, s := range scopes {
			if !validScopes[s] {
				return e.RequestEvent.BadRequestError(util.Errors.InvalidActiveWorkspace.ErrorCode, nil)
			}
		}

		return e.Next()
	})
}
