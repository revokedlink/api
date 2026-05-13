package tests

import (
	"net/http"
	"revoked/tests/testutils"
	"revoked/util"
	"testing"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase/core"
)

func TestApiKeyScopes_Permissions(t *testing.T) {
	baseURL, app := testutils.SetupTestApp(t)
	userA, _ := app.FindCollectionByNameOrId(util.Coll.Users)
	recordA := core.NewRecord(userA)
	recordA.Set("email", "userA@test.com")
	recordA.Set("verified", true)
	recordA.SetPassword("password12345")
	if err := app.Save(recordA); err != nil {
		t.Fatalf("Failed to save User A: %v", err)
	}

	workspaces, _ := app.FindCollectionByNameOrId(util.Coll.Workspaces)
	wsA := core.NewRecord(workspaces)
	wsA.Set("name", "Workspace A")
	wsA.Set("type", util.TypeBusiness)
	wsA.Set("slug", "ws-a-"+uuid.New().String()[:8])
	if err := app.Save(wsA); err != nil {
		t.Fatalf("Failed to save Workspace A: %v", err)
	}

	apiKeys, _ := app.FindCollectionByNameOrId(util.Coll.ApiKeys)
	apiKeyAToken := uuid.New().String()
	akA := core.NewRecord(apiKeys)
	akA.Set("token", apiKeyAToken)
	akA.Set("user", recordA.Id)
	akA.Set("workspace", wsA.Id)
	akA.Set("scopes", []string{util.ScopeDocumentsCreate, util.ScopeDocumentsRead})
	akA.Set("email", "aka@test.com")
	akA.SetPassword("password12345")
	if err := app.Save(akA); err != nil {
		t.Fatalf("Failed to save API Key A: %v", err)
	}

	recordB := core.NewRecord(userA)
	recordB.Set("email", "userB@test.com")
	recordB.Set("verified", true)
	recordB.SetPassword("password12345")
	if err := app.Save(recordB); err != nil {
		t.Fatalf("Failed to save User B: %v", err)
	}

	wsB := core.NewRecord(workspaces)
	wsB.Set("name", "Workspace B")
	wsB.Set("type", util.TypeBusiness)
	wsB.Set("slug", "ws-b-"+uuid.New().String()[:8])
	if err := app.Save(wsB); err != nil {
		t.Fatalf("Failed to save Workspace B: %v", err)
	}

	apiKeyBToken := uuid.New().String()
	akB := core.NewRecord(apiKeys)
	akB.Set("token", apiKeyBToken)
	akB.Set("user", recordB.Id)
	akB.Set("workspace", wsB.Id)
	akB.Set("scopes", []string{util.ScopeDocumentsRead})
	akB.Set("email", "akb@test.com")
	akB.SetPassword("password12345")
	if err := app.Save(akB); err != nil {
		t.Fatalf("Failed to save API Key B: %v", err)
	}

	api := testutils.NewPBClient(t, baseURL)

	t.Run("User A creates document successfully", func(t *testing.T) {
		res := api.AssertStatus(api.Create("documents", apiKeyAToken, map[string]any{
			"workspace": wsA.Id,
			"title":     "Doc from A",
			"content":   "Hello world",
		}), http.StatusOK)

		res.JSON().Object().Value("title").String().IsEqual("Doc from A")
	})

	t.Run("User B fails to create document (missing scope)", func(t *testing.T) {
		api.AssertStatus(api.Create("documents", apiKeyBToken, map[string]any{
			"workspace": wsB.Id,
			"title":     "Doc from B",
			"content":   "Should fail",
		}), http.StatusBadRequest)
	})
	t.Run("User A fails to create document in Workspace B (workspace mismatch)", func(t *testing.T) {
		api.AssertStatus(api.Create("documents", apiKeyAToken, map[string]any{
			"workspace": wsB.Id,
			"title":     "Doc from A in B",
		}), http.StatusBadRequest)
	})
}
