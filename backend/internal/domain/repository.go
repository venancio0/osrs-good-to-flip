package domain

import (
	"context"
	"time"
)

// PriceSnapshot represents a single price observation from the provider.
// Volume may be zero when the endpoint does not return it (e.g., /latest).
type PriceSnapshot struct {
	High   int
	Low    int
	Volume int
}

// PriceProvider defines the interface for fetching prices from external sources
type PriceProvider interface {
	FetchLatestPrices(ctx context.Context) (map[int]PriceSnapshot, error)
}

// ItemRepository defines the interface for item data operations
type ItemRepository interface {
	SavePrices(ctx context.Context, prices []ItemPrice) error
	GetItemByID(ctx context.Context, id int) (*ItemPrice, error)
	SearchItems(ctx context.Context, query string) ([]ItemPrice, error)
	GetAllItems(ctx context.Context) ([]ItemPrice, error)
	SearchItemsPaginated(ctx context.Context, query string, params PaginationParams) (PaginatedResult[ItemPrice], error)
	GetAllItemsPaginated(ctx context.Context, params PaginationParams) (PaginatedResult[ItemPrice], error)
	SavePriceHistory(ctx context.Context, itemID int, price int, date time.Time) error
	GetPriceHistory(ctx context.Context, itemID int, days int) ([]PriceHistory, error)
}
