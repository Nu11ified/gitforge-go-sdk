package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
)

type CreateTokenOptions struct {
	TTLSeconds int      `json:"ttlSeconds"`
	Scopes     []string `json:"scopes,omitempty"`
	Type       *string  `json:"type,omitempty"`
}

type TokensResource struct {
	client *httpClient
	repoID string
}

func (t *TokensResource) Create(ctx context.Context, opts *CreateTokenOptions) (*RepoToken, error) {
	raw, err := t.client.post(ctx, fmt.Sprintf("/repos/%s/tokens", t.repoID), opts)
	if err != nil {
		return nil, err
	}
	var token RepoToken
	if err := json.Unmarshal(raw, &token); err != nil {
		return nil, fmt.Errorf("unmarshal token: %w", err)
	}
	return &token, nil
}
