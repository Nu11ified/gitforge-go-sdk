package gitforge

import (
	"context"
	"encoding/json"
	"testing"
)

func TestMirrorsList(t *testing.T) {
	credID := "cred-1"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: []MirrorConfig{
			{
				ID:        "mirror-1",
				SourceURL: "https://github.com/org/repo",
				Interval:  3600,
				Enabled:   true,
				CreatedAt: "2026-03-26T00:00:00Z",
				Direction: "pull",
				Provider:  "github",
			},
			{
				ID:           "mirror-2",
				SourceURL:    "https://gitlab.com/org/repo",
				Interval:     7200,
				Enabled:      false,
				CreatedAt:    "2026-03-25T00:00:00Z",
				Direction:    "push",
				Provider:     "gitlab",
				CredentialID: &credID,
			},
		},
	})

	m := &MirrorsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	mirrors, err := m.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mirrors) != 2 {
		t.Fatalf("expected 2 mirrors, got %d", len(mirrors))
	}
	if mirrors[0].ID != "mirror-1" {
		t.Errorf("expected first mirror id 'mirror-1', got %q", mirrors[0].ID)
	}
	if mirrors[0].SourceURL != "https://github.com/org/repo" {
		t.Errorf("expected sourceUrl 'https://github.com/org/repo', got %q", mirrors[0].SourceURL)
	}
	if mirrors[0].Direction != "pull" {
		t.Errorf("expected direction 'pull', got %q", mirrors[0].Direction)
	}
	if mirrors[1].ID != "mirror-2" {
		t.Errorf("expected second mirror id 'mirror-2', got %q", mirrors[1].ID)
	}
	if mirrors[1].CredentialID == nil || *mirrors[1].CredentialID != credID {
		t.Errorf("expected credentialId %q, got %v", credID, mirrors[1].CredentialID)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/mirrors" {
		t.Errorf("expected path /repos/repo-123/mirrors, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestMirrorsCreateRequired(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: MirrorConfig{
			ID:        "mirror-abc",
			SourceURL: "https://github.com/org/repo",
			Interval:  3600,
			Enabled:   true,
			CreatedAt: "2026-03-26T00:00:00Z",
			Direction: "pull",
			Provider:  "github",
		},
	})

	m := &MirrorsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	mirror, err := m.Create(context.Background(), &CreateMirrorOptions{
		SourceURL: "https://github.com/org/repo",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mirror.ID != "mirror-abc" {
		t.Errorf("expected id 'mirror-abc', got %q", mirror.ID)
	}
	if mirror.SourceURL != "https://github.com/org/repo" {
		t.Errorf("expected sourceUrl 'https://github.com/org/repo', got %q", mirror.SourceURL)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/mirrors" {
		t.Errorf("expected path /repos/repo-123/mirrors, got %s", call.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["sourceUrl"] != "https://github.com/org/repo" {
		t.Errorf("expected sourceUrl in body, got %v", body["sourceUrl"])
	}
}

func TestMirrorsCreateAllOptions(t *testing.T) {
	credID := "cred-xyz"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: MirrorConfig{
			ID:           "mirror-xyz",
			SourceURL:    "https://github.com/org/repo",
			Interval:     1800,
			Enabled:      true,
			CreatedAt:    "2026-03-26T00:00:00Z",
			Direction:    "push",
			Provider:     "github",
			CredentialID: &credID,
		},
	})

	m := &MirrorsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	mirror, err := m.Create(context.Background(), &CreateMirrorOptions{
		SourceURL:    "https://github.com/org/repo",
		Direction:    ptrStr("push"),
		Interval:     ptrInt(1800),
		Provider:     ptrStr("github"),
		CredentialID: ptrStr(credID),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mirror.ID != "mirror-xyz" {
		t.Errorf("expected id 'mirror-xyz', got %q", mirror.ID)
	}
	if mirror.Direction != "push" {
		t.Errorf("expected direction 'push', got %q", mirror.Direction)
	}
	if mirror.Interval != 1800 {
		t.Errorf("expected interval 1800, got %d", mirror.Interval)
	}
	if mirror.CredentialID == nil || *mirror.CredentialID != credID {
		t.Errorf("expected credentialId %q, got %v", credID, mirror.CredentialID)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["direction"] != "push" {
		t.Errorf("expected direction 'push' in body, got %v", body["direction"])
	}
	if body["interval"] != float64(1800) {
		t.Errorf("expected interval 1800 in body, got %v", body["interval"])
	}
	if body["provider"] != "github" {
		t.Errorf("expected provider 'github' in body, got %v", body["provider"])
	}
	if body["credentialId"] != credID {
		t.Errorf("expected credentialId %q in body, got %v", credID, body["credentialId"])
	}
}

func TestMirrorsUpdate(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: MirrorConfig{
			ID:        "mirror-abc",
			SourceURL: "https://github.com/org/repo",
			Interval:  900,
			Enabled:   false,
			CreatedAt: "2026-03-26T00:00:00Z",
			Direction: "pull",
			Provider:  "github",
		},
	})

	m := &MirrorsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	mirror, err := m.Update(context.Background(), "mirror-abc", &UpdateMirrorOptions{
		Interval: ptrInt(900),
		Enabled:  ptrBool(false),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mirror.ID != "mirror-abc" {
		t.Errorf("expected id 'mirror-abc', got %q", mirror.ID)
	}
	if mirror.Interval != 900 {
		t.Errorf("expected interval 900, got %d", mirror.Interval)
	}
	if mirror.Enabled {
		t.Errorf("expected enabled false, got true")
	}

	call := (*calls)[0]
	if call.Method != "PATCH" {
		t.Errorf("expected PATCH, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/mirrors/mirror-abc" {
		t.Errorf("expected path /repos/repo-123/mirrors/mirror-abc, got %s", call.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["interval"] != float64(900) {
		t.Errorf("expected interval 900 in body, got %v", body["interval"])
	}
	if body["enabled"] != false {
		t.Errorf("expected enabled false in body, got %v", body["enabled"])
	}
}

func TestMirrorsDelete(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204})

	m := &MirrorsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	err := m.Delete(context.Background(), "mirror-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	if call.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/mirrors/mirror-abc" {
		t.Errorf("expected path /repos/repo-123/mirrors/mirror-abc, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestMirrorsSync(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: SyncResult{
			Message: "sync triggered",
		},
	})

	m := &MirrorsResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := m.Sync(context.Background(), "mirror-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Message != "sync triggered" {
		t.Errorf("expected message 'sync triggered', got %q", result.Message)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/mirrors/mirror-abc/sync" {
		t.Errorf("expected path /repos/repo-123/mirrors/mirror-abc/sync, got %s", call.Path)
	}
}
