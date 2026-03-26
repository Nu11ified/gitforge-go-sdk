package gitforge

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

// TestCommitList verifies GET /repos/{id}/commits returns []Commit (not paginated).
func TestCommitList(t *testing.T) {
	commits := []Commit{
		{SHA: "abc123", Message: "first commit", Author: "Alice", AuthorEmail: "alice@example.com", Date: "2026-01-01T00:00:00Z"},
		{SHA: "def456", Message: "second commit", Author: "Bob", AuthorEmail: "bob@example.com", Date: "2026-01-02T00:00:00Z"},
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: commits})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	result, err := c.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 commits, got %d", len(result))
	}
	if result[0].SHA != "abc123" {
		t.Errorf("expected SHA abc123, got %q", result[0].SHA)
	}
	if result[1].SHA != "def456" {
		t.Errorf("expected SHA def456, got %q", result[1].SHA)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %q", call.Method)
	}
	if call.Path != "/repos/repo-123/commits" {
		t.Errorf("expected path /repos/repo-123/commits, got %q", call.Path)
	}
}

// TestCommitListWithRef verifies List with Ref includes the ref query param.
func TestCommitListWithRef(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: []Commit{}})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := c.List(context.Background(), &ListCommitsOptions{Ref: ptrStr("main")})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("ref") != "main" {
		t.Errorf("expected ref=main, got %q", parsed.Get("ref"))
	}
}

// TestCommitListWithEphemeral verifies List with Ephemeral=true includes ephemeral=true query param.
func TestCommitListWithEphemeral(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: []Commit{}})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := c.List(context.Background(), &ListCommitsOptions{Ephemeral: ptrBool(true)})
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

// TestCommitListWithLimit verifies List with Limit includes limit query param.
func TestCommitListWithLimit(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: []Commit{}})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := c.List(context.Background(), &ListCommitsOptions{Limit: ptrInt(25)})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	parsed, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if parsed.Get("limit") != "25" {
		t.Errorf("expected limit=25, got %q", parsed.Get("limit"))
	}
}

// TestCommitGet verifies GET /repos/{id}/commits/{sha} returns *CommitDetail.
func TestCommitGet(t *testing.T) {
	detail := CommitDetail{
		SHA:         "abc123",
		Message:     "fix: something",
		Author:      "Alice",
		AuthorEmail: "alice@example.com",
		Date:        "2026-01-01T00:00:00Z",
		Tree:        "tree456",
		Files: []FileChange{
			{Path: "README.md", Status: "modified"},
		},
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: detail})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	result, err := c.Get(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.SHA != "abc123" {
		t.Errorf("expected SHA abc123, got %q", result.SHA)
	}
	if result.Tree != "tree456" {
		t.Errorf("expected Tree tree456, got %q", result.Tree)
	}
	if len(result.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(result.Files))
	}
	if result.Files[0].Path != "README.md" {
		t.Errorf("expected file README.md, got %q", result.Files[0].Path)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %q", call.Method)
	}
	if call.Path != "/repos/repo-123/commits/abc123" {
		t.Errorf("expected path /repos/repo-123/commits/abc123, got %q", call.Path)
	}
}

// TestCommitGetDiff verifies GET /repos/{id}/commits/{sha}/diff returns []DiffEntry.
func TestCommitGetDiff(t *testing.T) {
	entries := []DiffEntry{
		{Path: "main.go", Status: "modified", Additions: 10, Deletions: 2, Patch: "@@ -1,2 +1,10 @@"},
		{Path: "new.go", Status: "added", Additions: 50, Deletions: 0, Patch: "@@ -0,0 +1,50 @@"},
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: entries})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	result, err := c.GetDiff(context.Background(), "abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 diff entries, got %d", len(result))
	}
	if result[0].Path != "main.go" {
		t.Errorf("expected path main.go, got %q", result[0].Path)
	}
	if result[1].Status != "added" {
		t.Errorf("expected status added, got %q", result[1].Status)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %q", call.Method)
	}
	if call.Path != "/repos/repo-123/commits/abc123/diff" {
		t.Errorf("expected path /repos/repo-123/commits/abc123/diff, got %q", call.Path)
	}
}

