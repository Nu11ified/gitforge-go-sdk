package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type SearchCodeOptions struct {
	Query    string
	Language *string
	Branch   *string
	PerPage  *int
	Page     *int
}

type SearchResource struct {
	client *httpClient
	repoID string
}

func (s *SearchResource) SearchCode(ctx context.Context, opts *SearchCodeOptions) (*SearchCodeResult, error) {
	q := url.Values{"q": {opts.Query}}
	if opts.Language != nil {
		q.Set("lang", *opts.Language)
	}
	if opts.Branch != nil {
		q.Set("branch", *opts.Branch)
	}
	if opts.PerPage != nil {
		q.Set("perPage", fmt.Sprintf("%d", *opts.PerPage))
	}
	if opts.Page != nil {
		q.Set("page", fmt.Sprintf("%d", *opts.Page))
	}
	raw, err := s.client.get(ctx, fmt.Sprintf("/repos/%s/search", s.repoID), q)
	if err != nil {
		return nil, err
	}
	var result SearchCodeResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal search: %w", err)
	}
	return &result, nil
}

func (s *SearchResource) Compare(ctx context.Context, base, head string) (*Comparison, error) {
	q := url.Values{"base": {base}, "head": {head}}
	raw, err := s.client.get(ctx, fmt.Sprintf("/repos/%s/compare", s.repoID), q)
	if err != nil {
		return nil, err
	}
	var comp Comparison
	if err := json.Unmarshal(raw, &comp); err != nil {
		return nil, fmt.Errorf("unmarshal comparison: %w", err)
	}
	return &comp, nil
}

func (s *SearchResource) CompareDiff(ctx context.Context, base, head string) ([]DiffEntry, error) {
	q := url.Values{"base": {base}, "head": {head}}
	raw, err := s.client.get(ctx, fmt.Sprintf("/repos/%s/compare/diff", s.repoID), q)
	if err != nil {
		return nil, err
	}
	var entries []DiffEntry
	if err := json.Unmarshal(raw, &entries); err != nil {
		return nil, fmt.Errorf("unmarshal diff: %w", err)
	}
	return entries, nil
}
