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

// Identity represents a git author/committer.
type Identity struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateNoteOptions are options for creating a git note.
type CreateNoteOptions struct {
	SHA            string   `json:"sha"`
	Note           string   `json:"note"`
	Author         Identity `json:"author"`
	ExpectedRefSHA *string  `json:"expectedRefSha,omitempty"`
}

// RestoreCommitOptions are options for restoring a commit.
type RestoreCommitOptions struct {
	TargetBranch    string    `json:"targetBranch"`
	TargetCommitSHA string    `json:"targetCommitSha"`
	Author          Identity  `json:"author"`
	Committer       *Identity `json:"committer,omitempty"`
	CommitMessage   *string   `json:"commitMessage,omitempty"`
	ExpectedHeadSHA *string   `json:"expectedHeadSha,omitempty"`
}

// PullUpstreamOptions are options for pulling from upstream.
type PullUpstreamOptions struct {
	Branch *string `json:"branch,omitempty"`
}

// NoteResponse is the response from note operations.
type NoteResponse struct {
	SHA     string `json:"sha"`
	RefSHA  string `json:"refSha"`
	Note    string `json:"note,omitempty"`
	Success bool   `json:"success,omitempty"`
}

// RestoreCommitResponse is the response from restore-commit.
type RestoreCommitResponse struct {
	CommitSHA    string `json:"commitSha"`
	TreeSHA      string `json:"treeSha"`
	TargetBranch string `json:"targetBranch"`
	Success      bool   `json:"success"`
}

// FilesMetadataResponse is the response from list-files-with-metadata.
type FilesMetadataResponse struct {
	Files   []FileMetadata        `json:"files"`
	Commits map[string]CommitInfo `json:"commits"`
	Ref     string                `json:"ref"`
}

// FileMetadata holds metadata for a single file entry.
type FileMetadata struct {
	Path          string  `json:"path"`
	Mode          string  `json:"mode"`
	Size          int64   `json:"size"`
	LastCommitSHA *string `json:"last_commit_sha"`
}

// CommitInfo holds summary information about a commit.
type CommitInfo struct {
	Author  string `json:"author"`
	Date    string `json:"date"`
	Message string `json:"message"`
}

// PullUpstreamResponse is the response from pull-upstream.
type PullUpstreamResponse struct {
	Status  string `json:"status"`
	OldSHA  string `json:"oldSha"`
	NewSHA  string `json:"newSha"`
	Branch  string `json:"branch"`
	Success bool   `json:"success"`
}

// DetachResponse is the response from detach-upstream.
type DetachResponse struct {
	Message string `json:"message"`
}

func (r *ReposResource) CreateNote(ctx context.Context, repoID string, opts *CreateNoteOptions) (*NoteResponse, error) {
	body := map[string]any{
		"sha": opts.SHA, "action": "add", "note": opts.Note, "author": opts.Author,
	}
	if opts.ExpectedRefSHA != nil {
		body["expectedRefSha"] = *opts.ExpectedRefSHA
	}
	raw, err := r.client.post(ctx, "/repos/"+repoID+"/notes", body)
	if err != nil {
		return nil, err
	}
	var resp NoteResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal note response: %w", err)
	}
	return &resp, nil
}

func (r *ReposResource) AppendNote(ctx context.Context, repoID, sha, note string, author Identity) (*NoteResponse, error) {
	body := map[string]any{
		"sha": sha, "action": "append", "note": note, "author": author,
	}
	raw, err := r.client.post(ctx, "/repos/"+repoID+"/notes", body)
	if err != nil {
		return nil, err
	}
	var resp NoteResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal note response: %w", err)
	}
	return &resp, nil
}

func (r *ReposResource) GetNote(ctx context.Context, repoID, sha string) (*NoteResponse, error) {
	raw, err := r.client.get(ctx, fmt.Sprintf("/repos/%s/notes/%s", repoID, sha), nil)
	if err != nil {
		return nil, err
	}
	var resp NoteResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal note: %w", err)
	}
	return &resp, nil
}

func (r *ReposResource) DeleteNote(ctx context.Context, repoID, sha string, author *Identity) (*NoteResponse, error) {
	var body any
	if author != nil {
		body = map[string]any{"author": author}
	}
	raw, err := r.client.delWithBody(ctx, fmt.Sprintf("/repos/%s/notes/%s", repoID, sha), body)
	if err != nil {
		return nil, err
	}
	var resp NoteResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal note: %w", err)
	}
	return &resp, nil
}

func (r *ReposResource) RestoreCommit(ctx context.Context, repoID string, opts *RestoreCommitOptions) (*RestoreCommitResponse, error) {
	raw, err := r.client.post(ctx, "/repos/"+repoID+"/restore-commit", opts)
	if err != nil {
		return nil, err
	}
	var resp RestoreCommitResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal restore-commit: %w", err)
	}
	return &resp, nil
}

func (r *ReposResource) ListFilesWithMetadata(ctx context.Context, repoID string, ref *string, ephemeral *bool) (*FilesMetadataResponse, error) {
	q := url.Values{}
	if ref != nil {
		q.Set("ref", *ref)
	}
	if ephemeral != nil {
		q.Set("ephemeral", fmt.Sprintf("%t", *ephemeral))
	}
	raw, err := r.client.get(ctx, "/repos/"+repoID+"/files/metadata", q)
	if err != nil {
		return nil, err
	}
	var resp FilesMetadataResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal files-metadata: %w", err)
	}
	return &resp, nil
}

func (r *ReposResource) PullUpstream(ctx context.Context, repoID string, opts *PullUpstreamOptions) (*PullUpstreamResponse, error) {
	body := map[string]any{}
	if opts != nil && opts.Branch != nil {
		body["branch"] = *opts.Branch
	}
	raw, err := r.client.post(ctx, "/repos/"+repoID+"/pull-upstream", body)
	if err != nil {
		return nil, err
	}
	var resp PullUpstreamResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal pull-upstream: %w", err)
	}
	return &resp, nil
}

func (r *ReposResource) DetachUpstream(ctx context.Context, repoID string) (*DetachResponse, error) {
	if err := r.client.del(ctx, "/repos/"+repoID+"/base", nil); err != nil {
		return nil, err
	}
	return &DetachResponse{Message: "repository detached"}, nil
}

func (r *ReposResource) GetRawFile(ctx context.Context, repoID, ref, path string) ([]byte, error) {
	q := url.Values{}
	q.Set("path", path)
	return r.client.getRaw(ctx, fmt.Sprintf("/repos/%s/raw/%s", repoID, ref), q)
}
