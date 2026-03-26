package gitforge

import (
	"context"
	"errors"
	"testing"
)

func ptrInt(v int) *int { return &v }

// collectItems drains the channel returned by Paginate and returns items and
// the first error encountered (if any).
func collectItems[T any](ch <-chan PaginateItem[T]) ([]T, error) {
	var items []T
	var firstErr error
	for pi := range ch {
		if pi.Error != nil {
			if firstErr == nil {
				firstErr = pi.Error
			}
			continue
		}
		items = append(items, pi.Item)
	}
	return items, firstErr
}

func TestPaginate_SinglePage(t *testing.T) {
	fetcher := func(_ context.Context, limit, offset int) (*PaginatedResponse[int], error) {
		return &PaginatedResponse[int]{
			Data:    []int{1, 2, 3},
			Total:   3,
			Limit:   limit,
			Offset:  offset,
			HasMore: false,
		}, nil
	}

	ch := Paginate(context.Background(), fetcher, nil)
	items, err := collectItems(ch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	for i, v := range items {
		if v != i+1 {
			t.Errorf("item[%d]: expected %d, got %d", i, i+1, v)
		}
	}
}

func TestPaginate_MultiplePages(t *testing.T) {
	pages := []struct {
		data    []int
		hasMore bool
	}{
		{[]int{1, 2}, true},
		{[]int{3, 4}, true},
		{[]int{5}, false},
	}
	call := 0

	fetcher := func(_ context.Context, limit, offset int) (*PaginatedResponse[int], error) {
		p := pages[call]
		call++
		return &PaginatedResponse[int]{
			Data:    p.data,
			Total:   5,
			Limit:   limit,
			Offset:  offset,
			HasMore: p.hasMore,
		}, nil
	}

	ch := Paginate(context.Background(), fetcher, &PaginateOptions{PageSize: 2})
	items, err := collectItems(ch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 5 {
		t.Fatalf("expected 5 items, got %d", len(items))
	}
	for i, v := range items {
		if v != i+1 {
			t.Errorf("item[%d]: expected %d, got %d", i, i+1, v)
		}
	}
}

func TestPaginate_MaxItemsLimit(t *testing.T) {
	// 10 items spread across pages of 4
	allItems := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fetcher := func(_ context.Context, limit, offset int) (*PaginatedResponse[int], error) {
		end := offset + limit
		if end > len(allItems) {
			end = len(allItems)
		}
		data := allItems[offset:end]
		hasMore := end < len(allItems)
		return &PaginatedResponse[int]{
			Data:    data,
			Total:   len(allItems),
			Limit:   limit,
			Offset:  offset,
			HasMore: hasMore,
		}, nil
	}

	ch := Paginate(context.Background(), fetcher, &PaginateOptions{
		PageSize: 4,
		MaxItems: ptrInt(3),
	})
	items, err := collectItems(ch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(items))
	}
	for i, v := range items {
		if v != i+1 {
			t.Errorf("item[%d]: expected %d, got %d", i, i+1, v)
		}
	}
}

func TestPaginate_EmptyFirstPage(t *testing.T) {
	fetcher := func(_ context.Context, limit, offset int) (*PaginatedResponse[int], error) {
		return &PaginatedResponse[int]{
			Data:    []int{},
			Total:   0,
			Limit:   limit,
			Offset:  offset,
			HasMore: false,
		}, nil
	}

	ch := Paginate(context.Background(), fetcher, nil)
	items, err := collectItems(ch)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected 0 items, got %d", len(items))
	}
}

func TestPaginate_ErrorPropagation(t *testing.T) {
	fetchErr := errors.New("server error")
	call := 0

	fetcher := func(_ context.Context, limit, offset int) (*PaginatedResponse[int], error) {
		call++
		if call == 1 {
			return &PaginatedResponse[int]{
				Data:    []int{10, 20},
				Total:   4,
				Limit:   limit,
				Offset:  offset,
				HasMore: true,
			}, nil
		}
		return nil, fetchErr
	}

	ch := Paginate(context.Background(), fetcher, &PaginateOptions{PageSize: 2})
	items, err := collectItems(ch)

	if err == nil {
		t.Fatal("expected an error, got nil")
	}
	if !errors.Is(err, fetchErr) {
		t.Fatalf("expected fetchErr, got %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected 2 items from first page, got %d", len(items))
	}
}
