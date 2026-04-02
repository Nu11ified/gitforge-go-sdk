package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

// ---------------------------------------------------------------------------
// Edit types
// ---------------------------------------------------------------------------

// ApplyEditsOptions are options for applying structured edits.
type ApplyEditsOptions struct {
	Ref       *string        `json:"ref,omitempty"`
	Edits     []EditEntry    `json:"edits"`
	Commit    bool           `json:"commit,omitempty"`
	Validate  bool           `json:"validate,omitempty"`
	Message   *string        `json:"message,omitempty"`
	Author    *EditAuthor    `json:"author,omitempty"`
	SessionID *string        `json:"sessionId,omitempty"`
}

// EditAuthor identifies the edit/commit author.
type EditAuthor struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// EditEntry is a single edit operation. Set exactly one of TextPatch,
// Metadata, or BinaryPatch and the matching Type string.
type EditEntry struct {
	Type string `json:"type"` // "text-patch", "metadata", or "binary-patch"

	// text-patch fields
	Path        string       `json:"path"`
	Anchor      *TextAnchor  `json:"anchor,omitempty"`
	Mode        *string      `json:"mode,omitempty"`
	Content     *string      `json:"content,omitempty"`
	Indentation interface{}  `json:"indentation,omitempty"` // "auto", "preserve", or IndentationSpec

	// metadata fields
	Format     *string              `json:"format,omitempty"`
	Operations []MetadataOperation  `json:"operations,omitempty"`

	// binary-patch fields
	Patches       []BinaryPatchEntry `json:"patches,omitempty"`
	ValidatePatch *bool              `json:"validate,omitempty"`
}

// TextAnchor locates lines to edit. Use either StartLine/EndLine or Pattern.
type TextAnchor struct {
	StartLine *int    `json:"startLine,omitempty"`
	EndLine   *int    `json:"endLine,omitempty"`
	Pattern   *string `json:"pattern,omitempty"`
	Offset    *int    `json:"offset,omitempty"`
}

