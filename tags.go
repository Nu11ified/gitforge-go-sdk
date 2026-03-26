package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type ListTagsOptions struct {
	Limit  *int
	Offset *int
}

type CreateTagOptions struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
}

type TagsResource struct {
	client *httpClient
	repoID string
}

func (t *TagsResource) List(ctx context.Context, opts *ListTagsOptions) (*PaginatedResponse[Tag], error) {
	q := url.Values{}
	if opts != nil {
		if opts.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *opts.Limit))
		}
		if opts.Offset != nil {
			q.Set("offset", fmt.Sprintf("%d", *opts.Offset))
		}
	}
	raw, err := t.client.get(ctx, fmt.Sprintf("/repos/%s/tags", t.repoID), q)
	if err != nil {
		return nil, err
	}
	var result PaginatedResponse[Tag]
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal tags: %w", err)
	}
	return &result, nil
}

func (t *TagsResource) Create(ctx context.Context, opts *CreateTagOptions) (*Tag, error) {
	raw, err := t.client.post(ctx, fmt.Sprintf("/repos/%s/tags", t.repoID), opts)
	if err != nil {
		return nil, err
	}
	var tag Tag
	if err := json.Unmarshal(raw, &tag); err != nil {
		return nil, fmt.Errorf("unmarshal tag: %w", err)
	}
	return &tag, nil
}

func (t *TagsResource) Delete(ctx context.Context, name string) error {
	return t.client.del(ctx, fmt.Sprintf("/repos/%s/tags/%s", t.repoID, url.PathEscape(name)), nil)
}
