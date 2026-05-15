package hooks

import (
	"revoked/util"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/security"
)

func BindApiKeyHooks(app core.App) {
	app.OnRecordCreateRequest(util.Coll.ApiKeys).BindFunc(func(e *core.RecordRequestEvent) error {
		if e.Auth == nil {
			return e.RequestEvent.UnauthorizedError(util.Errors.NotAuthenticated.ErrorCode, nil)
		}

		// Only apply user-specific restrictions if the creator is a user record
		if e.Auth.Collection().Name == util.Coll.Users {
			user, err := e.App.FindRecordById(util.Coll.Users, e.Auth.Id)
			if err != nil {
				return e.RequestEvent.InternalServerError("Failed to retrieve user context", nil)
			}

			// Automatically set the user if not provided
			if e.Record.GetString(util.Fields.ApiKey.User) == "" {
				e.Record.Set(util.Fields.ApiKey.User, user.Id)
			}

			// Verify the user is creating the key for their own active workspace
			requestedWS := e.Record.GetString(util.Fields.ApiKey.Workspace)
			activeWS := user.GetString(util.Fields.User.ActiveWorkspace)

			if activeWS == "" {
				return e.RequestEvent.BadRequestError(util.Errors.InvalidActiveWorkspace.ErrorCode, nil)
			}

			if requestedWS == "" {
				e.Record.Set(util.Fields.ApiKey.Workspace, activeWS)
				requestedWS = activeWS
			}

			if requestedWS != activeWS {
				return e.RequestEvent.BadRequestError(util.Errors.ForbiddenWorkspaceAccess.ErrorCode, nil)
			}
		}

		// Generate a secure token if not provided
		if e.Record.GetString(util.Fields.ApiKey.Token) == "" {
			e.Record.Set(util.Fields.ApiKey.Token, security.RandomString(48))
		}

		// Validate Scopes
		var raw struct {
			Scopes []string `json:"scopes"`
		}
		if err := e.RequestEvent.BindBody(&raw); err == nil && len(raw.Scopes) > 0 {
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

func BindApiKeyAuthMiddleware(app core.App) {
	// Handler for single record requests
	recordHandler := func(e *core.RecordRequestEvent) error {
		token := e.Request.Header.Get("X-API-Key")
		if token == "" {
			return e.Next()
		}

		apiKey, err := app.FindFirstRecordByFilter(util.Coll.ApiKeys, "token = {:token}", map[string]any{"token": token})
		if err == nil && apiKey != nil {
			e.Auth = apiKey
		}

		return e.Next()
	}

	app.OnRecordViewRequest().BindFunc(recordHandler)
	app.OnRecordCreateRequest().BindFunc(recordHandler)
	app.OnRecordUpdateRequest().BindFunc(recordHandler)
	app.OnRecordDeleteRequest().BindFunc(recordHandler)

	// Handler for list requests
	app.OnRecordsListRequest().BindFunc(func(e *core.RecordsListRequestEvent) error {
		token := e.Request.Header.Get("X-API-Key")
		if token == "" {
			return e.Next()
		}

		apiKey, err := app.FindFirstRecordByFilter(util.Coll.ApiKeys, "token = {:token}", map[string]any{"token": token})
		if err == nil && apiKey != nil {
			e.Auth = apiKey
		}

		return e.Next()
	})
}
