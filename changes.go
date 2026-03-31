package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
)

// ChangeInfo represents a mutable change in a repository.
type ChangeInfo struct {
	ID              string  `json:"id"`
	ChangeID        string  `json:"changeId"`
	RepoID          string  `json:"repoId"`
	OwnerID         string  `json:"ownerId"`
	ParentChangeID  *string `json:"parentChangeId,omitempty"`
	CommitSHA       *string `json:"commitSha,omitempty"`
	TreeSHA         *string `json:"treeSha,omitempty"`
	BaseCommitSHA   string  `json:"baseCommitSha"`
	Description     *string `json:"description,omitempty"`
	Status          string  `json:"status"`
	ConflictDetails *string `json:"conflictDetails,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// CreateChangeOptions are options for creating a change.
type CreateChangeOptions struct {
	BaseRef        string            `json:"baseRef,omitempty"`
	Description    *string           `json:"description,omitempty"`
	ParentChangeID *string           `json:"parentChangeId,omitempty"`
	Files          []ChangeFileEntry `json:"files,omitempty"`
}

// ChangeFileEntry represents a file to include in a change.
type ChangeFileEntry struct {
	Path     string  `json:"path"`
	Content  string  `json:"content"`
	Encoding *string `json:"encoding,omitempty"`
}

// AmendOptions are options for amending a change.
type AmendOptions struct {
	Files   []ChangeFileEntry `json:"files,omitempty"`
	Deletes []string          `json:"deletes,omitempty"`
}

// SquashOptions are options for squashing a change into its parent.
type SquashOptions struct {
	Files []string `json:"files,omitempty"`
}

// SplitOptions are options for splitting a change.
type SplitOptions struct {
	Files []string `json:"files"`
}

// SplitResult is the result of splitting a change.
type SplitResult struct {
	First     ChangeInfo `json:"first"`
	Remainder ChangeInfo `json:"remainder"`
}

// SquashResult is the result of squashing a change into its parent.
type SquashResult struct {
	Parent ChangeInfo `json:"parent"`
	Child  ChangeInfo `json:"child"`
}

// ChangeMaterializeResult is the result of materializing a change to a branch.
type ChangeMaterializeResult struct {
	Branch string `json:"branch"`
	SHA    string `json:"sha"`
}

// ChangesService provides change operations.
type ChangesService struct {
	client *httpClient
}

func (s *ChangesService) Create(ctx context.Context, repoID string, opts CreateChangeOptions) (*ChangeInfo, error) {
	raw, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/changes", repoID), opts)
	if err != nil {
		return nil, err
	}
	var result ChangeInfo
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal change: %w", err)
	}
	return &result, nil
}

func (s *ChangesService) List(ctx context.Context, repoID string) ([]ChangeInfo, error) {
	raw, err := s.client.get(ctx, fmt.Sprintf("/repos/%s/changes", repoID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Items []ChangeInfo `json:"items"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal changes: %w", err)
	}
	return result.Items, nil
}

func (s *ChangesService) Get(ctx context.Context, repoID, changeID string) (*ChangeInfo, error) {
	raw, err := s.client.get(ctx, fmt.Sprintf("/repos/%s/changes/%s", repoID, changeID), nil)
	if err != nil {
		return nil, err
	}
	var result ChangeInfo
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal change: %w", err)
	}
	return &result, nil
}

func (s *ChangesService) Abandon(ctx context.Context, repoID, changeID string) (*ChangeInfo, error) {
	if err := s.client.del(ctx, fmt.Sprintf("/repos/%s/changes/%s", repoID, changeID), nil); err != nil {
		return nil, err
	}
	return nil, nil
}

func (s *ChangesService) Amend(ctx context.Context, repoID, changeID string, opts AmendOptions) (*ChangeInfo, error) {
	raw, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/changes/%s/amend", repoID, changeID), opts)
	if err != nil {
		return nil, err
	}
	var result ChangeInfo
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal change: %w", err)
	}
	return &result, nil
}

func (s *ChangesService) Squash(ctx context.Context, repoID, changeID string, opts *SquashOptions) (*SquashResult, error) {
	body := opts
	if body == nil {
		body = &SquashOptions{}
	}
	raw, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/changes/%s/squash", repoID, changeID), body)
	if err != nil {
		return nil, err
	}
	var result SquashResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal squash result: %w", err)
	}
	return &result, nil
}

func (s *ChangesService) Split(ctx context.Context, repoID, changeID string, opts SplitOptions) (*SplitResult, error) {
	raw, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/changes/%s/split", repoID, changeID), opts)
	if err != nil {
		return nil, err
	}
	var result SplitResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal split result: %w", err)
	}
	return &result, nil
}

func (s *ChangesService) Materialize(ctx context.Context, repoID, changeID, branch string) (*ChangeMaterializeResult, error) {
	raw, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/changes/%s/materialize", repoID, changeID), map[string]string{"branch": branch})
	if err != nil {
		return nil, err
	}
	var result ChangeMaterializeResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal materialize result: %w", err)
	}
	return &result, nil
}
