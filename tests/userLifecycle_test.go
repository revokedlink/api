package tests

import (
	"net/http"
	"revoked/tests/testutils"
	"revoked/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserLifecycle_Refactored(t *testing.T) {
	baseURL, _ := testutils.SetupTestApp(t)
	api := testutils.NewPBClient(t, baseURL)

	var token, userID, personalWorkspaceID string

	email := "lifecycle@test.com"
	pass := "password12345"

	t.Run("create user and authenticate", func(t *testing.T) {
		res := api.Create(util.Coll.Users, "", map[string]any{
			"email":           email,
			"password":        pass,
			"passwordConfirm": pass,
		}).Expect().Status(http.StatusOK)

		userID = testutils.ExtractString(res, "id")
		assert.NotEmpty(t, userID)

		authRes := api.AuthWithPassword(util.Coll.Users, email, pass).
			Expect().Status(http.StatusOK)

		token = testutils.ExtractString(authRes, "token")
		assert.NotEmpty(t, token)
	})

	t.Run("get personal workspace", func(t *testing.T) {
		res := api.List(util.Coll.Workspaces, token).
			Expect().Status(http.StatusOK).JSON().Object()

		items := res.Value("items").Array()
		items.Length().IsEqual(1)

		personalWorkspaceID = items.First().Object().Value("id").String().Raw()
		assert.NotEmpty(t, personalWorkspaceID)
	})

	t.Run("reject second personal workspace", func(t *testing.T) {
		resp := api.Create(util.Coll.Workspaces, token, map[string]any{
			"name": "Second Personal",
			"type": util.TypePersonal,
			"slug": "personal-2",
		}).Expect()

		testutils.AssertBadRequestErrors(t, resp, map[string]util.AppError{
			"user": util.Errors.PersonalWorkspaceLimitReached,
		})
	})

	for _, tc := range []struct {
		name string
		slug string
	}{
		{"create business workspace 1", "b-1"},
		{"create business workspace 2", "b-2"},
		{"create business workspace 3", "b-3"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			api.Create(util.Coll.Workspaces, token, map[string]any{
				"name": "Business",
				"type": util.TypeBusiness,
				"slug": tc.slug,
			}).Expect().Status(http.StatusOK)
		})
	}

	t.Run("reject 4th business workspace", func(t *testing.T) {
		resp := api.Create(util.Coll.Workspaces, token, map[string]any{
			"name": "B4",
			"type": util.TypeBusiness,
			"slug": "b-4",
		}).Expect()

		testutils.AssertBadRequestErrors(t, resp, map[string]util.AppError{
			"user": util.Errors.BusinessWorkspaceLimitReached,
		})
	})

	t.Run("delete personal workspace and verify user context", func(t *testing.T) {
		api.Delete(util.Coll.Workspaces, personalWorkspaceID, token).
			Expect().Status(http.StatusNoContent)

		user := api.Get(util.Coll.Users, userID, token).
			Expect().Status(http.StatusOK).JSON().Object()

		user.Value(util.Fields.User.ActiveWorkspace).IsEqual("")
	})

	t.Run("delete user completely", func(t *testing.T) {
		api.Delete(util.Coll.Users, userID, token).
			Expect().Status(http.StatusNoContent)

		// Verify user is truly gone by attempting login
		api.AuthWithPassword(util.Coll.Users, email, pass).
			Expect().Status(http.StatusBadRequest)
	})
}
