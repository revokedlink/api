package tests

import (
	"net/http"
	"revoked/tests/testutils"
	"revoked/util"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkspaceAdvanced_SideEffects_Refactored(t *testing.T) {
	baseURL, _ := testutils.SetupTestApp(t)
	api := testutils.NewPBClient(t, baseURL)

	// User A: The Workspace Owner
	userA_ID, userA_Token, _ := testutils.CreateRandomUser(baseURL)
	// User B: The Regular Member
	userB_ID, userB_Token, _ := testutils.CreateRandomUser(baseURL)

	var workspaceID string

	t.Run("Setup: User A creates a business workspace and User B is added", func(t *testing.T) {
		api := api.T(t)
		// 1. User A creates a business workspace
		ws := api.Create(util.Coll.Workspaces, userA_Token, map[string]any{
			"name": "Shared Business",
			"type": util.TypeBusiness,
		}).Expect().Status(http.StatusOK)

		workspaceID = testutils.ExtractString(ws, "id")

		// 2. User A MUST activate the workspace to manage members
		api.Update(util.Coll.Users, userA_ID, userA_Token, map[string]any{
			"activeWorkspace": workspaceID,
			"activeRole":      util.RoleAdmin,
		}).Expect().Status(http.StatusOK)

		// 3. User A adds User B as a 'member' (not admin)
		api.Create(util.Coll.WorkspaceMembers, userA_Token, map[string]any{
			"user":      userB_ID,
			"workspace": workspaceID,
			"role":      util.RoleMember,
		}).Expect().Status(http.StatusOK)
	})

	t.Run("Security: Admin cannot add users to a workspace that is not their active context", func(t *testing.T) {
		api := api.T(t)
		// 1. User A creates a SECOND business workspace
		ws := api.Create(util.Coll.Workspaces, userA_Token, map[string]any{
			"name": "User A Private Workspace",
			"type": util.TypeBusiness,
		}).Expect().Status(http.StatusOK)

		privateWorkspaceID := testutils.ExtractString(ws, "id")

		// 2. User A switches their active context BACK to the first "Shared Business" workspace
		api.Update(util.Coll.Users, userA_ID, userA_Token, map[string]any{
			"activeWorkspace": workspaceID,
			"activeRole":      util.RoleAdmin,
		}).Expect().Status(http.StatusOK)

		// 3. User A tries to add User B to the SECOND (private) workspace (Expect 403)
		api.Create(util.Coll.WorkspaceMembers, userA_Token, map[string]any{
			"user":      userB_ID,
			"workspace": privateWorkspaceID,
			"role":      util.RoleMember,
		}).Expect().Status(http.StatusForbidden)
	})

	t.Run("Security: Regular member cannot delete or rename workspace", func(t *testing.T) {
		api := api.T(t)
		// 1. User B switches to the shared workspace as a 'member'
		api.Update(util.Coll.Users, userB_ID, userB_Token, map[string]any{
			"activeWorkspace": workspaceID,
			"activeRole":      util.RoleMember,
		}).Expect().Status(http.StatusOK)

		// 2. User B tries to rename the workspace (Expect 404 because rules hide records you can't touch)
		api.Update(util.Coll.Workspaces, workspaceID, userB_Token, map[string]any{
			"name": "Hacked Name",
		}).Expect().Status(http.StatusNotFound)

		// 3. User B tries to delete the workspace (Expect 404)
		api.Delete(util.Coll.Workspaces, workspaceID, userB_Token).
			Expect().Status(http.StatusNotFound)
	})

	t.Run("Security: User cannot elevate their own role to admin in a workspace they only have member access to", func(t *testing.T) {
		api := api.T(t)
		// User B tries to set themselves as 'admin' for the shared workspace
		api.Update(util.Coll.Users, userB_ID, userB_Token, map[string]any{
			"activeWorkspace": workspaceID,
			"activeRole":      util.RoleAdmin,
		}).Expect().Status(http.StatusForbidden)
	})

	t.Run("Side Effect: Workspace deletion clears context for all members", func(t *testing.T) {
		api := api.T(t)
		// 1. Ensure User B is active in the shared workspace as 'member'
		api.Update(util.Coll.Users, userB_ID, userB_Token, map[string]any{
			"activeWorkspace": workspaceID,
			"activeRole":      util.RoleMember,
		}).Expect().Status(http.StatusOK)

		// 2. User A (the real admin) deletes the workspace
		api.Delete(util.Coll.Workspaces, workspaceID, userA_Token).
			Expect().Status(http.StatusNoContent)

		// 3. Verify User B's context was cleared automatically
		res := api.Get(util.Coll.Users, userB_ID, userB_Token).
			Expect().Status(http.StatusOK)

		assert.Empty(t, testutils.ExtractString(res, "activeWorkspace"))
		assert.Empty(t, testutils.ExtractString(res, "activeRole"))
	})
}