// MetadataOperation is a single structured-data mutation.
type MetadataOperation struct {
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

// BinaryPatchEntry describes a single binary patch segment.
type BinaryPatchEntry struct {
	Offset   int    `json:"offset"`
	Data     string `json:"data"`
	Encoding string `json:"encoding,omitempty"` // "base64" or "hex"
}

// IndentationSpec explicitly sets indentation style.
type IndentationSpec struct {
	Style string `json:"style"` // "tabs" or "spaces"
	Width int    `json:"width"`
}

// ---------------------------------------------------------------------------
// Apply results
// ---------------------------------------------------------------------------

// ApplyEditsResult is the response from applying edits.
type ApplyEditsResult struct {
	OK         bool             `json:"ok"`
	Results    []EditResultItem `json:"results"`
	Validation *ValidationResult `json:"validation,omitempty"`
	Commit     *EditCommitInfo  `json:"commit,omitempty"`
}

// EditResultItem is the per-file result of an edit.
type EditResultItem struct {
	Path    string  `json:"path"`
	Status  string  `json:"status"` // "applied" or "error"
	Diff    *string `json:"diff,omitempty"`
	Changes *int    `json:"changes,omitempty"`
	Error   *string `json:"error,omitempty"`
}

// ValidationResult holds LSP validation output.
type ValidationResult struct {
	OK          bool                   `json:"ok"`
	Diagnostics []ValidationDiagnostic `json:"diagnostics"`
}

// ValidationDiagnostic is a single LSP diagnostic.
type ValidationDiagnostic struct {
	File     string  `json:"file"`
	Line     int     `json:"line"`
	Column   int     `json:"column"`
	Severity string  `json:"severity"` // "error", "warning", "info"
	Message  string  `json:"message"`
	Source   *string `json:"source,omitempty"`
}

// EditCommitInfo is the commit created when commit=true.
type EditCommitInfo struct {
	CommitSHA string `json:"commitSha"`
	TreeSHA   string `json:"treeSha"`
	Branch    string `json:"branch"`
	Ref       string `json:"ref"`
}

// ---------------------------------------------------------------------------
// Context types
// ---------------------------------------------------------------------------

// ContextOptions are options for reading file context.
type ContextOptions struct {
	Paths            []string
	Ref              *string
	SurroundingLines *int
}

// ContextResult is the response from the context endpoint.
type ContextResult struct {
	Ref       string             `json:"ref"`
	CommitSHA string             `json:"commitSha"`
	Files     []ContextFileEntry `json:"files"`
}

// ContextFileEntry is a single file in the context response.
type ContextFileEntry struct {
	Path    string  `json:"path"`
	Content *string `json:"content"`
	Lines   int     `json:"lines"`
	Error   *string `json:"error,omitempty"`
}

// ---------------------------------------------------------------------------
// Session types
// ---------------------------------------------------------------------------

// CreateSessionOptions are options for creating an edit session.
type CreateSessionOptions struct {
	RepoID      string  `json:"repoId"`
	Branch      string  `json:"branch"`
	SourceRef   string  `json:"sourceRef,omitempty"`
	Description *string `json:"description,omitempty"`
	TTLHours    *int    `json:"ttlHours,omitempty"`
}

// CreateSessionResult is the response from creating a session.
type CreateSessionResult struct {
	ID              string `json:"id"`
	Branch          string `json:"branch"`
	SourceRef       string `json:"sourceRef"`
	SourceCommitSHA string `json:"sourceCommitSha"`
	ExpiresAt       string `json:"expiresAt"`
}

// EditSession is the full session status.
type EditSession struct {
	ID          string   `json:"id"`
	RepoID      string   `json:"repoId"`
	Branch      string   `json:"branch"`
	SourceRef   string   `json:"sourceRef"`
	Status      string   `json:"status"`
	Description *string  `json:"description,omitempty"`
	Commits     []string `json:"commits"`
	CreatedAt   *string  `json:"createdAt,omitempty"`
	UpdatedAt   *string  `json:"updatedAt,omitempty"`
	ExpiresAt   *string  `json:"expiresAt,omitempty"`
}

// SubmitSessionOptions are options for submitting a session as a PR.
type SubmitSessionOptions struct {
	Title        string  `json:"title"`
	Body         *string `json:"body,omitempty"`
	TargetBranch string  `json:"targetBranch"`
}

// SubmitSessionResult is the response from submitting a session.
type SubmitSessionResult struct {
	PRID      string `json:"prId"`
	PRNumber  int    `json:"prNumber"`
	SessionID string `json:"sessionId"`
	Status    string `json:"status"` // "submitted"
}

// ---------------------------------------------------------------------------
// Resource
// ---------------------------------------------------------------------------

// EditResource provides structured file editing operations.
type EditResource struct {
	client *httpClient
}

// Apply applies structured edits to repository files.
func (r *EditResource) Apply(ctx context.Context, repoID string, opts *ApplyEditsOptions) (*ApplyEditsResult, error) {
	raw, err := r.client.post(ctx, "/edit/"+repoID, opts)
	if err != nil {
		return nil, err
	}
	var result ApplyEditsResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal apply edits: %w", err)
	}
	return &result, nil
}

// Context reads file context for agent-driven editing.
func (r *EditResource) Context(ctx context.Context, repoID string, opts *ContextOptions) (*ContextResult, error) {
	q := url.Values{}
	q.Set("paths", strings.Join(opts.Paths, ","))
	if opts.Ref != nil {
		q.Set("ref", *opts.Ref)
	}
	if opts.SurroundingLines != nil {
		q.Set("surroundingLines", fmt.Sprintf("%d", *opts.SurroundingLines))
	}
	raw, err := r.client.get(ctx, "/edit/"+repoID+"/context", q)
	if err != nil {
		return nil, err
	}
	var result ContextResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal context: %w", err)
	}
	return &result, nil
}

// CreateSession creates a multi-step edit session with its own branch.
func (r *EditResource) CreateSession(ctx context.Context, opts *CreateSessionOptions) (*CreateSessionResult, error) {
	raw, err := r.client.post(ctx, "/edit/sessions", opts)
	if err != nil {
		return nil, err
	}
	var result CreateSessionResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal create session: %w", err)
	}
	return &result, nil
}

// GetSession retrieves the current status of an edit session.
func (r *EditResource) GetSession(ctx context.Context, sessionID string) (*EditSession, error) {
	raw, err := r.client.get(ctx, "/edit/sessions/"+sessionID, nil)
	if err != nil {
		return nil, err
	}
	var result EditSession
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}
	return &result, nil
}

// SubmitSession submits a completed edit session as a pull request.
func (r *EditResource) SubmitSession(ctx context.Context, sessionID string, opts *SubmitSessionOptions) (*SubmitSessionResult, error) {
	raw, err := r.client.post(ctx, "/edit/sessions/"+sessionID+"/submit", opts)
	if err != nil {
		return nil, err
	}
	var result SubmitSessionResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal submit session: %w", err)
	}
	return &result, nil
}
