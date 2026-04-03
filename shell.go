package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
)

// ShellExecOptions configures a single-repo shell command execution.
type ShellExecOptions struct {
	Command   string            `json:"command"`
	SessionID *string           `json:"sessionId,omitempty"`
	Ref       *string           `json:"ref,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
}

// ShellExecResult is the response from a single-repo shell execution.
type ShellExecResult struct {
	SessionID      string  `json:"sessionId"`
	Stdout         string  `json:"stdout"`
	Stderr         string  `json:"stderr"`
	ExitCode       int     `json:"exitCode"`
	Ref            *string `json:"ref"`
	HeadSha        *string `json:"headSha"`
	PendingChanges int     `json:"pendingChanges"`
}

// ShellMountOptions describes a single repo mount for multi-repo shell execution.
type ShellMountOptions struct {
	RepoID string `json:"repoId"`
	Path   string `json:"path"`
	Ref    string `json:"ref"`
}

// ShellMultiExecOptions configures a multi-repo shell command execution.
type ShellMultiExecOptions struct {
	Command   string              `json:"command"`
	SessionID *string             `json:"sessionId,omitempty"`
	Mounts    []ShellMountOptions `json:"mounts,omitempty"`
	Env       map[string]string   `json:"env,omitempty"`
}

// ShellMountResult describes the state of a single mount after multi-repo shell execution.
type ShellMountResult struct {
	Path           string `json:"path"`
	RepoID         string `json:"repoId"`
	Ref            string `json:"ref"`
	HeadSha        string `json:"headSha"`
	PendingChanges int    `json:"pendingChanges"`
}

// ShellMultiExecResult is the response from a multi-repo shell execution.
type ShellMultiExecResult struct {
	SessionID string             `json:"sessionId"`
	Stdout    string             `json:"stdout"`
	Stderr    string             `json:"stderr"`
	ExitCode  int                `json:"exitCode"`
	Mounts    []ShellMountResult `json:"mounts"`
}

// ShellDestroyResult is the response from destroying a shell session.
type ShellDestroyResult struct {
	Destroyed        bool `json:"destroyed"`
	UncommittedFiles int  `json:"uncommittedFiles"`
}

// RepoShellResource provides shell execution operations scoped to a repository.
type RepoShellResource struct {
	client *httpClient
	repoID string
}

// Exec runs a command in a VFS shell session scoped to the repository.
func (r *RepoShellResource) Exec(ctx context.Context, opts *ShellExecOptions) (*ShellExecResult, error) {
	raw, err := r.client.post(ctx, fmt.Sprintf("/repos/%s/shell", r.repoID), opts)
	if err != nil {
		return nil, err
	}
	var result ShellExecResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal shell exec result: %w", err)
	}
	return &result, nil
}

// ShellResource provides top-level (multi-repo) shell operations.
type ShellResource struct {
	client *httpClient
}

// ExecMulti runs a command in a multi-repo VFS shell session.
func (r *ShellResource) ExecMulti(ctx context.Context, opts *ShellMultiExecOptions) (*ShellMultiExecResult, error) {
	raw, err := r.client.post(ctx, "/shell", opts)
	if err != nil {
		return nil, err
	}
	var result ShellMultiExecResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal shell multi exec result: %w", err)
	}
	return &result, nil
}

// Destroy tears down a shell session, discarding any uncommitted changes.
func (r *ShellResource) Destroy(ctx context.Context, sessionID string) (*ShellDestroyResult, error) {
	raw, err := r.client.delWithBody(ctx, fmt.Sprintf("/shell/%s", sessionID), nil)
	if err != nil {
		return nil, err
	}
	var result ShellDestroyResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal shell destroy result: %w", err)
	}
	return &result, nil
}
