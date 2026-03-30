package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

// PatchSetInfo represents a patch set.
type PatchSetInfo struct {
	ID                 string  `json:"id"`
	RepoID             string  `json:"repoId"`
	Name               string  `json:"name"`
	Description        *string `json:"description,omitempty"`
	BaseRef            string  `json:"baseRef"`
	BaseSHA            string  `json:"baseSha"`
	MaterializedBranch *string `json:"materializedBranch,omitempty"`
	MaterializedSHA    *string `json:"materializedSha,omitempty"`
	Status             string  `json:"status"`
	AutoRebase         bool    `json:"autoRebase"`
	Visibility         string  `json:"visibility"`
	CreatedAt          string  `json:"createdAt"`
	UpdatedAt          string  `json:"updatedAt"`
}

// PatchInfo represents a single patch in a patch set.
type PatchInfo struct {
	ID              string  `json:"id"`
	Name            string  `json:"name"`
	Description     *string `json:"description,omitempty"`
	Order           int     `json:"order"`
	Diff            string  `json:"diff"`
	Status          string  `json:"status"`
	ConflictDetails *string `json:"conflictDetails,omitempty"`
	AuthorName      *string `json:"authorName,omitempty"`
	AuthorEmail     *string `json:"authorEmail,omitempty"`
	CreatedAt       string  `json:"createdAt"`
	UpdatedAt       string  `json:"updatedAt"`
}

// PatchSetWithPatches is a patch set with its patches.
type PatchSetWithPatches struct {
	PatchSetInfo
	Patches []PatchInfo `json:"patches"`
}

// CreatePatchSetOptions are options for creating a patch set.
type CreatePatchSetOptions struct {
	RepoID      string  `json:"repoId"`
	Name        string  `json:"name"`
	BaseRef     *string `json:"baseRef,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdatePatchSetOptions are options for updating a patch set.
type UpdatePatchSetOptions struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	AutoRebase  *bool   `json:"autoRebase,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
}

// AddPatchOptions are options for adding a patch.
type AddPatchOptions struct {
	Name        string  `json:"name"`
	Diff        string  `json:"diff"`
	Description *string `json:"description,omitempty"`
	AuthorName  *string `json:"authorName,omitempty"`
	AuthorEmail *string `json:"authorEmail,omitempty"`
}

// UpdatePatchOptions are options for updating a patch.
type UpdatePatchOptions struct {
	Name   *string `json:"name,omitempty"`
	Status *string `json:"status,omitempty"`
	Order  *int    `json:"order,omitempty"`
}

// PatchSetCreateResult is the result of creating a patch set.
type PatchSetCreateResult struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	MaterializedBranch string `json:"materializedBranch"`
}

// AddPatchResult is the result of adding a patch.
type AddPatchResult struct {
	ID    string `json:"id"`
	Order int    `json:"order"`
}

// RebaseResult is the result of rebasing a patch set.
type RebaseResult struct {
	Status          string  `json:"status"`
	ConflictedPatch *string `json:"conflictedPatch,omitempty"`
}

// MaterializeResult is the result of materializing a patch set.
type MaterializeResult struct {
	HeadSHA string `json:"headSha"`
	Status  string `json:"status"`
}

// PatchSetsResource provides patch set operations.
type PatchSetsResource struct {
	client *httpClient
}

func (r *PatchSetsResource) Create(ctx context.Context, opts *CreatePatchSetOptions) (*PatchSetCreateResult, error) {
	raw, err := r.client.post(ctx, "/patch-sets", opts)
	if err != nil {
		return nil, err
	}
	var result PatchSetCreateResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal patch set create: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) List(ctx context.Context, repoID *string) ([]PatchSetInfo, error) {
	q := url.Values{}
	if repoID != nil {
		q.Set("repoId", *repoID)
	}
	raw, err := r.client.get(ctx, "/patch-sets", q)
	if err != nil {
		return nil, err
	}
	var sets []PatchSetInfo
	if err := json.Unmarshal(raw, &sets); err != nil {
		return nil, fmt.Errorf("unmarshal patch sets: %w", err)
	}
	return sets, nil
}

func (r *PatchSetsResource) Get(ctx context.Context, setID string) (*PatchSetWithPatches, error) {
	raw, err := r.client.get(ctx, "/patch-sets/"+setID, nil)
	if err != nil {
		return nil, err
	}
	var result PatchSetWithPatches
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal patch set: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) Update(ctx context.Context, setID string, opts *UpdatePatchSetOptions) (*PatchSetInfo, error) {
	raw, err := r.client.patch(ctx, "/patch-sets/"+setID, opts)
	if err != nil {
		return nil, err
	}
	var result PatchSetInfo
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal patch set: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) Delete(ctx context.Context, setID string) error {
	return r.client.del(ctx, "/patch-sets/"+setID, nil)
}

func (r *PatchSetsResource) AddPatch(ctx context.Context, setID string, opts *AddPatchOptions) (*AddPatchResult, error) {
	raw, err := r.client.post(ctx, "/patch-sets/"+setID+"/patches", opts)
	if err != nil {
		return nil, err
	}
	var result AddPatchResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal add patch: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) UpdatePatch(ctx context.Context, setID, patchID string, opts *UpdatePatchOptions) error {
	_, err := r.client.patch(ctx, fmt.Sprintf("/patch-sets/%s/patches/%s", setID, patchID), opts)
	return err
}

