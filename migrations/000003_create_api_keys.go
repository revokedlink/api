package migrations

import (
	"revoked/util"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	migrations.Register(func(app core.App) error {
		users, err := app.FindCollectionByNameOrId(util.Coll.Users)
		if err != nil {
			return err
		}

		workspaces, err := app.FindCollectionByNameOrId(util.Coll.Workspaces)
		if err != nil {
			return err
		}

		apiKeys := core.NewAuthCollection(util.Coll.ApiKeys)
		apiKeys.Fields.Add(
			&core.TextField{
				Name:     util.Fields.ApiKey.Token,
				Required: true,
				Min:      32,
				Max:      128,
			},
			&core.RelationField{
				Name:          util.Fields.ApiKey.User,
				CollectionId:  users.Id,
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.RelationField{
				Name:          util.Fields.ApiKey.Workspace,
				CollectionId:  workspaces.Id,
				Required:      true,
				MaxSelect:     1,
				CascadeDelete: true,
			},
			&core.SelectField{
				Name:      util.Fields.ApiKey.Scopes,
				Required:  false,
				MaxSelect: len(util.AllScopes),
				Values:    util.AllScopes,
			},
		)

		apiKeys.AddIndex("idxApiKeyToken", true, util.Fields.ApiKey.Token, "")

		// Only the authenticated user can view, create and delete his own keys. Updating is not allowed
		// API access is not permitted for any CRUD
		apiKeys.ListRule = types.Pointer(util.UserSelfOnly())
		apiKeys.ViewRule = types.Pointer(util.UserSelfOnly())
		apiKeys.DeleteRule = types.Pointer(util.UserSelfOnly())
		apiKeys.CreateRule = types.Pointer(util.WorkspaceAdminSelfOnly("", util.Fields.ApiKey.Workspace))
		apiKeys.UpdateRule = nil

		if err := app.Save(apiKeys); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		col, err := app.FindCollectionByNameOrId(util.Coll.ApiKeys)
		if err != nil {
			return nil
		}
		return app.Delete(col)
	})
}
