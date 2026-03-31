package gitforge

import (
	"context"
	"encoding/json"
	"fmt"
)

// OperationInfo represents a recorded operation in the operation log.
type OperationInfo struct {
	ID            string                 `json:"id"`
	ChangeID      *string                `json:"changeId,omitempty"`
	RepoID        string                 `json:"repoId"`
	OperationType string                 `json:"operationType"`
	PreviousState map[string]interface{} `json:"previousState"`
	NewState      map[string]interface{} `json:"newState"`
	CreatedAt     string                 `json:"createdAt"`
}

// OperationsService provides operation log and undo operations.
type OperationsService struct {
	client *httpClient
}

func (s *OperationsService) List(ctx context.Context, repoID string) ([]OperationInfo, error) {
	raw, err := s.client.get(ctx, fmt.Sprintf("/repos/%s/operations", repoID), nil)
	if err != nil {
		return nil, err
	}
	var result struct {
		Items []OperationInfo `json:"items"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal operations: %w", err)
	}
	return result.Items, nil
}

func (s *OperationsService) Undo(ctx context.Context, repoID string, operationID *string) (*OperationInfo, error) {
	body := map[string]interface{}{}
	if operationID != nil {
		body["operationId"] = *operationID
	}
	raw, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/operations/undo", repoID), body)
	if err != nil {
		return nil, err
	}
	var result struct {
		UndoneOperation OperationInfo `json:"undoneOperation"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("unmarshal undo result: %w", err)
	}
	return &result.UndoneOperation, nil
}

func (s *OperationsService) Restore(ctx context.Context, repoID, operationID string) error {
	_, err := s.client.post(ctx, fmt.Sprintf("/repos/%s/operations/%s/restore", repoID, operationID), nil)
	return err
}
