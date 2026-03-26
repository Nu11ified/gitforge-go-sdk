package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type CreateWebhookOptions struct {
	URL    string   `json:"url"`
	Secret *string  `json:"secret,omitempty"`
	Events []string `json:"events,omitempty"`
}

type ListWebhooksOptions struct {
	Limit  *int
	Offset *int
}

type WebhooksResource struct {
	client *httpClient
	repoID string
}

func (w *WebhooksResource) Create(ctx context.Context, opts *CreateWebhookOptions) (*Webhook, error) {
	raw, err := w.client.post(ctx, fmt.Sprintf("/repos/%s/webhooks", w.repoID), opts)
	if err != nil {
		return nil, err
	}
	var webhook Webhook
	if err := json.Unmarshal(raw, &webhook); err != nil {
		return nil, fmt.Errorf("unmarshal webhook: %w", err)
	}
	return &webhook, nil
}

func (w *WebhooksResource) List(ctx context.Context, opts *ListWebhooksOptions) (*PaginatedResponse[Webhook], error) {
	q := url.Values{}
	if opts != nil {
		if opts.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *opts.Limit))
		}
		if opts.Offset != nil {
			q.Set("offset", fmt.Sprintf("%d", *opts.Offset))
		}
	}
	raw, err := w.client.get(ctx, fmt.Sprintf("/repos/%s/webhooks", w.repoID), q)
	if err != nil {
		return nil, err
	}
	var result PaginatedResponse[Webhook]
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal webhooks: %w", err)
	}
	return &result, nil
}

func (w *WebhooksResource) Delete(ctx context.Context, webhookID string) error {
	return w.client.del(ctx, fmt.Sprintf("/repos/%s/webhooks/%s", w.repoID, webhookID), nil)
}

func (w *WebhooksResource) Test(ctx context.Context, webhookID string) (*WebhookTestResult, error) {
	raw, err := w.client.post(ctx, fmt.Sprintf("/repos/%s/webhooks/%s/test", w.repoID, webhookID), nil)
	if err != nil {
		return nil, err
	}
	var result WebhookTestResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal test result: %w", err)
	}
	return &result, nil
}

func (w *WebhooksResource) Deliveries(ctx context.Context, webhookID string, opts *ListWebhooksOptions) (*PaginatedResponse[WebhookDelivery], error) {
	q := url.Values{}
	if opts != nil {
		if opts.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *opts.Limit))
		}
		if opts.Offset != nil {
			q.Set("offset", fmt.Sprintf("%d", *opts.Offset))
		}
	}
	raw, err := w.client.get(ctx, fmt.Sprintf("/repos/%s/webhooks/%s/deliveries", w.repoID, webhookID), q)
	if err != nil {
		return nil, err
	}
	var result PaginatedResponse[WebhookDelivery]
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal deliveries: %w", err)
	}
	return &result, nil
}
