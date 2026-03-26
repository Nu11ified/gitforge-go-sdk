package gitforge

import (
	"context"
	"testing"
)

func TestSearchCodeBasic(t *testing.T) {
	resp := SearchCodeResult{
		Results: []SearchResult{
			{
				RepoID:   "repo-123",
				RepoName: "my-repo",
				FilePath: "main.go",
				Branch:   "main",
				Matches:  []SearchMatch{{Line: 10, Content: "func main()", Highlight: "main"}},
			},
		},
		Total:   1,
		Page:    1,
		PerPage: 20,
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: resp})
	s := &SearchResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	got, err := s.SearchCode(context.Background(), &SearchCodeOptions{Query: "func main"})
	if err != nil {
		t.Fatalf("SearchCode: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want %q", call.Method, "GET")
	}
	if call.Path != "/repos/repo-123/search" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos/repo-123/search")
	}
	if call.Query != "q=func+main" {
		t.Errorf("query: got %q, want %q", call.Query, "q=func+main")
	}

	if got.Total != 1 {
		t.Errorf("Total: got %d, want 1", got.Total)
	}
	if len(got.Results) != 1 {
		t.Fatalf("Results length: got %d, want 1", len(got.Results))
	}
	if got.Results[0].FilePath != "main.go" {
		t.Errorf("Results[0].FilePath: got %q, want %q", got.Results[0].FilePath, "main.go")
	}
}

func TestSearchCodeAllOptions(t *testing.T) {
	lang := "go"
	branch := "dev"
	resp := SearchCodeResult{
		Results: []SearchResult{},
		Total:   0,
		Page:    2,
		PerPage: 10,
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: resp})
	s := &SearchResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	opts := &SearchCodeOptions{
		Query:    "hello",
		Language: &lang,
		Branch:   &branch,
		PerPage:  ptrInt(10),
		Page:     ptrInt(2),
	}
	_, err := s.SearchCode(context.Background(), opts)
	if err != nil {
		t.Fatalf("SearchCode: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want %q", call.Method, "GET")
	}
	if call.Path != "/repos/repo-123/search" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos/repo-123/search")
	}

	// Verify all query params are present
	for _, param := range []string{"q=hello", "lang=go", "branch=dev", "perPage=10", "page=2"} {
		found := false
		for _, part := range splitQuery(call.Query) {
			if part == param {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("query param %q not found in %q", param, call.Query)
		}
	}
}

func TestSearchCompare(t *testing.T) {
	resp := Comparison{
		Ahead:  3,
		Behind: 1,
		Commits: []CommitSummary{
			{SHA: "aaa111", Message: "add feature", Author: "Alice", Date: "2026-01-01T00:00:00Z"},
		},
		Files: []FileChange{
			{Path: "foo.go", Status: "modified"},
		},
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: resp})
	s := &SearchResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	got, err := s.Compare(context.Background(), "main", "feature-branch")
	if err != nil {
		t.Fatalf("Compare: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want %q", call.Method, "GET")
	}
	if call.Path != "/repos/repo-123/compare" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos/repo-123/compare")
	}
	for _, param := range []string{"base=main", "head=feature-branch"} {
		found := false
		for _, part := range splitQuery(call.Query) {
			if part == param {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("query param %q not found in %q", param, call.Query)
		}
	}

	if got.Ahead != 3 {
		t.Errorf("Ahead: got %d, want 3", got.Ahead)
	}
	if got.Behind != 1 {
		t.Errorf("Behind: got %d, want 1", got.Behind)
	}
	if len(got.Commits) != 1 {
		t.Fatalf("Commits length: got %d, want 1", len(got.Commits))
	}
	if got.Commits[0].SHA != "aaa111" {
		t.Errorf("Commits[0].SHA: got %q, want %q", got.Commits[0].SHA, "aaa111")
	}
}

func TestSearchCompareDiff(t *testing.T) {
	resp := []DiffEntry{
		{Path: "foo.go", Status: "modified", Additions: 5, Deletions: 2, Patch: "@@ -1,2 +1,5 @@"},
		{Path: "bar.go", Status: "added", Additions: 10, Deletions: 0, Patch: "@@ -0,0 +1,10 @@"},
	}
	srv, calls := setupMock(t, mockResponse{Status: 200, Body: resp})
	s := &SearchResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}

	got, err := s.CompareDiff(context.Background(), "main", "feature-branch")
	if err != nil {
		t.Fatalf("CompareDiff: unexpected error: %v", err)
	}

	if len(*calls) != 1 {
		t.Fatalf("expected 1 call, got %d", len(*calls))
	}
	call := (*calls)[0]

	if call.Method != "GET" {
		t.Errorf("method: got %q, want %q", call.Method, "GET")
	}
	if call.Path != "/repos/repo-123/compare/diff" {
		t.Errorf("path: got %q, want %q", call.Path, "/repos/repo-123/compare/diff")
	}
	for _, param := range []string{"base=main", "head=feature-branch"} {
		found := false
		for _, part := range splitQuery(call.Query) {
			if part == param {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("query param %q not found in %q", param, call.Query)
		}
	}

	if len(got) != 2 {
		t.Fatalf("DiffEntries length: got %d, want 2", len(got))
	}
	if got[0].Path != "foo.go" {
		t.Errorf("DiffEntries[0].Path: got %q, want %q", got[0].Path, "foo.go")
	}
	if got[0].Additions != 5 {
		t.Errorf("DiffEntries[0].Additions: got %d, want 5", got[0].Additions)
	}
	if got[1].Status != "added" {
		t.Errorf("DiffEntries[1].Status: got %q, want %q", got[1].Status, "added")
	}
}

// splitQuery splits a raw query string into individual key=value pairs.
func splitQuery(q string) []string {
	if q == "" {
		return nil
	}
	var parts []string
	start := 0
	for i := 0; i < len(q); i++ {
		if q[i] == '&' {
			parts = append(parts, q[start:i])
			start = i + 1
		}
	}
	parts = append(parts, q[start:])
	return parts
}
