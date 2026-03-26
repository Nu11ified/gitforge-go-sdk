package gitforge

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"
)

func TestWebhooksCreate(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: Webhook{
			ID:     "wh-1",
			URL:    "https://example.com/hook",
			Events: []string{"push", "pr"},
			Active: true,
		},
	})

	w := &WebhooksResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	webhook, err := w.Create(context.Background(), &CreateWebhookOptions{
		URL:    "https://example.com/hook",
		Events: []string{"push", "pr"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.ID != "wh-1" {
		t.Errorf("expected id 'wh-1', got %q", webhook.ID)
	}
	if webhook.URL != "https://example.com/hook" {
		t.Errorf("expected url 'https://example.com/hook', got %q", webhook.URL)
	}
	if len(webhook.Events) != 2 {
		t.Errorf("expected 2 events, got %d", len(webhook.Events))
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/webhooks" {
		t.Errorf("expected path /repos/repo-123/webhooks, got %s", call.Path)
	}
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["url"] != "https://example.com/hook" {
		t.Errorf("expected url in body, got %v", body["url"])
	}
	if _, ok := body["active"]; ok {
		t.Errorf("active field should not be present in request body")
	}
}

func TestWebhooksCreateWithSecret(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: Webhook{
			ID:     "wh-2",
			URL:    "https://example.com/hook",
			Events: []string{"push"},
			Active: true,
		},
	})

	secret := "my-secret"
	w := &WebhooksResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	webhook, err := w.Create(context.Background(), &CreateWebhookOptions{
		URL:    "https://example.com/hook",
		Secret: &secret,
		Events: []string{"push"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if webhook.ID != "wh-2" {
		t.Errorf("expected id 'wh-2', got %q", webhook.ID)
	}

	call := (*calls)[0]
	var body map[string]any
	if err := json.Unmarshal(call.Body, &body); err != nil {
		t.Fatalf("failed to parse request body: %v", err)
	}
	if body["secret"] != "my-secret" {
		t.Errorf("expected secret 'my-secret' in body, got %v", body["secret"])
	}
}

func TestWebhooksList(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PaginatedResponse[Webhook]{
			Data: []Webhook{
				{ID: "wh-1", URL: "https://example.com/hook1", Events: []string{"push"}, Active: true},
				{ID: "wh-2", URL: "https://example.com/hook2", Events: []string{"pr"}, Active: false},
			},
			Total:   2,
			Limit:   10,
			Offset:  0,
			HasMore: false,
		},
	})

	w := &WebhooksResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := w.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 2 {
		t.Fatalf("expected 2 webhooks, got %d", len(result.Data))
	}
	if result.Data[0].ID != "wh-1" {
		t.Errorf("expected first webhook id 'wh-1', got %q", result.Data[0].ID)
	}
	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/webhooks" {
		t.Errorf("expected path /repos/repo-123/webhooks, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestWebhooksListWithOpts(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PaginatedResponse[Webhook]{
			Data:    []Webhook{{ID: "wh-3", URL: "https://example.com/hook3", Events: []string{"push"}, Active: true}},
			Total:   10,
			Limit:   5,
			Offset:  5,
			HasMore: false,
		},
	})

	w := &WebhooksResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := w.List(context.Background(), &ListWebhooksOptions{
		Limit:  ptrInt(5),
		Offset: ptrInt(5),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 webhook, got %d", len(result.Data))
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
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

func TestWebhooksDelete(t *testing.T) {
	srv, calls := setupMock(t, mockResponse{Status: 204})

	w := &WebhooksResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	err := w.Delete(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	call := (*calls)[0]
	if call.Method != "DELETE" {
		t.Errorf("expected DELETE, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/webhooks/wh-1" {
		t.Errorf("expected path /repos/repo-123/webhooks/wh-1, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}

func TestWebhooksTest(t *testing.T) {
	status := 200
	body := "ok"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: WebhookTestResult{
			Success:      true,
			Status:       &status,
			ResponseBody: &body,
			DurationMs:   42,
		},
	})

	w := &WebhooksResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := w.Test(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected success=true, got false")
	}
	if result.Status == nil || *result.Status != 200 {
		t.Errorf("expected status=200, got %v", result.Status)
	}
	if result.DurationMs != 42 {
		t.Errorf("expected durationMs=42, got %d", result.DurationMs)
	}

	call := (*calls)[0]
	if call.Method != "POST" {
		t.Errorf("expected POST, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/webhooks/wh-1/test" {
		t.Errorf("expected path /repos/repo-123/webhooks/wh-1/test, got %s", call.Path)
	}
}

func TestWebhooksDeliveries(t *testing.T) {
	respStatus := 200
	respBody := `{"event":"push"}`
	deliveredAt := "2026-03-26T10:00:00Z"
	srv, calls := setupMock(t, mockResponse{
		Status: 200,
		Body: PaginatedResponse[WebhookDelivery]{
			Data: []WebhookDelivery{
				{
					ID:             "del-1",
					EventType:      "push",
					Payload:        `{"ref":"refs/heads/main"}`,
					ResponseStatus: &respStatus,
					ResponseBody:   &respBody,
					DeliveredAt:    &deliveredAt,
					CreatedAt:      "2026-03-26T10:00:00Z",
				},
			},
			Total:   1,
			Limit:   10,
			Offset:  0,
			HasMore: false,
		},
	})

	w := &WebhooksResource{client: newHTTPClient(srv.URL, "tok", nil), repoID: "repo-123"}
	result, err := w.Deliveries(context.Background(), "wh-1", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data) != 1 {
		t.Fatalf("expected 1 delivery, got %d", len(result.Data))
	}
	if result.Data[0].ID != "del-1" {
		t.Errorf("expected delivery id 'del-1', got %q", result.Data[0].ID)
	}
	if result.Data[0].EventType != "push" {
		t.Errorf("expected eventType 'push', got %q", result.Data[0].EventType)
	}
	if result.Total != 1 {
		t.Errorf("expected total 1, got %d", result.Total)
	}

	call := (*calls)[0]
	if call.Method != "GET" {
		t.Errorf("expected GET, got %s", call.Method)
	}
	if call.Path != "/repos/repo-123/webhooks/wh-1/deliveries" {
		t.Errorf("expected path /repos/repo-123/webhooks/wh-1/deliveries, got %s", call.Path)
	}
	if call.Query != "" {
		t.Errorf("expected empty query, got %s", call.Query)
	}
}
