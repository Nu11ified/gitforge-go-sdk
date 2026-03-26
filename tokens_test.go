package gitforge

import (
	"context"
	"encoding/json"
	"testing"
)

func TestTokensCreateRequired(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: RepoToken{
			Token:     "gf_abc123",
			PatID:     "pat-456",
			ExpiresAt: "2026-04-01T00:00:00Z",
			RemoteURL: "https://gitforge.example.com/repos/repo-123.git",
		},
	})

	tok := &TokensResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := tok.Create(context.Background(), &CreateTokenOptions{
		TTLSeconds: 3600,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "gf_abc123" {
		t.Errorf("expected token 'gf_abc123', got %q", result.Token)
	}
	if result.PatID != "pat-456" {
		t.Errorf("expected patId 'pat-456', got %q", result.PatID)
	}
	if result.ExpiresAt != "2026-04-01T00:00:00Z" {
		t.Errorf("expected expiresAt '2026-04-01T00:00:00Z', got %q", result.ExpiresAt)
	}
	if result.RemoteURL != "https://gitforge.example.com/repos/repo-123.git" {
		t.Errorf("expected remoteUrl, got %q", result.RemoteURL)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/tokens" {
		t.Errorf("expected path /repos/repo-123/tokens, got %s", call.Path)
	}
}

func TestTokensCreateAllOptions(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: RepoToken{
			Token:     "gf_xyz789",
			PatID:     "pat-999",
			ExpiresAt: "2026-05-01T00:00:00Z",
			RemoteURL: "https://gitforge.example.com/repos/repo-123.git",
		},
	})

	tok := &TokensResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := tok.Create(context.Background(), &CreateTokenOptions{
		TTLSeconds: 7200,
		Scopes:     []string{"read", "write"},
		Type:       ptrStr("ephemeral"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Token != "gf_xyz789" {
		t.Errorf("expected token 'gf_xyz789', got %q", result.Token)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/tokens" {
		t.Errorf("expected path /repos/repo-123/tokens, got %s", call.Path)
	}
}

func TestTokensCreateJSONBody(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: RepoToken{
			Token:     "gf_body_test",
			PatID:     "pat-body",
			ExpiresAt: "2026-06-01T00:00:00Z",
			RemoteURL: "https://gitforge.example.com/repos/repo-123.git",
		},
	})

	tok := &TokensResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	_, err := tok.Create(context.Background(), &CreateTokenOptions{
		TTLSeconds: 1800,
		Scopes:     []string{"read"},
		Type:       ptrStr("ci"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}

	// ttlSeconds must be an int (float64 after JSON decode)
	if body["ttlSeconds"] != float64(1800) {
		t.Errorf("expected ttlSeconds=1800, got %v", body["ttlSeconds"])
	}

	// scopes must be an array
	scopes, ok := body["scopes"].([]any)
	if !ok {
		t.Fatalf("expected scopes to be an array, got %T", body["scopes"])
	}
	if len(scopes) != 1 {
		t.Fatalf("expected 1 scope, got %d", len(scopes))
	}
	if scopes[0] != "read" {
		t.Errorf("expected scopes[0]='read', got %v", scopes[0])
	}

	// type must be a string
	if body["type"] != "ci" {
		t.Errorf("expected type='ci', got %v", body["type"])
	}
}
