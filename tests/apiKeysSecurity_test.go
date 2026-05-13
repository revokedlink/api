package tests

import (
	"fmt"
	"net/http"
	"revoked/tests/testutils"
	"revoked/util"
	"testing"
	"time"
)

func TestApiKeySecurity_Rigorous_Refactored(t *testing.T) {
	baseURL, _ := testutils.SetupTestApp(t)
	api := testutils.NewPBClient(t, baseURL)

	userA_ID, userA_Token, _ := testutils.CreateRandomUser(baseURL)
	userB_ID, userB_Token, _ := testutils.CreateRandomUser(baseURL)

	var workspaceID string

	t.Run("Setup: Admin creates a business workspace", func(t *testing.T) {
		api := api.T(t)
		slug := fmt.Sprintf("admin-workspace-%d", time.Now().UnixNano())

		// Create Workspace
		ws := api.Create(util.Coll.Workspaces, userA_Token, map[string]any{
			"name": "Admin Workspace",
			"slug": slug,
			"type": util.TypeBusiness,
		}).Expect().Status(http.StatusOK)

		workspaceID = testutils.ExtractString(ws, "id")

		// Setup User A (Admin)
		api.Update(util.Coll.Users, userA_ID, userA_Token, map[string]any{
			"activeWorkspace": workspaceID,
			"activeRole":      util.RoleAdmin,
		}).Expect().Status(http.StatusOK)

		refreshA := api.AuthRefresh(util.Coll.Users, userA_Token).Expect().Status(http.StatusOK)
		userA_Token = testutils.ExtractString(refreshA, "token")

		// Setup User B (Member)
		api.Create(util.Coll.WorkspaceMembers, userA_Token, map[string]any{
			"user":      userB_ID,
			"workspace": workspaceID,
			"role":      util.RoleMember,
		}).Expect().Status(http.StatusOK)

		api.Update(util.Coll.Users, userB_ID, userB_Token, map[string]any{
			"activeWorkspace": workspaceID,
			"activeRole":      util.RoleMember,
		}).Expect().Status(http.StatusOK)

		refreshB := api.AuthRefresh(util.Coll.Users, userB_Token).Expect().Status(http.StatusOK)
		userB_Token = testutils.ExtractString(refreshB, "token")
	})

	t.Run("Security: Admin can create an API Key for themselves", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()
		api.AssertStatus(api.Create(util.Coll.ApiKeys, userA_Token, map[string]any{
			"username":        fmt.Sprintf("admin-key-%d", ts),
			"email":           fmt.Sprintf("admin-key-%d@test.com", ts),
			"token":           fmt.Sprintf("token-32-chars-long-and-secure-%d", ts),
			"workspace":       workspaceID,
			"user":            userA_ID,
			"scopes":          []string{util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}), http.StatusOK)
	})

	t.Run("Security: Regular member cannot create an API Key", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()
		api.AssertStatus(api.Create(util.Coll.ApiKeys, userB_Token, map[string]any{
			"username":        fmt.Sprintf("member-key-%d", ts),
			"email":           fmt.Sprintf("member-key-%d@test.com", ts),
			"token":           fmt.Sprintf("token-32-chars-long-and-secure-%d", ts),
			"workspace":       workspaceID,
			"user":            userB_ID,
			"scopes":          []string{util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}), http.StatusForbidden)
	})

	t.Run("Security: Admin cannot create an API Key for another user (SelfOnly check)", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()
		api.AssertStatus(api.Create(util.Coll.ApiKeys, userA_Token, map[string]any{
			"username":        fmt.Sprintf("spoofed-key-%d", ts),
			"email":           fmt.Sprintf("spoofed-key-%d@test.com", ts),
			"token":           fmt.Sprintf("token-32-chars-long-and-secure-%d", ts),
			"workspace":       workspaceID,
			"user":            userB_ID, // Admin trying to spoof user B
			"scopes":          []string{util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}), http.StatusForbidden)
	})

	t.Run("Security: API Key cannot create another API Key (Scope-less check)", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()
		apiKeyToken := fmt.Sprintf("initial-api-key-for-test-%d", ts)

		// Create base key
		api.Create(util.Coll.ApiKeys, userA_Token, map[string]any{
			"username":        fmt.Sprintf("escalation-base-%d", ts),
			"email":           fmt.Sprintf("escalation-base-%d@test.com", ts),
			"token":           apiKeyToken,
			"workspace":       workspaceID,
			"user":            userA_ID,
			"scopes":          []string{util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}).Expect().Status(http.StatusOK)

		// Try to use that key to create another key
		api.AssertStatus(api.Create(util.Coll.ApiKeys, apiKeyToken, map[string]any{
			"username":        fmt.Sprintf("escalated-key-%d", ts),
			"email":           fmt.Sprintf("escalated-key-%d@test.com", ts),
			"token":           fmt.Sprintf("token-32-chars-long-and-secure-%d", ts),
			"workspace":       workspaceID,
			"user":            userA_ID,
			"scopes":          []string{util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}), http.StatusForbidden)
	})

	t.Run("Security: API Key is immutable (No Updates)", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()
		token := fmt.Sprintf("immutable-test-token-%d", ts)

		res := api.Create(util.Coll.ApiKeys, userA_Token, map[string]any{
			"username":        fmt.Sprintf("immutable-%d", ts),
			"email":           fmt.Sprintf("immutable-%d@test.com", ts),
			"token":           token,
			"workspace":       workspaceID,
			"user":            userA_ID,
			"scopes":          []string{util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}).Expect().Status(http.StatusOK)

		actualID := testutils.ExtractString(res, "id")

		// Admin trying to update
		api.AssertStatus(api.Update(util.Coll.ApiKeys, actualID, userA_Token, map[string]any{
			"scopes": []string{util.ScopeDocumentsRead, util.ScopeDocumentsCreate},
		}), http.StatusForbidden)

		// API Key trying to update itself
		api.AssertStatus(api.Update(util.Coll.ApiKeys, actualID, token, map[string]any{
			"scopes": []string{util.ScopeDocumentsRead, util.ScopeDocumentsCreate},
		}), http.StatusForbidden)
	})

	t.Run("Security: API Key cannot be deleted by anyone except the owner", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()

		res := api.Create(util.Coll.ApiKeys, userA_Token, map[string]any{
			"username":        fmt.Sprintf("deletion-%d", ts),
			"email":           fmt.Sprintf("deletion-%d@test.com", ts),
			"token":           fmt.Sprintf("deletion-test-token-%d", ts),
			"workspace":       workspaceID,
			"user":            userA_ID,
			"scopes":          []string{util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}).Expect().Status(http.StatusOK)

		keyID := testutils.ExtractString(res, "id")

		api.AssertStatus(api.Delete(util.Coll.ApiKeys, keyID, userB_Token), http.StatusForbidden)
		api.AssertStatus(api.Delete(util.Coll.ApiKeys, keyID, userA_Token), http.StatusNoContent)
	})

	t.Run("Integrity: Cannot create an API Key with an invalid scope", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()
		api.AssertStatus(api.Create(util.Coll.ApiKeys, userA_Token, map[string]any{
			"username":        fmt.Sprintf("invalid-scope-%d", ts),
			"email":           fmt.Sprintf("invalid-scope-%d@test.com", ts),
			"token":           fmt.Sprintf("token-invalid-scope-%d", ts),
			"workspace":       workspaceID,
			"user":            userA_ID,
			"scopes":          []string{"invalid:scope:name"},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}), http.StatusBadRequest)
	})

	t.Run("Integrity: Cannot create an API Key with duplicate scopes", func(t *testing.T) {
		api := api.T(t)
		ts := time.Now().UnixNano()
		api.AssertStatus(api.Create(util.Coll.ApiKeys, userA_Token, map[string]any{
			"username":        fmt.Sprintf("duplicate-scope-%d", ts),
			"email":           fmt.Sprintf("duplicate-scope-%d@test.com", ts),
			"token":           fmt.Sprintf("token-duplicate-scope-%d", ts),
			"workspace":       workspaceID,
			"user":            userA_ID,
			"scopes":          []string{util.ScopeDocumentsRead, util.ScopeDocumentsRead},
			"password":        "password123456",
			"passwordConfirm": "password123456",
		}), http.StatusBadRequest)
	})
}
