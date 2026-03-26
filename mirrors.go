package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
)

type CreateMirrorOptions struct {
	SourceURL    string  `json:"sourceUrl"`
	Direction    *string `json:"direction,omitempty"`
	Interval     *int    `json:"interval,omitempty"`
	Provider     *string `json:"provider,omitempty"`
	CredentialID *string `json:"credentialId,omitempty"`
}

type UpdateMirrorOptions struct {
	SourceURL    *string `json:"sourceUrl,omitempty"`
	Interval     *int    `json:"interval,omitempty"`
	Enabled      *bool   `json:"enabled,omitempty"`
	Direction    *string `json:"direction,omitempty"`
	Provider     *string `json:"provider,omitempty"`
	CredentialID *string `json:"credentialId,omitempty"`
}

type MirrorsResource struct {
	client *httpClient
	repoID string
}

func (m *MirrorsResource) List(ctx context.Context) ([]MirrorConfig, error) {
	raw, err := m.client.get(ctx, fmt.Sprintf("/repos/%s/mirrors", m.repoID), nil)
	if err != nil {
		return nil, err
	}
	var mirrors []MirrorConfig
	if err := json.Unmarshal(raw, &mirrors); err != nil {
		return nil, fmt.Errorf("unmarshal mirrors: %w", err)
	}
	return mirrors, nil
}

func (m *MirrorsResource) Create(ctx context.Context, opts *CreateMirrorOptions) (*MirrorConfig, error) {
	raw, err := m.client.post(ctx, fmt.Sprintf("/repos/%s/mirrors", m.repoID), opts)
	if err != nil {
		return nil, err
	}
	var mirror MirrorConfig
	if err := json.Unmarshal(raw, &mirror); err != nil {
		return nil, fmt.Errorf("unmarshal mirror: %w", err)
	}
	return &mirror, nil
}

func (m *MirrorsResource) Update(ctx context.Context, mirrorID string, opts *UpdateMirrorOptions) (*MirrorConfig, error) {
	raw, err := m.client.patch(ctx, fmt.Sprintf("/repos/%s/mirrors/%s", m.repoID, mirrorID), opts)
	if err != nil {
		return nil, err
	}
	var mirror MirrorConfig
	if err := json.Unmarshal(raw, &mirror); err != nil {
		return nil, fmt.Errorf("unmarshal mirror: %w", err)
	}
	return &mirror, nil
}

func (m *MirrorsResource) Delete(ctx context.Context, mirrorID string) error {
	return m.client.del(ctx, fmt.Sprintf("/repos/%s/mirrors/%s", m.repoID, mirrorID), nil)
}

func (m *MirrorsResource) Sync(ctx context.Context, mirrorID string) (*SyncResult, error) {
	raw, err := m.client.post(ctx, fmt.Sprintf("/repos/%s/mirrors/%s/sync", m.repoID, mirrorID), nil)
	if err != nil {
		return nil, err
	}
	var result SyncResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal sync: %w", err)
	}
	return &result, nil
}