// TestCommitBuilderAddFileChaining verifies AddFile chaining sends POST with correct body.
func TestCommitBuilderAddFileChaining(t *testing.T) {
	result := CommitResult{CommitSHA: "newsha", Branch: "main", NewSHA: "newsha", OldSHA: "oldsha"}
	srv, calls := setupMock(t, mockResponse{Status: 201, Body: result})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	res, err := c.Create(&CreateCommitOptions{
		Branch:      "main",
		Message:     "add two files",
		AuthorName:  "Alice",
		AuthorEmail: "alice@example.com",
	}).
		AddFile("README.md", "# Hello", nil).
		AddFile("main.go", "package main", &FileOptions{Mode: "100644"}).
		Send(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.CommitSHA != "newsha" {
		t.Errorf("expected CommitSHA newsha, got %q", res.CommitSHA)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %q", call.Method)
	}
	if call.Path != "/repos/repo-123/commits" {
		t.Errorf("expected path /repos/repo-123/commits, got %q", call.Path)
	}

	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	if body["branch"] != "main" {
		t.Errorf("expected branch=main, got %v", body["branch"])
	}
	if body["message"] != "add two files" {
		t.Errorf("expected message='add two files', got %v", body["message"])
	}

	files, ok := body["files"].([]any)
	if !ok {
		t.Fatalf("expected files to be array, got %T", body["files"])
	}
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	f0 := files[0].(map[string]any)
	if f0["path"] != "README.md" {
		t.Errorf("expected first file path README.md, got %v", f0["path"])
	}
	if f0["content"] != "# Hello" {
		t.Errorf("expected first file content '# Hello', got %v", f0["content"])
	}

	f1 := files[1].(map[string]any)
	if f1["path"] != "main.go" {
		t.Errorf("expected second file path main.go, got %v", f1["path"])
	}
	if f1["mode"] != "100644" {
		t.Errorf("expected second file mode 100644, got %v", f1["mode"])
	}
}

// TestCommitBuilderDeleteFile verifies DeleteFile sends deletes array in body.
func TestCommitBuilderDeleteFile(t *testing.T) {
	result := CommitResult{CommitSHA: "delsha", Branch: "main"}
	srv, calls := setupMock(t, mockResponse{Status: 201, Body: result})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := c.Create(&CreateCommitOptions{
		Branch:      "main",
		Message:     "delete file",
		AuthorName:  "Alice",
		AuthorEmail: "alice@example.com",
	}).
		DeleteFile("old.go").
		Send(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	deletes, ok := body["deletes"].([]any)
	if !ok {
		t.Fatalf("expected deletes to be array, got %T", body["deletes"])
	}
	if len(deletes) != 1 {
		t.Fatalf("expected 1 delete entry, got %d", len(deletes))
	}
	if deletes[0] != "old.go" {
		t.Errorf("expected delete path old.go, got %v", deletes[0])
	}
}

// TestCommitBuilderEphemeral verifies Ephemeral(true) includes "ephemeral":true in body.
func TestCommitBuilderEphemeral(t *testing.T) {
	result := CommitResult{CommitSHA: "epsha", Branch: "eph-branch"}
	srv, calls := setupMock(t, mockResponse{Status: 201, Body: result})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := c.Create(&CreateCommitOptions{
		Branch:      "eph-branch",
		Message:     "ephemeral commit",
		AuthorName:  "Alice",
		AuthorEmail: "alice@example.com",
	}).
		Ephemeral(true).
		Send(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	ephemeral, ok := body["ephemeral"]
	if !ok {
		t.Fatal("expected ephemeral field in body")
	}
	if ephemeral != true {
		t.Errorf("expected ephemeral=true, got %v", ephemeral)
	}
}

// TestCommitBuilderExpectedHeadSHA verifies ExpectedHeadSHA includes "expectedHeadSha" in body.
func TestCommitBuilderExpectedHeadSHA(t *testing.T) {
	result := CommitResult{CommitSHA: "newsha", Branch: "main"}
	srv, calls := setupMock(t, mockResponse{Status: 201, Body: result})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	_, err := c.Create(&CreateCommitOptions{
		Branch:      "main",
		Message:     "safe commit",
		AuthorName:  "Alice",
		AuthorEmail: "alice@example.com",
	}).
		ExpectedHeadSHA("headsha123").
		Send(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	headSHA, ok := body["expectedHeadSha"]
	if !ok {
		t.Fatal("expected expectedHeadSha field in body")
	}
	if headSHA != "headsha123" {
		t.Errorf("expected expectedHeadSha=headsha123, got %v", headSHA)
	}
}

// TestCommitBuilderWithBaseBranch verifies Create with BaseBranch includes "baseBranch" in body.
func TestCommitBuilderWithBaseBranch(t *testing.T) {
	result := CommitResult{CommitSHA: "branchsha", Branch: "feature"}
	srv, calls := setupMock(t, mockResponse{Status: 201, Body: result})
	c := &CommitsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	baseBranch := "main"
	_, err := c.Create(&CreateCommitOptions{
		Branch:      "feature",
		Message:     "branch from main",
		AuthorName:  "Alice",
		AuthorEmail: "alice@example.com",
		BaseBranch:  &baseBranch,
	}).
		Send(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}

	bb, ok := body["baseBranch"]
	if !ok {
		t.Fatal("expected baseBranch field in body")
	}
	if bb != "main" {
		t.Errorf("expected baseBranch=main, got %v", bb)
	}
}
