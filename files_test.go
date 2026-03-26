package gitforge

import (
	"context"
	"net/url"
	"testing"
)

// TestFilesListDefault verifies GET /repos/{id}/tree/{ref} returns []TreeEntry.
func TestFilesListDefault(t *testing.T) {
	entries := []TreeEntry{
		{Name: "README.md", Type: "blob", Mode: "100644", SHA: "abc123"},
		{Name: "src", Type: "tree", Mode: "040000", SHA: "def456"},
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: entries})
	f := &FilesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	result, err := f.ListFiles(context.Background(), "main", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(result))
	}
	if result[0].Name != "README.md" {
		t.Errorf("expected Name README.md, got %q", result[0].Name)
	}
	if result[0].Type != "blob" {
		t.Errorf("expected Type blob, got %q", result[0].Type)
	}
	if result[1].Name != "src" {
		t.Errorf("expected Name src, got %q", result[1].Name)
	}
	if result[1].Type != "tree" {
		t.Errorf("expected Type tree, got %q", result[1].Type)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %q", call.Method)
	}
	if call.Path != "/repos/repo-123/tree/main" {
		t.Errorf("expected path /repos/repo-123/tree/main, got %q", call.Path)
	}
}

// TestFilesListWithPath verifies ListFiles with Path includes path query param.
func TestFilesListWithPath(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: []TreeEntry{}})
	f := &FilesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := f.ListFiles(context.Background(), "main", &ListFilesOptions{Path: ptrStr("src/lib")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("path") != "src/lib" {
		t.Errorf("expected path=src/lib, got %q", parsed.Get("path"))
	}
}

// TestFilesListEphemeral verifies ListFiles with Ephemeral=true includes ephemeral=true query param.
func TestFilesListEphemeral(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: []TreeEntry{}})
	f := &FilesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := f.ListFiles(context.Background(), "eph-branch", &ListFilesOptions{Ephemeral: ptrBool(true)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("ephemeral") != "true" {
		t.Errorf("expected ephemeral=true, got %q", parsed.Get("ephemeral"))
	}
}

// TestFilesGetFile verifies GET /repos/{id}/blob/{ref} with path query returns *BlobContent.
func TestFilesGetFile(t *testing.T) {
	blob := BlobContent{Content: "aGVsbG8gd29ybGQ=", Size: 11}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: blob})
	f := &FilesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	result, err := f.GetFile(context.Background(), "main", "README.md", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Content != "aGVsbG8gd29ybGQ=" {
		t.Errorf("expected Content aGVsbG8gd29ybGQ=, got %q", result.Content)
	}
	if result.Size != 11 {
		t.Errorf("expected Size 11, got %d", result.Size)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %q", call.Method)
	}
	if call.Path != "/repos/repo-123/blob/main" {
		t.Errorf("expected path /repos/repo-123/blob/main, got %q", call.Path)
	}
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("path") != "README.md" {
		t.Errorf("expected path=README.md, got %q", parsed.Get("path"))
	}
}

// TestFilesGetFileEphemeral verifies GetFile with Ephemeral=true includes ephemeral=true query param.
func TestFilesGetFileEphemeral(t *testing.T) {
	blob := BlobContent{Content: "c29tZWNvbnRlbnQ=", Size: 11}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: blob})
	f := &FilesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := f.GetFile(context.Background(), "eph-branch", "config.json", &GetFileOptions{Ephemeral: ptrBool(true)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("ephemeral") != "true" {
		t.Errorf("expected ephemeral=true, got %q", parsed.Get("ephemeral"))
	}
	if parsed.Get("path") != "config.json" {
		t.Errorf("expected path=config.json, got %q", parsed.Get("path"))
	}
}
