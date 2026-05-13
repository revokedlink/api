package testutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"revoked/util"
	"time"

	"github.com/gavv/httpexpect/v2"
	"github.com/google/uuid"
)

// CreateRandomUser creates a new user with a random email and password via the HTTP API,
// and returns the user's record ID and their associated JWT auth token.
//
// This is primarily intended for integration tests where a real JWT is needed
// to authenticate subsequent API requests.
func CreateRandomUser(baseURL string) (id string, token string, err error) {
	email := fmt.Sprintf("test-%s@example.com", uuid.New().String()[:8])
	password := "password12345"

	createData := map[string]any{
		"email":           email,
		"password":        password,
		"passwordConfirm": password,
	}
	createBody, _ := json.Marshal(createData)

	createURL := fmt.Sprintf("%s/api/collections/%s/records", baseURL, util.Coll.Users)
	resp, err := http.Post(createURL, "application/json", bytes.NewBuffer(createBody))
	if err != nil {
		return "", "", fmt.Errorf("failed to send create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errBody any
		json.NewDecoder(resp.Body).Decode(&errBody)
		return "", "", fmt.Errorf("create user failed with status %d: %v", resp.StatusCode, errBody)
	}

	var createResult struct {
		Id string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&createResult); err != nil {
		return "", "", fmt.Errorf("failed to decode create response: %w", err)
	}

	authData := map[string]any{
		"identity": email,
		"password": password,
	}
	authBody, _ := json.Marshal(authData)
	authURL := fmt.Sprintf("%s/api/collections/%s/auth-with-password", baseURL, util.Coll.Users)

	for i := 0; i < 5; i++ {
		authResp, err := http.Post(authURL, "application/json", bytes.NewBuffer(authBody))
		if err == nil && authResp.StatusCode == http.StatusOK {
			var authResult struct {
				Token string `json:"token"`
			}
			if err := json.NewDecoder(authResp.Body).Decode(&authResult); err == nil {
				authResp.Body.Close()
				return createResult.Id, authResult.Token, nil
			}
		}
		if authResp != nil {
			authResp.Body.Close()
		}
		time.Sleep(200 * time.Millisecond)
	}

	return "", "", fmt.Errorf("failed to authenticate user after creation at %s", authURL)
}

// ExtractString is a shorthand to grab a top-level string field from a JSON response
func ExtractString(res *httpexpect.Response, key string) string {
	return res.JSON().Object().Value(key).String().Raw()
}

// List fetches a paginated list of records from a collection
func (c *PBClient) List(collection string, token string) *httpexpect.Request {
	req := c.Request("GET", collection, "/records")
	return applyAuth(req, token)
}
