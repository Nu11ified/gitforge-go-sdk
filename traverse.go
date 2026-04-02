package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ---------------------------------------------------------------------------
// Option types
// ---------------------------------------------------------------------------

// ListTraverseReposOptions are options for listing repos with index summaries.
type ListTraverseReposOptions struct {
	Q           *string
	Language    *string
	BuildSystem *string
	IsMonorepo  *bool
	Sort        *string // "recent", "name", "size", "fileCount"
	Limit       *int
	Offset      *int
}

// GetTraverseRepoOptions are options for retrieving repo traversal data.
type GetTraverseRepoOptions struct {
	Ref     *string  // branch/tag/SHA, defaults to "main"
	Depth   *string  // "L1", "L2", or "L3"
	Path    *string  // subtree path filter
	Include []string // specific sections to include (e.g. "tree", "symbols")
}

// ImpactOptions are options for impact analysis.
type ImpactOptions struct {
	Paths []string // required — files to analyse
	Ref   *string  // branch/tag/SHA, defaults to "main"
}

// ---------------------------------------------------------------------------
// Result types
// ---------------------------------------------------------------------------

// TraverseRepoSummary is a repo entry returned by the repos listing.
type TraverseRepoSummary struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Slug          *string            `json:"slug,omitempty"`
	DefaultBranch string             `json:"defaultBranch"`
	FileCount     *int               `json:"fileCount,omitempty"`
	Languages     map[string]float64 `json:"languages,omitempty"`
	BuildSystem   *string            `json:"buildSystem,omitempty"`
	IsMonorepo    *bool              `json:"isMonorepo,omitempty"`
	IndexStatus   *IndexStatus       `json:"indexStatus,omitempty"`
}

// IndexStatus describes the build state of each index level.
type IndexStatus struct {
	L1 string `json:"l1"`
	L2 string `json:"l2"`
	L3 string `json:"l3"`
}

// TraverseRepoResult is the full traversal response for a single repo.
type TraverseRepoResult struct {
	IndexStatus   IndexStatus              `json:"indexStatus"`
	Head          string                   `json:"head"`
	Ref           string                   `json:"ref"`
	Tree          []TraverseTreeEntry      `json:"tree,omitempty"`
	Configs       map[string]interface{}   `json:"configs,omitempty"`
	Languages     map[string]float64       `json:"languages,omitempty"`
	BuildSystem   *string                  `json:"buildSystem,omitempty"`
	IsMonorepo    *bool                    `json:"isMonorepo,omitempty"`
	FileCount     *int                     `json:"fileCount,omitempty"`
	Symbols       map[string]FileSymbols   `json:"symbols,omitempty"`
	EntryPoints   []string                 `json:"entryPoints,omitempty"`
	TestMap       map[string]string        `json:"testMap,omitempty"`
	Architecture  *ArchitectureInfo        `json:"architecture,omitempty"`
	RelevanceTags []string                 `json:"relevanceTags,omitempty"`
	Summaries     []FileSummary            `json:"summaries,omitempty"`
}

// TraverseTreeEntry is a file or directory in the traversal tree.
type TraverseTreeEntry struct {
	Path     string  `json:"path"`
	Type     string  `json:"type"` // "blob" or "tree"
	Mode     string  `json:"mode"`
	Size     *int64  `json:"size,omitempty"`
	SHA      *string `json:"sha,omitempty"`
	Language *string `json:"language,omitempty"`
}

// FileSymbols holds exports and imports for a single file.
type FileSymbols struct {
	Exports []SymbolEntry `json:"exports"`
	Imports []ImportEntry `json:"imports"`
}

// SymbolEntry is an exported symbol.
type SymbolEntry struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
	Line int    `json:"line"`
}

// ImportEntry is an import statement.
type ImportEntry struct {
	Source  string   `json:"source"`
	Symbols []string `json:"symbols"`
}

// ArchitectureInfo describes high-level architecture.
type ArchitectureInfo struct {
	Layers      []string `json:"layers"`
	EntryPoints []string `json:"entryPoints"`
	Description *string  `json:"description,omitempty"`
}

// FileSummary is a per-file natural-language summary.
type FileSummary struct {
	Path          string   `json:"path"`
	Summary       string   `json:"summary"`
	RelevanceTags []string `json:"relevanceTags,omitempty"`
}

// ImpactedFile is a file affected by a set of changes.
type ImpactedFile struct {
	Path   string `json:"path"`
	Reason string `json:"reason"`
	Depth  int    `json:"depth"`
}

// ImpactResult is the response from impact analysis.
type ImpactResult struct {
	Impacted           []ImpactedFile `json:"impacted"`
	TestFiles          []string       `json:"testFiles"`
	TotalImpactedFiles int            `json:"totalImpactedFiles"`
}

// ---------------------------------------------------------------------------
// Resource
// ---------------------------------------------------------------------------

// TraverseResource provides repository traversal operations.
type TraverseResource struct {
	client *httpClient
}

// Repos lists repositories with their index summaries.
func (r *TraverseResource) Repos(ctx context.Context, opts *ListTraverseReposOptions) (*PaginatedResponse[TraverseRepoSummary], error) {
	q := url.Values{}
	if opts != nil {
		if opts.Q != nil {
			q.Set("q", *opts.Q)
		}
		if opts.Language != nil {
			q.Set("language", *opts.Language)
		}
		if opts.BuildSystem != nil {
			q.Set("buildSystem", *opts.BuildSystem)
		}
		if opts.IsMonorepo != nil {
			q.Set("isMonorepo", fmt.Sprintf("%t", *opts.IsMonorepo))
		}
		if opts.Sort != nil {
			q.Set("sort", *opts.Sort)
		}
		if opts.Limit != nil {
			q.Set("limit", fmt.Sprintf("%d", *opts.Limit))
		}
		if opts.Offset != nil {
			q.Set("offset", fmt.Sprintf("%d", *opts.Offset))
		}
	}
	raw, err := r.client.get(ctx, "/traverse/repos", q)
	if err != nil {
		return nil, err
	}
	var result PaginatedResponse[TraverseRepoSummary]
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal traverse repos: %w", err)
	}
	return &result, nil
}

// Repo retrieves traversal data for a single repository.
func (r *TraverseResource) Repo(ctx context.Context, repoID string, opts *GetTraverseRepoOptions) (*TraverseRepoResult, error) {
	q := url.Values{}
	if opts != nil {
		if opts.Ref != nil {
			q.Set("ref", *opts.Ref)
		}
		if opts.Depth != nil {
			q.Set("depth", *opts.Depth)
		}
		if opts.Path != nil {
			q.Set("path", *opts.Path)
		}
		if len(opts.Include) > 0 {
			q.Set("include", strings.Join(opts.Include, ","))
		}
	}
	raw, err := r.client.get(ctx, "/traverse/"+repoID, q)
	if err != nil {
		return nil, err
	}
	var result TraverseRepoResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal traverse repo: %w", err)
	}
	return &result, nil
}

// Impact runs impact analysis for a set of changed file paths.
func (r *TraverseResource) Impact(ctx context.Context, repoID string, opts *ImpactOptions) (*ImpactResult, error) {
	q := url.Values{}
	q.Set("paths", strings.Join(opts.Paths, ","))
	if opts.Ref != nil {
		q.Set("ref", *opts.Ref)
	}
	raw, err := r.client.get(ctx, "/traverse/"+repoID+"/impact", q)
	if err != nil {
		return nil, err
	}
	var result ImpactResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal impact: %w", err)
	}
	return &result, nil
}
