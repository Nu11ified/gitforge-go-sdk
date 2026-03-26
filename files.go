package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type ListFilesOptions struct {
	Path      *string
	Ephemeral *bool
}

type GetFileOptions struct {
	Ephemeral *bool
}

type FilesResource struct {
	client *httpClient
	repoID string
}

func (f *FilesResource) ListFiles(ctx context.Context, ref string, opts *ListFilesOptions) ([]TreeEntry, error) {
	q := url.Values{}
	if opts != nil {
		if opts.Path != nil {
			q.Set("path", *opts.Path)
		}
		if opts.Ephemeral != nil && *opts.Ephemeral {
			q.Set("ephemeral", "true")
		}
	}
	raw, err := f.client.get(ctx, fmt.Sprintf("/repos/%s/tree/%s", f.repoID, ref), q)
	if err != nil {
		return nil, err
	}
	var entries []TreeEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("unmarshal tree: %w", err)
	}
	return entries, nil
}

func (f *FilesResource) GetFile(ctx context.Context, ref, path string, opts *GetFileOptions) (*BlobContent, error) {
	q := url.Values{"path": {path}}
	if opts != nil && opts.Ephemeral != nil && *opts.Ephemeral {
		q.Set("ephemeral", "true")
	}
	raw, err := f.client.get(ctx, fmt.Sprintf("/repos/%s/blob/%s", f.repoID, ref), q)
	if err != nil {
		return nil, err
	}
	var blob BlobContent
	if err := json.Unmarshal(raw, &blob); err != nil {
		return nil, fmt.Errorf("unmarshal blob: %w", err)
	}
	return &blob, nil
}
