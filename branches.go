package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type ListBranchesOptions struct {
	Limit     *int
	Offset    *int
	Namespace *string
}

type CreateBranchOptions struct {
	Name              string  `json:"name"`
	BaseBranch        *string `json:"baseBranch,omitempty"`
	SHA               *string `json:"sha,omitempty"`
	TargetIsEphemeral *bool   `json:"targetIsEphemeral,omitempty"`
	BaseIsEphemeral   *bool   `json:"baseIsEphemeral,omitempty"`
	TTLSeconds        *int    `json:"ttlSeconds,omitempty"`
}

type DeleteBranchOptions struct {
	Namespace *string
}

type PromoteBranchOptions struct {
	BaseBranch   string  `json:"baseBranch"`
	TargetBranch *string `json:"targetBranch,omitempty"`
}

type BranchesResource struct {
	client *httpClient
	repoID string
}

func (b *BranchesResource) List(ctx context.Context, opts *ListBranchesOptions) (*PaginatedResponse[Branch], error) {
	q := url.Values{}
	if opts != nil {
		if opts.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *opts.Limit))
		}
		if opts.Offset != nil {
			q.Set("offset", fmt.Sprintf("%d", *opts.Offset))
		}
		if opts.Namespace != nil {
			q.Set("namespace", *opts.Namespace)
		}
	}
	raw, err := b.client.get(ctx, fmt.Sprintf("/repos/%s/branches", b.repoID), q)
	if err != nil {
		return nil, err
	}
	var result PaginatedResponse[Branch]
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal branches: %w", err)
	}
	return &result, nil
}

func (b *BranchesResource) Create(ctx context.Context, opts *CreateBranchOptions) (*Branch, error) {
	raw, err := b.client.post(ctx, fmt.Sprintf("/repos/%s/branches", b.repoID), opts)
	if err != nil {
		return nil, err
	}
	var branch Branch
	if err := json.Unmarshal(raw, &branch); err != nil {
		return nil, fmt.Errorf("unmarshal branch: %w", err)
	}
	return &branch, nil
}

func (b *BranchesResource) Delete(ctx context.Context, name string, opts *DeleteBranchOptions) error {
	q := url.Values{}
	if opts != nil && opts.Namespace != nil {
		q.Set("namespace", *opts.Namespace)
	}
	return b.client.del(ctx, fmt.Sprintf("/repos/%s/branches/%s", b.repoID, url.PathEscape(name)), q)
}

func (b *BranchesResource) Promote(ctx context.Context, opts *PromoteBranchOptions) (*PromoteResult, error) {
	raw, err := b.client.post(ctx, fmt.Sprintf("/repos/%s/branches/promote", b.repoID), opts)
	if err != nil {
		return nil, err
	}
	var result PromoteResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal promote: %w", err)
	}
	return &result, nil
}
