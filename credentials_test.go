package gitforge

import (
	"context"
	"encoding/json"
	"testing"
)

func TestCredentialsCreate(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: GitCredential{
			ID:        "cred-abc",
			Provider:  "github",
			CreatedAt: "2026-03-26T00:00:00Z",
		},
	})

	c := &CredentialsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	cred, err := c.Create(context.Background(), &CreateCredentialOptions{
		Provider: "github",
		Token:    "ghp_secret",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.ID != "cred-abc" {
		t.Errorf("expected id 'cred-abc', got %q", cred.ID)
	}
	if cred.Provider != "github" {
		t.Errorf("expected provider 'github', got %q", cred.Provider)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/credentials" {
		t.Errorf("expected path /repos/repo-123/credentials, got %s", call.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["provider"] != "github" {
		t.Errorf("expected provider 'github' in body, got %v", body["provider"])
	}
	if body["token"] != "ghp_secret" {
		t.Errorf("expected token 'ghp_secret' in body, got %v", body["token"])
	}
}

func TestCredentialsCreateWithAllOptions(t *testing.T) {
	label := "My GitHub Credential"
	username := "octocat"
	sourceURL := "https://github.com"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: GitCredential{
			ID:        "cred-xyz",
			Provider:  "github",
			Username:  &username,
			Label:     &label,
			CreatedAt: "2026-03-26T00:00:00Z",
		},
	})

	c := &CredentialsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	cred, err := c.Create(context.Background(), &CreateCredentialOptions{
		Provider:  "github",
		Token:     "ghp_secret2",
		Username:  ptrStr(username),
		Label:     ptrStr(label),
		SourceURL: ptrStr(sourceURL),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.ID != "cred-xyz" {
		t.Errorf("expected id 'cred-xyz', got %q", cred.ID)
	}
	if cred.Username == nil || *cred.Username != username {
		t.Errorf("expected username %q, got %v", username, cred.Username)
	}
	if cred.Label == nil || *cred.Label != label {
		t.Errorf("expected label %q, got %v", label, cred.Label)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["username"] != username {
		t.Errorf("expected username %q in body, got %v", username, body["username"])
	}
	if body["label"] != label {
		t.Errorf("expected label %q in body, got %v", label, body["label"])
	}
	if body["sourceUrl"] != sourceURL {
		t.Errorf("expected sourceUrl %q in body, got %v", sourceURL, body["sourceUrl"])
	}
}

func TestCredentialsList(t *testing.T) {
	username := "octocat"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: []GitCredential{
			{ID: "cred-1", Provider: "github", Username: &username, CreatedAt: "2026-03-26T00:00:00Z"},
			{ID: "cred-2", Provider: "gitlab", CreatedAt: "2026-03-25T00:00:00Z"},
		},
	})

	c := &CredentialsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	creds, err := c.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(creds) != 2 {
		t.Fatalf("expected 2 credentials, got %d", len(creds))
	}
	if creds[0].ID != "cred-1" {
		t.Errorf("expected first cred id 'cred-1', got %q", creds[0].ID)
	}
	if creds[0].Provider != "github" {
		t.Errorf("expected first cred provider 'github', got %q", creds[0].Provider)
	}
	if creds[1].ID != "cred-2" {
		t.Errorf("expected second cred id 'cred-2', got %q", creds[1].ID)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/credentials" {
		t.Errorf("expected path /repos/repo-123/credentials, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestCredentialsUpdate(t *testing.T) {
	newLabel := "Updated Label"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: GitCredential{
			ID:        "cred-abc",
			Provider:  "github",
			Label:     &newLabel,
			CreatedAt: "2026-03-26T00:00:00Z",
		},
	})

	c := &CredentialsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	cred, err := c.Update(context.Background(), "cred-abc", &UpdateCredentialOptions{
		Label: ptrStr(newLabel),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred.ID != "cred-abc" {
		t.Errorf("expected id 'cred-abc', got %q", cred.ID)
	}
	if cred.Label == nil || *cred.Label != newLabel {
		t.Errorf("expected label %q, got %v", newLabel, cred.Label)
	}

	call := (*calls)[0]
	if call.Method != "PATCH" {
		t.Errorf("expected PATCH, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/credentials/cred-abc" {
		t.Errorf("expected path /repos/repo-123/credentials/cred-abc, got %s", call.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["label"] != newLabel {
		t.Errorf("expected label %q in body, got %v", newLabel, body["label"])
	}
}

func TestCredentialsDelete(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204})

	c := &CredentialsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	err := c.Delete(context.Background(), "cred-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	if call.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/credentials/cred-abc" {
		t.Errorf("expected path /repos/repo-123/credentials/cred-abc, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}
