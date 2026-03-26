package gitforge

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

// TestHTTPBearerTokenGet verifies the Authorization header is sent on GET requests.
func TestHTTPBearerTokenGet(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: map[string]any{"ok": true}})
	c := newHTTPClient(srv.URL, "mytoken", nil)

	_, err := c.get(context.Background(), "/repos", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	auth := (*calls)[0].Headers.Get("Authorization")
	if auth != "Bearer mytoken" {
		t.Errorf("expected Authorization: Bearer mytoken, got %q", auth)
	}
}

// TestHTTPBearerTokenPost verifies the Authorization header is sent on POST requests.
func TestHTTPBearerTokenPost(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 201, Body: map[string]any{"id": "1"}})
	c := newHTTPClient(srv.URL, "posttoken", nil)

	_, err := c.post(context.Background(), "/repos", map[string]string{"name": "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	auth := (*calls)[0].Headers.Get("Authorization")
	if auth != "Bearer posttoken" {
		t.Errorf("expected Bearer posttoken, got %q", auth)
	}
}

// TestHTTPBearerTokenPatch verifies the Authorization header is sent on PATCH requests.
func TestHTTPBearerTokenPatch(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: map[string]any{"updated": true}})
	c := newHTTPClient(srv.URL, "patchtoken", nil)

	_, err := c.patch(context.Background(), "/repos/1", map[string]string{"name": "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	auth := (*calls)[0].Headers.Get("Authorization")
	if auth != "Bearer patchtoken" {
		t.Errorf("expected Bearer patchtoken, got %q", auth)
	}
}

// TestHTTPBearerTokenDelete verifies the Authorization header is sent on DELETE requests.
func TestHTTPBearerTokenDelete(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204, Body: nil})
	c := newHTTPClient(srv.URL, "deltoken", nil)

	err := c.del(context.Background(), "/repos/1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	auth := (*calls)[0].Headers.Get("Authorization")
	if auth != "Bearer deltoken" {
		t.Errorf("expected Bearer deltoken, got %q", auth)
	}
}

// TestHTTPURLBuilding verifies URLs are built from baseURL + path.
func TestHTTPURLBuilding(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: map[string]any{}})
	c := newHTTPClient(srv.URL, "tok", nil)

	_, err := c.get(context.Background(), "/api/v1/repos", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*calls)[0].Path != "/api/v1/repos" {
		t.Errorf("expected path /api/v1/repos, got %q", (*calls)[0].Path)
	}
}

// TestHTTPTrailingSlashStripped verifies trailing slashes are stripped from baseURL.
func TestHTTPTrailingSlashStripped(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: map[string]any{}})
	// Append trailing slashes to the server URL
	c := newHTTPClient(srv.URL+"/", "tok", nil)

	_, err := c.get(context.Background(), "/repos", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if (*calls)[0].Path != "/repos" {
		t.Errorf("expected path /repos, got %q", (*calls)[0].Path)
	}
}

// TestHTTPQueryParamsAppended verifies query params are appended on GET.
func TestHTTPQueryParamsAppended(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: map[string]any{}})
	c := newHTTPClient(srv.URL, "tok", nil)

	q := url.Values{}
	q.Set("limit", "10")
	q.Set("offset", "20")

	_, err := c.get(context.Background(), "/repos", q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	call := (*calls)[0]
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("limit") != "10" {
		t.Errorf("expected limit=10, got %q", parsed.Get("limit"))
	}
	if parsed.Get("offset") != "20" {
		t.Errorf("expected offset=20, got %q", parsed.Get("offset"))
	}
}

// TestHTTPPostJSONBody verifies POST sends JSON body with Content-Type header.
func TestHTTPPostJSONBody(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 201, Body: map[string]any{"id": "abc"}})
	c := newHTTPClient(srv.URL, "tok", nil)

	payload := map[string]string{"name": "myrepo", "visibility": "private"}
	_, err := c.post(context.Background(), "/repos", payload)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %q", call.Method)
	}
	ct := call.Headers.Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected Content-Type: application/json, got %q", ct)
	}

	var decoded map[string]string
	if err := json.Unmarshal(call.Body, &decoded); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if decoded["name"] != "myrepo" {
		t.Errorf("expected name=myrepo, got %q", decoded["name"])
	}
	if decoded["visibility"] != "private" {
		t.Errorf("expected visibility=private, got %q", decoded["visibility"])
	}
}

