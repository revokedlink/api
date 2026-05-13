package testutils

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

// PBClient wraps httpexpect to provide PocketBase-specific, declarative methods
type PBClient struct {
	E       *httpexpect.Expect
	baseURL string
	t       *testing.T
}

func NewPBClient(t *testing.T, baseURL string) *PBClient {
	return &PBClient{
		E:       NewExpect(t, baseURL),
		baseURL: baseURL,
		t:       t,
	}
}

// T returns a new PBClient scoped to the provided testing.T instance.
// This is essential for correct output attribution in subtests (t.Run).
func (c *PBClient) T(t *testing.T) *PBClient {
	return &PBClient{
		E:       NewExpect(t, c.baseURL),
		baseURL: c.baseURL,
		t:       t,
	}
}

// Request is a base builder that handles the repetitive "/api/collections/..." URL generation
func (c *PBClient) Request(method string, collection string, pathSuffix string) *httpexpect.Request {
	path := fmt.Sprintf("/api/collections/%s%s", collection, pathSuffix)
	return c.E.Request(method, path)
}

func (c *PBClient) Create(collection string, token string, body any) *httpexpect.Request {
	req := c.Request("POST", collection, "/records").WithJSON(body)
	return applyAuth(req, token)
}

func (c *PBClient) Get(collection string, id string, token string) *httpexpect.Request {
	req := c.Request("GET", collection, "/records/"+id)
	return applyAuth(req, token)
}

func (c *PBClient) Update(collection string, id string, token string, body any) *httpexpect.Request {
	req := c.Request("PATCH", collection, "/records/"+id).WithJSON(body)
	return applyAuth(req, token)
}

func (c *PBClient) Delete(collection string, id string, token string) *httpexpect.Request {
	req := c.Request("DELETE", collection, "/records/"+id)
	return applyAuth(req, token)
}

func (c *PBClient) AuthWithPassword(collection string, email, password string) *httpexpect.Request {
	return c.Request("POST", collection, "/auth-with-password").WithJSON(map[string]any{
		"identity": email,
		"password": password,
	})
}

// applyAuth intelligently applies the correct header based on token type
func applyAuth(req *httpexpect.Request, token string) *httpexpect.Request {
	if token == "" {
		return req
	}
	if len(token) > 100 {
		return req.WithHeader("Authorization", token)
	}
	return req.WithHeader("X-API-Key", token)
}

func (c *PBClient) AuthRefresh(collection string, token string) *httpexpect.Request {
	req := c.Request("POST", collection, "/auth-refresh")
	return applyAuth(req, token)
}

// AssertStatus wraps httpexpect's status check to handle PocketBase's specific API Rule quirks.
func (c *PBClient) AssertStatus(req *httpexpect.Request, expected int) *httpexpect.Response {
	resp := req.Expect()
	actual := resp.Raw().StatusCode

	if actual != expected {
		isPBQuirk := expected == http.StatusForbidden && (actual == http.StatusBadRequest || actual == http.StatusNotFound)

		if !isPBQuirk {
			c.t.Logf("\n--- FAILURE BODY ---\n%s\n--------------------\n", resp.JSON().Raw())
		}
	}

	if !(expected == http.StatusForbidden && (actual == http.StatusBadRequest || actual == http.StatusNotFound)) {
		resp.Status(expected)
	}

	return resp
}
