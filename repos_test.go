package gitforge

import (
	"context"
	"encoding/json"
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
