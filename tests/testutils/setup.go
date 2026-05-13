package testutils

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"revoked/cmd/revoked/hooks"
	_ "revoked/migrations"
	"revoked/util"
	"sync"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/stretchr/testify/assert"
)

var (
	testApp                *pocketbase.PocketBase
	testAppURL             = "127.0.0.1:5559"
	testAppUrlWithProtocol = "http://127.0.0.1:5559"
	serverOnce             sync.Once
	serverCtx              context.Context
	cancelServer           context.CancelFunc
)

// SetupTestApp initializes a shared PocketBase instance for the test suite
func SetupTestApp(t testing.TB) (string, *pocketbase.PocketBase) {
	t.Helper()

	serverOnce.Do(func() {
		testDataDir := "./pb_test_data"
		_ = os.RemoveAll(testDataDir)

		testApp = pocketbase.NewWithConfig(pocketbase.Config{
			DefaultDataDir: testDataDir,
		})

		migratecmd.MustRegister(testApp, testApp.RootCmd, migratecmd.Config{
			Automigrate: true,
		})

		if err := testApp.Bootstrap(); err != nil {
			t.Fatalf("Failed to bootstrap app: %v", err)
		}

		os.Args = []string{"pb", "migrate", "up"}
		if err := testApp.RootCmd.ExecuteContext(context.Background()); err != nil {
			fmt.Printf("Migration notice: %v\n", err)
		}

		hooks.BindHooksAndRoutes(testApp)

		serverCtx, cancelServer = context.WithCancel(context.Background())

		os.Args = []string{"pb", "serve", "--http=" + testAppURL}

		go func() {
			if err := testApp.Start(); err != nil {
				fmt.Printf("Server stopped: %v\n", err)
			}
		}()
		print(fmt.Sprintf("http://%s/healthz", testAppURL))

		waitForHealthy(t, fmt.Sprintf("http://%s/healthz", testAppURL))
	})

	return fmt.Sprintf("http://%s", testAppURL), testApp
}

func waitForHealthy(t testing.TB, url string) {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			t.Fatal("Timeout: PocketBase server failed to start")
		case <-ticker.C:
			resp, err := http.Get(url)
			if err == nil && resp.StatusCode == http.StatusOK {
				return
			}
		}
	}
}

// ClearCollections wipes records from specific collections using the modern API
func ClearCollections(t testing.TB, app *pocketbase.PocketBase, names ...string) {
	t.Helper()

	for _, name := range names {
		collection, err := app.FindCollectionByNameOrId(name)
		if err != nil {
			t.Fatalf("Collection %s not found: %v", name, err)
		}

		if err := app.TruncateCollection(collection); err != nil {
			t.Fatalf("Failed to truncate %s: %v", name, err)
		}
	}
}

// ClearAllCustomData loops through your app's collections and wipes them
func ClearAllCustomData(t testing.TB, app *pocketbase.PocketBase) {
	t.Helper()

	collections := []string{"workspaces", "workspace_members", "users"}
	ClearCollections(t, app, collections...)
}

func AssertBadRequestErrors(t *testing.T, resp *httpexpect.Response, expectations map[string]util.AppError) {
	resp.Status(http.StatusBadRequest)
	body := resp.JSON().Object()

	body.Value("status").Number().IsEqual(http.StatusBadRequest)
	actualMsg := resp.JSON().Object().Value("message").String().Raw()
	if actualMsg != "Something went wrong while processing your request." && actualMsg != util.Errors.PersonalWorkspaceLimitReached.ErrorText && actualMsg != util.Errors.BusinessWorkspaceLimitReached.ErrorText {
		assert.Equal(t, "Something went wrong while processing your request.", actualMsg)
	}

	data := body.Value("data").Object()

	for key, appErr := range expectations {
		errorContext := data.Value(key).Object()
		errorContext.Value("code").String().IsEqual(appErr.ErrorCode)
		errorContext.Value("message").String().IsEqual(appErr.ErrorText)
	}
}

// AssertErrorResponse validates a standard top-level error response (like 401, 403, or 404).
// It verifies both the HTTP Header status and the JSON body fields.
func AssertErrorResponse(t *testing.T, resp *httpexpect.Response, expectedStatus int, expectedMessage string) {
	resp.Status(expectedStatus)
	body := resp.JSON().Object()

	body.Value("status").Number().IsEqual(expectedStatus)
	body.Value("message").String().IsEqual(expectedMessage)
}

// NewExpect function to generate a scoped httpexpect instance
func NewExpect(t *testing.T, baseURL string) *httpexpect.Expect {
	return httpexpect.WithConfig(httpexpect.Config{
		BaseURL:  baseURL,
		Reporter: httpexpect.NewAssertReporter(t),
		Printers: []httpexpect.Printer{
			httpexpect.NewCompactPrinter(t),
		},
	})
}
