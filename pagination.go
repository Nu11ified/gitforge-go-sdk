package gitforge

import "context"

// PageFetcher fetches a page of results.
type PageFetcher[T any] func(ctx context.Context, limit, offset int) (*PaginatedResponse[T], error)

// PaginateOptions configures pagination behavior.
type PaginateOptions struct {
	PageSize int
	MaxItems *int
}

// PaginateItem wraps a single item or error from pagination.
type PaginateItem[T any] struct {
	Item  T
	Error error
}

// Paginate returns a channel that yields items from paginated API responses.
func Paginate[T any](ctx context.Context, fetcher PageFetcher[T], opts *PaginateOptions) <-chan PaginateItem[T] {
	ch := make(chan PaginateItem[T])
	pageSize := 20
	var maxItems int = -1

	if opts != nil {
		if opts.PageSize > 0 {
			pageSize = opts.PageSize
		}
		if opts.MaxItems != nil {
			maxItems = *opts.MaxItems
		}
	}

	go func() {
		defer close(ch)
		offset := 0
		yielded := 0

		for {
			if ctx.Err() != nil {
				return
			}

			page, err := fetcher(ctx, pageSize, offset)
			if err != nil {
				select {
				case ch <- PaginateItem[T]{Error: err}:
				case <-ctx.Done():
				}
				return
			}

			for _, item := range page.Data {
				if maxItems >= 0 && yielded >= maxItems {
					return
				}
				select {
				case ch <- PaginateItem[T]{Item: item}:
					yielded++
				case <-ctx.Done():
					return
				}
			}

			if !page.HasMore || len(page.Data) == 0 {
				return
			}
			offset += len(page.Data)
		}
	}()

	return ch
}
