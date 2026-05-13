package hooks

import (
	"github.com/pocketbase/pocketbase/core"
)

// BindCreateSuperuserAccount automatically creates a superuser if one doesn't exist
func BindCreateSuperuserAccount(app core.App, email string, password string) {
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {

		_, err := se.App.FindFirstRecordByFilter(
			core.CollectionNameSuperusers,
			"email = {:email}",
			map[string]any{"email": email},
		)

		if err != nil {
			superusersCollection, err := se.App.FindCollectionByNameOrId(core.CollectionNameSuperusers)
			if err == nil {
				record := core.NewRecord(superusersCollection)
				record.Set("email", email)
				record.Set("password", password)
				record.Set("passwordConfirm", password)

				if err := se.App.Save(record); err == nil {
					se.App.Logger().Info("🌱 Superuser seeded successfully!")
				} else {
					se.App.Logger().Error("Failed to seed superuser", "error", err)
				}
			} else {
				se.App.Logger().Error("Could not find superusers collection", "error", err)
			}
		}

		return se.Next()
	})
}
