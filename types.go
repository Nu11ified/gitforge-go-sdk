package gitforge

// Repo represents a repository.
type Repo struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Slug                *string  `json:"slug,omitempty"`
	OwnerSlug           *string  `json:"ownerSlug,omitempty"`
	Description         *string  `json:"description,omitempty"`
	Visibility          string   `json:"visibility"`
	DefaultBranch       string   `json:"defaultBranch"`
	LFSEnabled          bool     `json:"lfsEnabled"`
	IsArchived          bool     `json:"isArchived"`
	ForkedFrom          *string  `json:"forkedFrom,omitempty"`
	CreatedAt           *string  `json:"createdAt,omitempty"`
	UpdatedAt           *string  `json:"updatedAt,omitempty"`
	StarCount           *int     `json:"starCount,omitempty"`
	OpenPrCount         *int     `json:"openPrCount,omitempty"`
	OpenIssueCount      *int     `json:"openIssueCount,omitempty"`
	Topics              []string `json:"topics,omitempty"`
	MergeCommitTemplate *string  `json:"mergeCommitTemplate,omitempty"`
}

// Branch represents a branch reference.
type Branch struct {
	Name      string  `json:"name"`
	SHA       string  `json:"sha"`
	ExpiresAt *string `json:"expiresAt,omitempty"`
}

// Tag represents a tag reference.
type Tag struct {
	Name string `json:"name"`
	SHA  string `json:"sha"`
}

// Commit represents a commit object.
type Commit struct {
	SHA         string   `json:"sha"`
	Message     string   `json:"message"`
	Author      string   `json:"author"`
	AuthorEmail string   `json:"authorEmail"`
	Date        string   `json:"date"`
	ParentSHAs  []string `json:"parentShas"`
}

// CommitDetail is a commit with file changes.
type CommitDetail struct {
	SHA         string       `json:"sha"`
	Message     string       `json:"message"`
	Author      string       `json:"author"`
	AuthorEmail string       `json:"authorEmail"`
	Date        string       `json:"date"`
	ParentSHAs  []string     `json:"parentShas"`
	Tree        string       `json:"tree"`
	Files       []FileChange `json:"files"`
}

// FileChange represents a file changed in a commit.
type FileChange struct {
	Path   string `json:"path"`
	Status string `json:"status"`
}

// DiffEntry represents a diff for a file comparison.
type DiffEntry struct {
	Path      string `json:"path"`
	Status    string `json:"status"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Patch     string `json:"patch"`
}

// TreeEntry represents a directory listing entry.
type TreeEntry struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Mode string `json:"mode"`
	SHA  string `json:"sha"`
}

// BlobContent represents file content.
type BlobContent struct {
	Content string `json:"content"`
	Size    int    `json:"size"`
}

// SearchMatch is a code search match within a file.
type SearchMatch struct {
	Line      int    `json:"line"`
	Content   string `json:"content"`
	Highlight string `json:"highlight"`
}

// SearchResult is a code search result for a single file.
type SearchResult struct {
	RepoID   string        `json:"repoId"`
	RepoName string        `json:"repoName"`
	FilePath string        `json:"filePath"`
	Branch   string        `json:"branch"`
	Language *string       `json:"language"`
	Matches  []SearchMatch `json:"matches"`
}

// SearchCodeResult is the response from a code search.
type SearchCodeResult struct {
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
	Page    int            `json:"page"`
	PerPage int            `json:"perPage"`
}

// Comparison represents a comparison between two refs.
type Comparison struct {
	Ahead   int             `json:"ahead"`
	Behind  int             `json:"behind"`
	Commits []CommitSummary `json:"commits"`
	Files   []FileChange    `json:"files"`
}

// CommitSummary is a short commit representation in comparisons.
type CommitSummary struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	Author  string `json:"author"`
	Date    string `json:"date"`
}

// CommitResult is the result of a direct commit API call.
type CommitResult struct {
	CommitSHA  string   `json:"commitSha"`
	TreeSHA    string   `json:"treeSha"`
	Branch     string   `json:"branch"`
	Ref        string   `json:"ref"`
	ParentSHAs []string `json:"parentShas"`
	OldSHA     string   `json:"oldSha"`
	NewSHA     string   `json:"newSha"`
}

// RepoToken is a repo-scoped token.
type RepoToken struct {
	Token     string `json:"token"`
	PatID     string `json:"patId"`
	ExpiresAt string `json:"expiresAt"`
	RemoteURL string `json:"remoteUrl"`
}

// GitCredential represents a git credential.
type GitCredential struct {
	ID        string  `json:"id"`
	Provider  string  `json:"provider"`
	Username  *string `json:"username"`
	Label     *string `json:"label"`
	CreatedAt string  `json:"createdAt"`
}

// MirrorConfig represents a mirror configuration.
type MirrorConfig struct {
	ID           string  `json:"id"`
	SourceURL    string  `json:"sourceUrl"`
	Interval     int     `json:"interval"`
	LastSyncAt   *string `json:"lastSyncAt"`
	LastError    *string `json:"lastError"`
	Enabled      bool    `json:"enabled"`
	CreatedAt    string  `json:"createdAt"`
	Direction    string  `json:"direction"`
	Provider     string  `json:"provider"`
	CredentialID *string `json:"credentialId"`
}

// Webhook represents a webhook configuration.
type Webhook struct {
	ID     string   `json:"id"`
	URL    string   `json:"url"`
	Events []string `json:"events"`
	Active bool     `json:"active"`
}

// WebhookDelivery represents a webhook delivery log.
type WebhookDelivery struct {
	ID             string  `json:"id"`
	EventType      string  `json:"eventType"`
	Payload        string  `json:"payload"`
	ResponseStatus *int    `json:"responseStatus"`
	ResponseBody   *string `json:"responseBody"`
	DeliveredAt    *string `json:"deliveredAt"`
	CreatedAt      string  `json:"createdAt"`
}

// WebhookTestResult is the result of a webhook test.
type WebhookTestResult struct {
	Success      bool    `json:"success"`
	Status       *int    `json:"status"`
	ResponseBody *string `json:"responseBody"`
	DurationMs   int     `json:"durationMs"`
	Error        *string `json:"error"`
}

// PromoteResult is the result of promoting an ephemeral branch.
type PromoteResult struct {
	TargetBranch string `json:"targetBranch"`
	CommitSHA    string `json:"commitSha"`
}

// PaginatedResponse wraps paginated API responses.
type PaginatedResponse[T any] struct {
	Data    []T  `json:"data"`
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasMore bool `json:"hasMore"`
}

// SyncResult is returned by mirror sync.
type SyncResult struct {
	Message string `json:"message"`
}

// CommitFileEntry is a file entry for direct commits.
type CommitFileEntry struct {
	Path     string `json:"path"`
	Content  string `json:"content"`
	Encoding string `json:"encoding,omitempty"`
	Mode     string `json:"mode,omitempty"`
}