func (r *PatchSetsResource) RemovePatch(ctx context.Context, setID, patchID string) error {
	return r.client.del(ctx, fmt.Sprintf("/patch-sets/%s/patches/%s", setID, patchID), nil)
}

func (r *PatchSetsResource) Rebase(ctx context.Context, setID string) (*RebaseResult, error) {
	raw, err := r.client.post(ctx, "/patch-sets/"+setID+"/rebase", nil)
	if err != nil {
		return nil, err
	}
	var result RebaseResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal rebase: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) Materialize(ctx context.Context, setID string) (*MaterializeResult, error) {
	raw, err := r.client.post(ctx, "/patch-sets/"+setID+"/materialize", nil)
	if err != nil {
		return nil, err
	}
	var result MaterializeResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal materialize: %w", err)
	}
	return &result, nil
}

// ---- Sharing methods ----

// PublishResult is the result of publishing a patch set.
type PublishResult struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Visibility string `json:"visibility"`
	OwnerName  string `json:"ownerName"`
	RepoName   string `json:"repoName"`
	PatchCount int    `json:"patchCount"`
}

// ForkResult is the result of forking a patch set.
type ForkResult struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	ForkedFromID string `json:"forkedFromId"`
	RepoID       string `json:"repoId"`
	PatchCount   int    `json:"patchCount"`
}

// ForkOptions are options for forking a patch set.
type ForkOptions struct {
	Name *string `json:"name,omitempty"`
}

// UpdateChange represents a single upstream change.
type UpdateChange struct {
	Type    string `json:"type"`
	PatchID string `json:"patchId"`
	Name    string `json:"name"`
	Order   int    `json:"order"`
}

// UpdatesResult is the result of checking for upstream updates.
type UpdatesResult struct {
	HasUpdates    bool           `json:"hasUpdates"`
	UpstreamSetID string         `json:"upstreamSetId"`
	Changes       []UpdateChange `json:"changes"`
}

// AcceptUpdatesOptions are options for accepting upstream updates.
type AcceptUpdatesOptions struct {
	Patches []string `json:"patches"`
}

// AcceptResult is the result of accepting upstream updates.
type AcceptResult struct {
	Accepted  int `json:"accepted"`
	Conflicts int `json:"conflicts"`
}

// ExploreOptions are options for browsing public patch sets.
type ExploreOptions struct {
	Query        *string
	Base         *string
	Page         *int
	PageSize     *int
}

// ExploreResult is the result of browsing public patch sets.
type ExploreResult struct {
	Items    []PatchSetInfo `json:"items"`
	Total    int            `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
}

func (r *PatchSetsResource) Publish(ctx context.Context, setID string) (*PublishResult, error) {
	raw, err := r.client.post(ctx, "/patch-sets/"+setID+"/publish", nil)
	if err != nil {
		return nil, err
	}
	var result PublishResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal publish: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) Unpublish(ctx context.Context, setID string) error {
	return r.client.del(ctx, "/patch-sets/"+setID+"/publish", nil)
}

func (r *PatchSetsResource) Fork(ctx context.Context, setID string, opts *ForkOptions) (*ForkResult, error) {
	raw, err := r.client.post(ctx, "/patch-sets/"+setID+"/fork", opts)
	if err != nil {
		return nil, err
	}
	var result ForkResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal fork: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) Subscribe(ctx context.Context, setID string) error {
	_, err := r.client.post(ctx, "/patch-sets/"+setID+"/subscribe", nil)
	return err
}

func (r *PatchSetsResource) Unsubscribe(ctx context.Context, setID string) error {
	return r.client.del(ctx, "/patch-sets/"+setID+"/subscribe", nil)
}

func (r *PatchSetsResource) GetUpdates(ctx context.Context, setID string) (*UpdatesResult, error) {
	raw, err := r.client.get(ctx, "/patch-sets/"+setID+"/updates", nil)
	if err != nil {
		return nil, err
	}
	var result UpdatesResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal updates: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) AcceptUpdates(ctx context.Context, setID string, opts *AcceptUpdatesOptions) (*AcceptResult, error) {
	raw, err := r.client.post(ctx, "/patch-sets/"+setID+"/updates/accept", opts)
	if err != nil {
		return nil, err
	}
	var result AcceptResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal accept: %w", err)
	}
	return &result, nil
}

func (r *PatchSetsResource) Explore(ctx context.Context, opts *ExploreOptions) (*ExploreResult, error) {
	q := url.Values{}
	if opts != nil {
		if opts.Query != nil {
			q.Set("q", *opts.Query)
		}
		if opts.Base != nil {
			q.Set("base", *opts.Base)
		}
		if opts.Page != nil {
			q.Set("page", fmt.Sprintf("%d", *opts.Page))
		}
		if opts.PageSize != nil {
			q.Set("pageSize", fmt.Sprintf("%d", *opts.PageSize))
		}
	}
	raw, err := r.client.get(ctx, "/explore/patch-sets", q)
	if err != nil {
		return nil, err
	}
	var result ExploreResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal explore: %w", err)
	}
	return &result, nil
}
