package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
)

// CreateSandboxUrlOptions configures sandbox URL creation.
type CreateSandboxUrlOptions struct {
	TTLSeconds int      `json:"ttlSeconds"`
	Scopes     []string `json:"scopes,omitempty"`
	Branch     *string  `json:"branch,omitempty"`
	Ephemeral  *bool    `json:"ephemeral,omitempty"`
}

// SandboxResource provides sandbox URL operations scoped to a repository.
type SandboxResource struct {
	client *httpClient
	repoID string
}

// CreateSandboxUrl generates a time-limited, scoped Git remote URL for sandbox use.
func (s *SandboxResource) CreateSandboxUrl(ctx context.Context, opts *CreateSandboxUrlOptions) (*SandboxUrl, error) {
	raw, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/sandbox-url", s.repoID), opts)
	if err != nil {
		return nil, err
	}
	var result SandboxUrl
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal sandbox url: %w", err)
	}
	return &result, nil
}