// TestHTTPDeleteReturnsNilOn204 verifies DELETE returns nil error on 204.
func TestHTTPDeleteReturnsNilOn204(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{Status: 204, Body: nil})
	c := newHTTPClient(srv.URL, "tok", nil)

	err := c.del(context.Background(), "/repos/abc", nil)
	if err != nil {
		t.Errorf("expected nil error on 204, got %v", err)
	}
}

// TestHTTPDeleteQueryParams verifies DELETE supports query params.
func TestHTTPDeleteQueryParams(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204, Body: nil})
	c := newHTTPClient(srv.URL, "tok", nil)

	q := url.Values{}
	q.Set("force", "true")
	err := c.del(context.Background(), "/repos/abc", q)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("force") != "true" {
		t.Errorf("expected force=true, got %q", parsed.Get("force"))
	}
}

// TestHTTPErrorMappingGitForgeError verifies non-2xx responses map to *GitForgeError.
func TestHTTPErrorMappingGitForgeError(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{
		Status: 403,
		Body:   map[string]any{"code": "forbidden", "message": "access denied"},
	})
	c := newHTTPClient(srv.URL, "tok", nil)

	_, err := c.get(context.Background(), "/repos", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var gfe *GitForgeError
	if !errors.As(err, &gfe) {
		t.Fatalf("expected *GitForgeError, got %T", err)
	}
	if gfe.StatusCode != 403 {
		t.Errorf("expected StatusCode 403, got %d", gfe.StatusCode)
	}
	if gfe.Code != "forbidden" {
		t.Errorf("expected Code forbidden, got %q", gfe.Code)
	}
	if gfe.Message != "access denied" {
		t.Errorf("expected Message 'access denied', got %q", gfe.Message)
	}
}

// TestHTTPErrorCodeKey verifies the "code" field is preferred for the error code.
func TestHTTPErrorCodeKey(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{
		Status: 404,
		Body:   map[string]any{"code": "not_found", "message": "repo not found"},
	})
	c := newHTTPClient(srv.URL, "tok", nil)

	_, err := c.get(context.Background(), "/repos/missing", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var gfe *GitForgeError
	if !errors.As(err, &gfe) {
		t.Fatalf("expected *GitForgeError, got %T", err)
	}
	if gfe.StatusCode != 404 {
		t.Errorf("expected StatusCode 404, got %d", gfe.StatusCode)
	}
	if gfe.Code != "not_found" {
		t.Errorf("expected Code not_found, got %q", gfe.Code)
	}
}

// TestHTTPErrorErrorKey verifies the "error" field is used as fallback for error code.
func TestHTTPErrorErrorKey(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{
		Status: 422,
		Body:   map[string]any{"error": "validation_failed", "message": "invalid input"},
	})
	c := newHTTPClient(srv.URL, "tok", nil)

	_, err := c.get(context.Background(), "/repos", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var gfe *GitForgeError
	if !errors.As(err, &gfe) {
		t.Fatalf("expected *GitForgeError, got %T", err)
	}
	if gfe.Code != "validation_failed" {
		t.Errorf("expected Code validation_failed, got %q", gfe.Code)
	}
	if gfe.Message != "invalid input" {
		t.Errorf("expected Message 'invalid input', got %q", gfe.Message)
	}
}

