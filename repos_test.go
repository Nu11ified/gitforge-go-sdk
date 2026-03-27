package gitforge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReposCreate(t *testing.T) {
	desc := "a test repo"
	vis := "private"
	repoResp := Repo{
		ID:            "repo-1",
		Name:          "my-repo",
		Description:   &desc,
		Visibility:    vis,
		DefaultBranch: "main",
	}

	srv, calls := setupMock(t, mockResponse{Status: 201, Body: repoResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	opts := &CreateRepoOptions{
		Name:        "my-repo",
		Description: &desc,
		Visibility:  &vis,
	}
	got, err := r.Create(context.Background(), opts)
	if err != nil {
		t.Fatalf("Create: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "POST" {
		t.Errorf("method: got %q, want %q", call.Method, "POST")
	}
	if call.Path != "/repos" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos")
	}

	var sentBody CreateRepoOptions
	if err := json.Unmarshal(call.Body, &sentBody); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if sentBody.Name != opts.Name {
		t.Errorf("body.name: got %q, want %q", sentBody.Name, opts.Name)
	}
	if sentBody.Description == nil || *sentBody.Description != desc {
		t.Errorf("body.description: got %v, want %q", sentBody.Description, desc)
	}
	if sentBody.Visibility == nil || *sentBody.Visibility != vis {
		t.Errorf("body.visibility: got %v, want %q", sentBody.Visibility, vis)
	}

	if got.ID != repoResp.ID {
		t.Errorf("response ID: got %q, want %q", got.ID, repoResp.ID)
	}
	if got.Name != repoResp.Name {
		t.Errorf("response Name: got %q, want %q", got.Name, repoResp.Name)
	}
	if got.Visibility != repoResp.Visibility {
		t.Errorf("response Visibility: got %q, want %q", got.Visibility, repoResp.Visibility)
	}
}

func TestReposListNoOpts(t *testing.T) {
	listResp := PaginatedResponse[Repo]{
		Data: []Repo{
			{ID: "repo-1", Name: "alpha", Visibility: "public", DefaultBranch: "main"},
			{ID: "repo-2", Name: "beta", Visibility: "private", DefaultBranch: "main"},
		},
		Total:   2,
		Limit:   20,
		Offset:  0,
		HasMore: false,
	}

	srv, calls := setupMock(t, mockResponse{Status: 200, Body: listResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	got, err := r.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("List: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want %q", call.Method, "GET")
	}
	if call.Path != "/repos" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos")
	}
	if call.Query != "" {
		t.Errorf("query: got %q, want empty", call.Query)
	}

	if len(got.Data) != 2 {
		t.Fatalf("response data length: got %d, want 2", len(got.Data))
	}
	if got.Data[0].ID != "repo-1" {
		t.Errorf("response Data[0].ID: got %q, want %q", got.Data[0].ID, "repo-1")
	}
	if got.Total != 2 {
		t.Errorf("response Total: got %d, want 2", got.Total)
	}
}

func TestReposListWithOpts(t *testing.T) {
	listResp := PaginatedResponse[Repo]{
		Data:    []Repo{{ID: "repo-3", Name: "gamma", Visibility: "public", DefaultBranch: "main"}},
		Total:   11,
		Limit:   10,
		Offset:  5,
		HasMore: false,
	}

	srv, calls := setupMock(t, mockResponse{Status: 200, Body: listResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	limit := 10
	offset := 5
	opts := &ListReposOptions{Limit: &limit, Offset: &offset}
	got, err := r.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want %q", call.Method, "GET")
	}
	if call.Path != "/repos" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos")
	}
	if call.Query != "limit=10&offset=5" {
		t.Errorf("query: got %q, want %q", call.Query, "limit=10&offset=5")
	}

	if len(got.Data) != 1 {
		t.Fatalf("response data length: got %d, want 1", len(got.Data))
	}
	if got.Data[0].ID != "repo-3" {
		t.Errorf("response Data[0].ID: got %q, want %q", got.Data[0].ID, "repo-3")
	}
	if got.Limit != 10 {
		t.Errorf("response Limit: got %d, want 10", got.Limit)
	}
	if got.Offset != 5 {
		t.Errorf("response Offset: got %d, want 5", got.Offset)
	}
}

func TestReposGet(t *testing.T) {
	repoResp := Repo{
		ID:            "repo-abc",
		Name:          "my-repo",
		Visibility:    "public",
		DefaultBranch: "main",
	}

	srv, calls := setupMock(t, mockResponse{Status: 200, Body: repoResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	got, err := r.Get(context.Background(), "repo-abc")
	if err != nil {
		t.Fatalf("Get: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want %q", call.Method, "GET")
	}
	if call.Path != "/repos/repo-abc" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos/repo-abc")
	}
	if call.Query != "" {
		t.Errorf("query: got %q, want empty", call.Query)
	}

	if got.ID != repoResp.ID {
		t.Errorf("response ID: got %q, want %q", got.ID, repoResp.ID)
	}
	if got.Name != repoResp.Name {
		t.Errorf("response Name: got %q, want %q", got.Name, repoResp.Name)
	}
	if got.DefaultBranch != repoResp.DefaultBranch {
		t.Errorf("response DefaultBranch: got %q, want %q", got.DefaultBranch, repoResp.DefaultBranch)
	}
}

func TestReposUpdate(t *testing.T) {
	newName := "renamed-repo"
	newBranch := "develop"
	repoResp := Repo{
		ID:            "repo-abc",
		Name:          newName,
		Visibility:    "public",
		DefaultBranch: newBranch,
	}

	srv, calls := setupMock(t, mockResponse{Status: 200, Body: repoResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	opts := &UpdateRepoOptions{
		Name:          &newName,
		DefaultBranch: &newBranch,
	}
	got, err := r.Update(context.Background(), "repo-abc", opts)
	if err != nil {
		t.Fatalf("Update: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "PATCH" {
		t.Errorf("method: got %q, want %q", call.Method, "PATCH")
	}
	if call.Path != "/repos/repo-abc" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos/repo-abc")
	}

	var sentBody UpdateRepoOptions
	if err := json.Unmarshal(call.Body, &sentBody); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if sentBody.Name == nil || *sentBody.Name != newName {
		t.Errorf("body.name: got %v, want %q", sentBody.Name, newName)
	}
	if sentBody.DefaultBranch == nil || *sentBody.DefaultBranch != newBranch {
		t.Errorf("body.defaultBranch: got %v, want %q", sentBody.DefaultBranch, newBranch)
	}

	if got.ID != repoResp.ID {
		t.Errorf("response ID: got %q, want %q", got.ID, repoResp.ID)
	}
	if got.Name != newName {
		t.Errorf("response Name: got %q, want %q", got.Name, newName)
	}
	if got.DefaultBranch != newBranch {
		t.Errorf("response DefaultBranch: got %q, want %q", got.DefaultBranch, newBranch)
	}
}

func TestReposDelete(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204, Body: nil})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	err := r.Delete(context.Background(), "repo-abc")
	if err != nil {
		t.Fatalf("Delete: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "DELETE" {
		t.Errorf("method: got %q, want %q", call.Method, "DELETE")
	}
	if call.Path != "/repos/repo-abc" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos/repo-abc")
	}
}

func TestCreateNote(t *testing.T) {
	noteResp := NoteResponse{SHA: "abc123", RefSHA: "def456", Success: true}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: noteResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	resp, err := r.CreateNote(context.Background(), "repo-1", &CreateNoteOptions{
		SHA: "abc123", Note: "LGTM", Author: Identity{Name: "Jane", Email: "j@e.com"},
	})
	if err != nil {
		t.Fatalf("CreateNote: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "POST" {
		t.Errorf("method: got %q, want POST", call.Method)
	}
	if call.Path != "/repos/repo-1/notes" {
		t.Errorf("path: got %q, want /repos/repo-1/notes", call.Path)
	}

	var sentBody map[string]any
	if err := json.Unmarshal(call.Body, &sentBody); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if sentBody["action"] != "add" {
		t.Errorf("body.action: got %v, want add", sentBody["action"])
	}
	if sentBody["sha"] != "abc123" {
		t.Errorf("body.sha: got %v, want abc123", sentBody["sha"])
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.SHA != "abc123" {
		t.Errorf("response SHA: got %q, want abc123", resp.SHA)
	}
}

func TestAppendNote(t *testing.T) {
	noteResp := NoteResponse{SHA: "abc123", RefSHA: "def456", Success: true}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: noteResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	resp, err := r.AppendNote(context.Background(), "repo-1", "abc123", "more text", Identity{Name: "Bob", Email: "b@e.com"})
	if err != nil {
		t.Fatalf("AppendNote: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "POST" {
		t.Errorf("method: got %q, want POST", call.Method)
	}
	if call.Path != "/repos/repo-1/notes" {
		t.Errorf("path: got %q, want /repos/repo-1/notes", call.Path)
	}

	var sentBody map[string]any
	if err := json.Unmarshal(call.Body, &sentBody); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if sentBody["action"] != "append" {
		t.Errorf("body.action: got %v, want append", sentBody["action"])
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
}

func TestGetNote(t *testing.T) {
	noteResp := NoteResponse{SHA: "abc123", RefSHA: "def456", Note: "Test note"}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: noteResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	resp, err := r.GetNote(context.Background(), "repo-1", "abc123")
	if err != nil {
		t.Fatalf("GetNote: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want GET", call.Method)
	}
	if call.Path != "/repos/repo-1/notes/abc123" {
		t.Errorf("path: got %q, want /repos/repo-1/notes/abc123", call.Path)
	}

	if resp.Note != "Test note" {
		t.Errorf("response Note: got %q, want Test note", resp.Note)
	}
	if resp.SHA != "abc123" {
		t.Errorf("response SHA: got %q, want abc123", resp.SHA)
	}
}

func TestDeleteNote(t *testing.T) {
	noteResp := NoteResponse{SHA: "abc123", RefSHA: "def456", Success: true}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: noteResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	author := Identity{Name: "Jane", Email: "j@e.com"}
	resp, err := r.DeleteNote(context.Background(), "repo-1", "abc123", &author)
	if err != nil {
		t.Fatalf("DeleteNote: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "DELETE" {
		t.Errorf("method: got %q, want DELETE", call.Method)
	}
	if call.Path != "/repos/repo-1/notes/abc123" {
		t.Errorf("path: got %q, want /repos/repo-1/notes/abc123", call.Path)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
}

func TestRestoreCommit(t *testing.T) {
	restoreResp := RestoreCommitResponse{
		CommitSHA: "newcommit", TreeSHA: "newtree", TargetBranch: "main", Success: true,
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: restoreResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	resp, err := r.RestoreCommit(context.Background(), "repo-1", &RestoreCommitOptions{
		TargetBranch: "main", TargetCommitSHA: "oldcommit", Author: Identity{Name: "A", Email: "a@e.com"},
	})
	if err != nil {
		t.Fatalf("RestoreCommit: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "POST" {
		t.Errorf("method: got %q, want POST", call.Method)
	}
	if call.Path != "/repos/repo-1/restore-commit" {
		t.Errorf("path: got %q, want /repos/repo-1/restore-commit", call.Path)
	}

	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.CommitSHA != "newcommit" {
		t.Errorf("response CommitSHA: got %q, want newcommit", resp.CommitSHA)
	}
	if resp.TargetBranch != "main" {
		t.Errorf("response TargetBranch: got %q, want main", resp.TargetBranch)
	}
}

func TestListFilesWithMetadata(t *testing.T) {
	commitSHA := "abc123"
	metaResp := FilesMetadataResponse{
		Files: []FileMetadata{
			{Path: "README.md", Mode: "100644", Size: 42, LastCommitSHA: &commitSHA},
		},
		Commits: map[string]CommitInfo{
			"abc123": {Author: "Jane", Date: "2026-01-01", Message: "initial"},
		},
		Ref: "main",
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: metaResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	ref := "main"
	resp, err := r.ListFilesWithMetadata(context.Background(), "repo-1", &ref, nil)
	if err != nil {
		t.Fatalf("ListFilesWithMetadata: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want GET", call.Method)
	}
	if call.Path != "/repos/repo-1/files/metadata" {
		t.Errorf("path: got %q, want /repos/repo-1/files/metadata", call.Path)
	}
	if call.Query != "ref=main" {
		t.Errorf("query: got %q, want ref=main", call.Query)
	}

	if len(resp.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(resp.Files))
	}
	if resp.Files[0].Path != "README.md" {
		t.Errorf("file path: got %q, want README.md", resp.Files[0].Path)
	}
	if resp.Ref != "main" {
		t.Errorf("ref: got %q, want main", resp.Ref)
	}
}

func TestPullUpstream(t *testing.T) {
	pullResp := PullUpstreamResponse{
		Status: "fast_forward", OldSHA: "aaa", NewSHA: "bbb", Branch: "main", Success: true,
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: pullResp})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	branch := "main"
	resp, err := r.PullUpstream(context.Background(), "repo-1", &PullUpstreamOptions{Branch: &branch})
	if err != nil {
		t.Fatalf("PullUpstream: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "POST" {
		t.Errorf("method: got %q, want POST", call.Method)
	}
	if call.Path != "/repos/repo-1/pull-upstream" {
		t.Errorf("path: got %q, want /repos/repo-1/pull-upstream", call.Path)
	}

	var sentBody map[string]any
	if err := json.Unmarshal(call.Body, &sentBody); err != nil {
		t.Fatalf("unmarshal request body: %v", err)
	}
	if sentBody["branch"] != "main" {
		t.Errorf("body.branch: got %v, want main", sentBody["branch"])
	}

	if resp.Status != "fast_forward" {
		t.Errorf("response Status: got %q, want fast_forward", resp.Status)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
}

func TestDetachUpstream(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204, Body: nil})
	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	resp, err := r.DetachUpstream(context.Background(), "repo-1")
	if err != nil {
		t.Fatalf("DetachUpstream: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "DELETE" {
		t.Errorf("method: got %q, want DELETE", call.Method)
	}
	if call.Path != "/repos/repo-1/base" {
		t.Errorf("path: got %q, want /repos/repo-1/base", call.Path)
	}

	if resp.Message != "repository detached" {
		t.Errorf("response Message: got %q, want repository detached", resp.Message)
	}
}

func TestGetRawFile(t *testing.T) {
	// Use a custom httptest server to return raw bytes (not JSON-encoded).
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("method: got %q, want GET", r.Method)
		}
		if r.URL.Path != "/repos/repo-1/raw/main" {
			t.Errorf("path: got %q, want /repos/repo-1/raw/main", r.URL.Path)
		}
		if r.URL.Query().Get("path") != "README.md" {
			t.Errorf("query path: got %q, want README.md", r.URL.Query().Get("path"))
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte("hello world"))
	}))
	t.Cleanup(srv.Close)

	r := &ReposResource{client: newHTTPClient(srv.URL, "test-token", nil)}

	data, err := r.GetRawFile(context.Background(), "repo-1", "main", "README.md")
	if err != nil {
		t.Fatalf("GetRawFile: unexpected error: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("response data: got %q, want hello world", string(data))
	}
}
