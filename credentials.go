package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
)

type CreateCredentialOptions struct {
	Provider  string  `json:"provider"`
	Token     string  `json:"token"`
	Username  *string `json:"username,omitempty"`
	Label     *string `json:"label,omitempty"`
	SourceURL *string `json:"sourceUrl,omitempty"`
}

type UpdateCredentialOptions struct {
	Token    *string `json:"token,omitempty"`
	Username *string `json:"username,omitempty"`
	Label    *string `json:"label,omitempty"`
}

type CredentialsResource struct {
	client *httpClient
	repoID string
}

func (c *CredentialsResource) Create(ctx context.Context, opts *CreateCredentialOptions) (*GitCredential, error) {
	raw, err := c.client.post(ctx, fmt.Sprintf("/repos/%s/credentials", c.repoID), opts)
	if err != nil {
		return nil, err
	}
	var cred GitCredential
	if err := json.Unmarshal(raw, &cred); err != nil {
		return nil, fmt.Errorf("unmarshal credential: %w", err)
	}
	return &cred, nil
}

func (c *CredentialsResource) List(ctx context.Context) ([]GitCredential, error) {
	raw, err := c.client.get(ctx, fmt.Sprintf("/repos/%s/credentials", c.repoID), nil)
	if err != nil {
		return nil, err
	}
	var creds []GitCredential
	if err := json.Unmarshal(raw, &creds); err != nil {
		return nil, fmt.Errorf("unmarshal credentials: %w", err)
	}
	return creds, nil
}

func (c *CredentialsResource) Update(ctx context.Context, credID string, opts *UpdateCredentialOptions) (*GitCredential, error) {
	raw, err := c.client.patch(ctx, fmt.Sprintf("/repos/%s/credentials/%s", c.repoID, credID), opts)
	if err != nil {
		return nil, err
	}
	var cred GitCredential
	if err := json.Unmarshal(raw, &cred); err != nil {
		return nil, fmt.Errorf("unmarshal credential: %w", err)
	}
	return &cred, nil
}

func (c *CredentialsResource) Delete(ctx context.Context, credID string) error {
	return c.client.del(ctx, fmt.Sprintf("/repos/%s/credentials/%s", c.repoID, credID), nil)
}