// TestHTTP409BranchMoved verifies 409 with branch_moved maps to *RefUpdateError.
func TestHTTP409BranchMoved(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{
		Status: 409,
		Body: map[string]any{
			"error":      "branch_moved",
			"currentSha": "abc123",
			"message":    "msg",
		},
	})
	c := newHTTPClient(srv.URL, "tok", nil)

	_, err := c.post(context.Background(), "/repos/r/branches/main", map[string]string{"sha": "old"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var rue *RefUpdateError
	if !errors.As(err, &rue) {
		t.Fatalf("expected *RefUpdateError, got %T: %v", err, err)
	}
	if rue.CurrentSHA != "abc123" {
		t.Errorf("expected CurrentSHA abc123, got %q", rue.CurrentSHA)
	}
	if rue.StatusCode != 409 {
		t.Errorf("expected StatusCode 409, got %d", rue.StatusCode)
	}
	if rue.Code != "branch_moved" {
		t.Errorf("expected Code branch_moved, got %q", rue.Code)
	}
	if rue.Message != "msg" {
		t.Errorf("expected Message msg, got %q", rue.Message)
	}
}

// TestHTTPNonJSONErrorBody verifies non-JSON error body returns GitForgeError with "unknown" code.
func TestHTTPNonJSONErrorBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(502)
		_, _ = w.Write([]byte("Bad Gateway"))
	}))
	t.Cleanup(srv.Close)

	c := newHTTPClient(srv.URL, "tok", nil)
	_, err := c.get(context.Background(), "/repos", nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	var gfe *GitForgeError
	if !errors.As(err, &gfe) {
		t.Fatalf("expected *GitForgeError, got %T", err)
	}
	if gfe.Code != "unknown" {
		t.Errorf("expected Code unknown, got %q", gfe.Code)
	}
	if gfe.StatusCode != 502 {
		t.Errorf("expected StatusCode 502, got %d", gfe.StatusCode)
	}
}

// TestHTTP404ErrorCode verifies 404 error code is mapped correctly.
func TestHTTP404ErrorCode(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{
		Status: 404,
		Body:   map[string]any{"code": "not_found", "message": "not found"},
	})
	c := newHTTPClient(srv.URL, "tok", nil)

	_, err := c.get(context.Background(), "/missing", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var gfe *GitForgeError
	if !errors.As(err, &gfe) {
		t.Fatalf("expected *GitForgeError, got %T", err)
	}
	if gfe.StatusCode != 404 {
		t.Errorf("expected 404, got %d", gfe.StatusCode)
	}
}

// TestHTTP500ErrorCode verifies 500 error code is mapped correctly.
func TestHTTP500ErrorCode(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{
		Status: 500,
		Body:   map[string]any{"code": "internal_error", "message": "server exploded"},
	})
	c := newHTTPClient(srv.URL, "tok", nil)

	_, err := c.get(context.Background(), "/repos", nil)
	if err == nil {
		t.Fatal("expected error")
	}
	var gfe *GitForgeError
	if !errors.As(err, &gfe) {
		t.Fatalf("expected *GitForgeError, got %T", err)
	}
	if gfe.StatusCode != 500 {
		t.Errorf("expected 500, got %d", gfe.StatusCode)
	}
	if gfe.Code != "internal_error" {
		t.Errorf("expected internal_error, got %q", gfe.Code)
	}
}

// TestHTTPSuccessParsedJSON verifies 200 response returns parsed JSON.
func TestHTTPSuccessParsedJSON(t *testing.T) {
	expected := map[string]any{"id": "repo1", "name": "myrepo"}
	srv, _ := setupMock(t, mockResponse{Status: 200, Body: expected})
	c := newHTTPClient(srv.URL, "tok", nil)

	raw, err := c.get(context.Background(), "/repos/repo1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if raw == nil {
		t.Fatal("expected non-nil response body")
	}

	var decoded map[string]any
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if decoded["id"] != "repo1" {
		t.Errorf("expected id=repo1, got %v", decoded["id"])
	}
	if decoded["name"] != "myrepo" {
		t.Errorf("expected name=myrepo, got %v", decoded["name"])
	}
}

// TestHTTPContextCancellation verifies context cancellation is propagated.
func TestHTTPContextCancellation(t *testing.T) {
	srv, _ := setupMock(t, mockResponse{Status: 200, Body: map[string]any{}})
	c := newHTTPClient(srv.URL, "tok", nil)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	_, err := c.get(ctx, "/repos", nil)
	if err == nil {
		t.Fatal("expected error from cancelled context, got nil")
	}
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}
