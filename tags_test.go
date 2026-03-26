package gitforge

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

func TestTagsList(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PaginatedResponse[Tag]{
			Data:    []Tag{{Name: "v1.0.0", SHA: "abc123"}, {Name: "v2.0.0", SHA: "def456"}},
			Total:   2,
			Limit:   10,
			Offset:  0,
			HasMore: false,
		},
	})

	tg := &TagsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := tg.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(result.Data))
	}
	if result.Data[0].Name != "v1.0.0" {
		t.Errorf("expected first tag 'v1.0.0', got %q", result.Data[0].Name)
	}
	if result.Data[1].SHA != "def456" {
		t.Errorf("expected second tag SHA 'def456', got %q", result.Data[1].SHA)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/tags" {
		t.Errorf("expected path /repos/repo-123/tags, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestTagsListWithOpts(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PaginatedResponse[Tag]{
			Data:    []Tag{{Name: "v3.0.0", SHA: "ghi789"}},
			Total:   10,
			Limit:   5,
			Offset:  5,
			HasMore: false,
		},
	})

	tg := &TagsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	limit := 5
	offset := 5
	result, err := tg.List(context.Background(), &ListTagsOptions{Limit: &limit, Offset: &offset})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(result.Data))
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/tags" {
		t.Errorf("expected path /repos/repo-123/tags, got %s", call.Path)
	}
	if call.Query == "" {
		t.Fatal("expected query params, got empty string")
	}
	vals, err := url.ParseQuery(call.Query)
	if err != nil {
		t.Fatalf("failed to parse query: %v", err)
	}
	if vals.Get("limit") != "5" {
		t.Errorf("expected limit=5, got %q", vals.Get("limit"))
	}
	if vals.Get("offset") != "5" {
		t.Errorf("expected offset=5, got %q", vals.Get("offset"))
	}
}

func TestTagsCreate(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 201,
		Body:   Tag{Name: "v1.0.0", SHA: "abc123"},
	})

	tg := &TagsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	tag, err := tg.Create(context.Background(), &CreateTagOptions{
		Name: "v1.0.0",
		SHA:  "abc123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tag.Name != "v1.0.0" {
		t.Errorf("expected tag name 'v1.0.0', got %q", tag.Name)
	}
	if tag.SHA != "abc123" {
		t.Errorf("expected tag SHA 'abc123', got %q", tag.SHA)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/tags" {
		t.Errorf("expected path /repos/repo-123/tags, got %s", call.Path)
	}

	var body map[string]json.RawMessage
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to unmarshal body: %v", err)
	}
	var name string
	if err := json.Unmarshal(body["name"], &name); err != nil || name != "v1.0.0" {
		t.Errorf("expected body name 'v1.0.0', got %q", name)
	}
	var sha string
	if err := json.Unmarshal(body["sha"], &sha); err != nil || sha != "abc123" {
		t.Errorf("expected body sha 'abc123', got %q", sha)
	}
}

func TestTagsDelete(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 204,
	})

	tg := &TagsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	err := tg.Delete(context.Background(), "v1.0.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	if call.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/tags/v1.0.0" {
		t.Errorf("expected path /repos/repo-123/tags/v1.0.0, got %s", call.Path)
	}
}
