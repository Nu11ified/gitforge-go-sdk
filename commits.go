package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type ListCommitsOptions struct {
	Ref       *string
	Ephemeral *bool
	Limit     *int
}

type CreateCommitOptions struct {
	Branch      string
	Message     string
	AuthorName  string
	AuthorEmail string
	BaseBranch  *string
}

type FileOptions struct {
	Encoding string
	Mode     string
}

type CommitBuilder struct {
	client      *httpClient
	repoID      string
	branch      string
	message     string
	authorName  string
	authorEmail string
	baseBranch  *string
	files       []commitFilePayload
	deletes     []string
	ephemeral   bool
	headSHA     *string
}

type commitFilePayload struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Encoding string `json:"encoding,omitempty"`
	Mode     string `json:"mode,omitempty"`
}

func (cb *CommitBuilder) AddFile(path, content string, opts *FileOptions) *CommitBuilder {
	f := commitFilePayload{Path: path, Content: content}
	if opts != nil {
		f.Encoding = opts.Encoding
		f.Mode = opts.Mode
	}
	cb.files = append(cb.files, f)
	return cb
}

func (cb *CommitBuilder) DeleteFile(path string) *CommitBuilder {
	cb.deletes = append(cb.deletes, path)
	return cb
}

func (cb *CommitBuilder) Ephemeral(value bool) *CommitBuilder {
	cb.ephemeral = value
	return cb
}

func (cb *CommitBuilder) ExpectedHeadSHA(sha string) *CommitBuilder {
	cb.headSHA = &sha
	return cb
}

func (cb *CommitBuilder) Send(ctx context.Context) (*CommitResult, error) {
	body := map[string]any{
		"branch":  cb.branch,
		"message": cb.message,
		"author":  map[string]string{"name": cb.authorName, "email": cb.authorEmail},
		"files":   cb.files,
		"deletes": cb.deletes,
	}
	if cb.baseBranch != nil {
		body["baseBranch"] = *cb.baseBranch
	}
	if cb.ephemeral {
		body["ephemeral"] = true
	}
	if cb.headSHA != nil {
		body["expectedHeadSha"] = *cb.headSHA
	}
	raw, err := cb.client.post(ctx, fmt.Sprintf("/repos/%s/commits", cb.repoID), body)
	if err != nil {
		return nil, err
	}
	var result CommitResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal commit result: %w", err)
	}
	return &result, nil
}

type CommitsResource struct {
	client *httpClient
	repoID string
}

func (c *CommitsResource) List(ctx context.Context, opts *ListCommitsOptions) ([]Commit, error) {
	q := url.Values{}
	if opts != nil {
		if opts.Ref != nil {
			q.Set("ref", *opts.Ref)
		}
		if opts.Ephemeral != nil && *opts.Ephemeral {
			q.Set("ephemeral", "true")
		}
		if opts.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *opts.Limit))
		}
	}
	raw, err := c.client.get(ctx, fmt.Sprintf("/repos/%s/commits", c.repoID), q)
	if err != nil {
		return nil, err
	}
	var commits []Commit
	if err := json.Unmarshal(raw, &commits); err != nil {
		return nil, fmt.Errorf("unmarshal commits: %w", err)
	}
	return commits, nil
}

func (c *CommitsResource) Get(ctx context.Context, sha string) (*CommitDetail, error) {
	raw, err := c.client.get(ctx, fmt.Sprintf("/repos/%s/commits/%s", c.repoID, sha), nil)
	if err != nil {
		return nil, err
	}
	var detail CommitDetail
	if err := json.Unmarshal(raw, &detail); err != nil {
		return nil, fmt.Errorf("unmarshal commit detail: %w", err)
	}
	return &detail, nil
}

func (c *CommitsResource) GetDiff(ctx context.Context, sha string) ([]DiffEntry, error) {
	raw, err := c.client.get(ctx, fmt.Sprintf("/repos/%s/commits/%s/diff", c.repoID, sha), nil)
	if err != nil {
		return nil, err
	}
	var entries []DiffEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("unmarshal diff: %w", err)
	}
	return entries, nil
}

func (c *CommitsResource) Create(opts *CreateCommitOptions) *CommitBuilder {
	return &CommitBuilder{
		client:      c.client,
		repoID:      c.repoID,
		branch:      opts.Branch,
		message:     opts.Message,
		authorName:  opts.AuthorName,
		authorEmail: opts.AuthorEmail,
		baseBranch:  opts.BaseBranch,
	}
}
