package gitforge

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

func ptrStr(s string) *string { return &s }
func ptrBool(b bool) *bool    { return &b }

func TestBranchList(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PaginatedResponse[Branch]{
			Data:    []Branch{{Name: "main", SHA: "abc123"}, {Name: "dev", SHA: "def456"}},
			Total:   2,
			Limit:   10,
			Offset:  0,
			HasMore: false,
		},
	})

	b := &BranchesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := b.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(result.Data))
	}
	if result.Data[0].Name != "main" {
		t.Errorf("expected first branch 'main', got %q", result.Data[0].Name)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/branches" {
		t.Errorf("expected path /repos/repo-123/branches, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestBranchListWithNamespace(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PaginatedResponse[Branch]{
			Data:    []Branch{{Name: "ephemeral/abc", SHA: "aaa111"}},
			Total:   1,
			Limit:   50,
			Offset:  0,
			HasMore: false,
		},
	})

	b := &BranchesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := b.List(context.Background(), &ListBranchesOptions{
		Namespace: ptrStr("ephemeral"),
		Limit:     ptrInt(50),
		Offset:    ptrInt(0),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(result.Data))
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/branches" {
		t.Errorf("expected path /repos/repo-123/branches, got %s", call.Path)
	}
	if call.Query == "" {
		t.Fatal("expected query params, got empty string")
	}
	vals, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if vals.Get("namespace") != "ephemeral" {
		t.Errorf("expected namespace=ephemeral, got %q", vals.Get("namespace"))
	}
	if vals.Get("limit") != "50" {
		t.Errorf("expected limit=50, got %q", vals.Get("limit"))
	}
	if vals.Get("offset") != "0" {
		t.Errorf("expected offset=0, got %q", vals.Get("offset"))
	}
}

func TestBranchCreateWithBaseBranch(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body:   Branch{Name: "feature/new", SHA: "aabbcc"},
	})

	b := &BranchesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	branch, err := b.Create(context.Background(), &CreateBranchOptions{
		Name:       "feature/new",
		BaseBranch: ptrStr("main"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch.Name != "feature/new" {
		t.Errorf("expected name 'feature/new', got %q", branch.Name)
	}
	if branch.SHA != "aabbcc" {
		t.Errorf("expected sha 'aabbcc', got %q", branch.SHA)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/branches" {
		t.Errorf("expected path /repos/repo-123/branches, got %s", call.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["name"] != "feature/new" {
		t.Errorf("expected name 'feature/new' in body, got %v", body["name"])
	}
	if body["baseBranch"] != "main" {
		t.Errorf("expected baseBranch 'main' in body, got %v", body["baseBranch"])
	}
}

func TestBranchCreateWithAllOptions(t *testing.T) {
	expiresAt := "2026-04-01T00:00:00Z"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body:   Branch{Name: "ephemeral/xyz", SHA: "112233", ExpiresAt: &expiresAt},
	})

	b := &BranchesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	branch, err := b.Create(context.Background(), &CreateBranchOptions{
		Name:              "ephemeral/xyz",
		BaseBranch:        ptrStr("main"),
		TargetIsEphemeral: ptrBool(true),
		BaseIsEphemeral:   ptrBool(false),
		TTLSeconds:        ptrInt(3600),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if branch.Name != "ephemeral/xyz" {
		t.Errorf("expected name 'ephemeral/xyz', got %q", branch.Name)
	}
	if branch.ExpiresAt == nil || *branch.ExpiresAt != "2026-04-01T00:00:00Z" {
		t.Errorf("expected expiresAt '2026-04-01T00:00:00Z', got %v", branch.ExpiresAt)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["targetIsEphemeral"] != true {
		t.Errorf("expected targetIsEphemeral=true, got %v", body["targetIsEphemeral"])
	}
	if body["baseIsEphemeral"] != false {
		t.Errorf("expected baseIsEphemeral=false, got %v", body["baseIsEphemeral"])
	}
	// JSON numbers are float64
	if body["ttlSeconds"] != float64(3600) {
		t.Errorf("expected ttlSeconds=3600, got %v", body["ttlSeconds"])
	}
}

func TestBranchDelete(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204})

	b := &BranchesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	err := b.Delete(context.Background(), "feature/test", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	if call.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", call.Method)
	}
	// net/http decodes %2F in r.URL.Path; the path segment is still correctly isolated
	if call.Path != "/repos/repo-123/branches/feature/test" {
		t.Errorf("expected path /repos/repo-123/branches/feature/test, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestBranchDeleteWithNamespace(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204})

	b := &BranchesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	err := b.Delete(context.Background(), "feature/test", &DeleteBranchOptions{
		Namespace: ptrStr("ephemeral"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	if call.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/branches/feature/test" {
		t.Errorf("expected path /repos/repo-123/branches/feature/test, got %s", call.Path)
	}
	vals, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if vals.Get("namespace") != "ephemeral" {
		t.Errorf("expected namespace=ephemeral, got %q", vals.Get("namespace"))
	}
}

func TestBranchPromote(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PromoteResult{
			TargetBranch: "main",
			CommitSHA:    "deadbeef",
		},
	})

	b := &BranchesResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := b.Promote(context.Background(), &PromoteBranchOptions{
		BaseBranch:   "ephemeral/feature-x",
		TargetBranch: ptrStr("main"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.TargetBranch != "main" {
		t.Errorf("expected targetBranch 'main', got %q", result.TargetBranch)
	}
	if result.CommitSHA != "deadbeef" {
		t.Errorf("expected commitSha 'deadbeef', got %q", result.CommitSHA)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/branches/promote" {
		t.Errorf("expected path /repos/repo-123/branches/promote, got %s", call.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["baseBranch"] != "ephemeral/feature-x" {
		t.Errorf("expected baseBranch 'ephemeral/feature-x', got %v", body["baseBranch"])
	}
	if body["targetBranch"] != "main" {
		t.Errorf("expected targetBranch 'main', got %v", body["targetBranch"])
	}
}
