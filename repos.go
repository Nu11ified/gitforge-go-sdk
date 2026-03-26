package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// CreateRepoOptions are options for creating a repository.
type CreateRepoOptions struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
}

// ListReposOptions are options for listing repositories.
type ListReposOptions struct {
	Limit  *int
	Offset *int
}

// UpdateRepoOptions are options for updating a repository.
type UpdateRepoOptions struct {
	Name                *string `json:"name,omitempty"`
	Description         *string `json:"description,omitempty"`
	DefaultBranch       *string `json:"defaultBranch,omitempty"`
	MergeCommitTemplate *string `json:"mergeCommitTemplate,omitempty"`
}

// ReposResource provides repository operations.
type ReposResource struct {
	client *httpClient
}

func (r *ReposResource) Create(ctx context.Context, opts *CreateRepoOptions) (*Repo, error) {
	raw, err := r.client.post(ctx, "/repos", opts)
	if err != nil {
		return nil, err
	}
	var repo Repo
	if err := json.Unmarshal(raw, &repo); err != nil {
		return nil, fmt.Errorf("unmarshal repo: %w", err)
	}
	return &repo, nil
}

func (r *ReposResource) List(ctx context.Context, opts *ListReposOptions) (*PaginatedResponse[Repo], error) {
	q := url.Values{}
	if opts != nil {
		if opts.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *opts.Limit))
		}
		if opts.Offset != nil {
			q.Set("offset", fmt.Sprintf("%d", *opts.Offset))
		}
	}
	raw, err := r.client.get(ctx, "/repos", q)
	if err != nil {
		return nil, err
	}
	var result PaginatedResponse[Repo]
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal repos: %w", err)
	}
	return &result, nil
}

func (r *ReposResource) Get(ctx context.Context, id string) (*Repo, error) {
	raw, err := r.client.get(ctx, "/repos/"+id, nil)
	if err != nil {
		return nil, err
	}
	var repo Repo
	if err := json.Unmarshal(raw, &repo); err != nil {
		return nil, fmt.Errorf("unmarshal repo: %w", err)
	}
	return &repo, nil
}

func (r *ReposResource) Update(ctx context.Context, id string, opts *UpdateRepoOptions) (*Repo, error) {
	raw, err := r.client.patch(ctx, "/repos/"+id, opts)
	if err != nil {
		return nil, err
	}
	var repo Repo
	if err := json.Unmarshal(raw, &repo); err != nil {
		return nil, fmt.Errorf("unmarshal repo: %w", err)
	}
	return &repo, nil
}

func (r *ReposResource) Delete(ctx context.Context, id string) error {
	return r.client.del(ctx, "/repos/"+id, nil)
}
