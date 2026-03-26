package gitforge

import (
	"errors"
	"fmt"
)

// GitForgeError represents an API error response.
type GitForgeError struct {
	StatusCode int
	Code       string
	Message    string
}

func (e *GitForgeError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// RefUpdateError is returned on 409 Conflict when a branch was moved (CAS failure).
type RefUpdateError struct {
	*GitForgeError
	CurrentSHA string
}

// IsRefUpdateError checks if err is a *RefUpdateError and returns it.
func IsRefUpdateError(err error) (*RefUpdateError, bool) {
	var rue *RefUpdateError
	if errors.As(err, &rue) {
		return rue, true
	}
	return nil, false
}

// IsGitForgeError checks if err is a *GitForgeError and returns it.
func IsGitForgeError(err error) (*GitForgeError, bool) {
	var gfe *GitForgeError
	if errors.As(err, &gfe) {
		return gfe, true
	}
	return nil, false
}
