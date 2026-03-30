package gitforge

import "net/http"

// ClientOptions configures the GitForge client.
type ClientOptions struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// Client is the GitForge API client.
type Client struct {
	Repos     *ReposResource
	PatchSets *PatchSetsResource
	http      *httpClient
}

// NewClient creates a new GitForge API client.
func NewClient(opts ClientOptions) *Client {
	h := newHTTPClient(opts.BaseURL, opts.Token, opts.HTTPClient)
	return &Client{
		Repos:     &ReposResource{client: h},
		PatchSets: &PatchSetsResource{client: h},
		http:      h,
	}
}

// Ptr returns a pointer to v. Useful for optional fields in option structs.
func Ptr[T any](v T) *T {
	return &v
}

// Repo returns a RepoScope for the given repository ID.
func (c *Client) Repo(repoID string) *RepoScope {
	return &RepoScope{
		Branches:    &BranchesResource{client: c.http, repoID: repoID},
		Tags:        &TagsResource{client: c.http, repoID: repoID},
		Commits:     &CommitsResource{client: c.http, repoID: repoID},
		Files:       &FilesResource{client: c.http, repoID: repoID},
		Search:      &SearchResource{client: c.http, repoID: repoID},
		Tokens:      &TokensResource{client: c.http, repoID: repoID},
		Credentials: &CredentialsResource{client: c.http, repoID: repoID},
		Mirrors:     &MirrorsResource{client: c.http, repoID: repoID},
		Webhooks:    &WebhooksResource{client: c.http, repoID: repoID},
		Sandbox:     &SandboxResource{client: c.http, repoID: repoID},
	}
}

// RepoScope provides access to repo-scoped resources.
type RepoScope struct {
	Branches    *BranchesResource
	Tags        *TagsResource
	Commits     *CommitsResource
	Files       *FilesResource
	Search      *SearchResource
	Tokens      *TokensResource
	Credentials *CredentialsResource
	Mirrors     *MirrorsResource
	Webhooks    *WebhooksResource
	Sandbox     *SandboxResource
}
