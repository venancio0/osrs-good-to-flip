package domain

import (
	"context"
	"time"
)

// PriceProvider defines the interface for fetching prices from external sources
type PriceProvider interface {
	FetchLatestPrices(ctx context.Context) (map[int]int, error)
}

// ItemRepository defines the interface for item data operations
type ItemRepository interface {
	SavePrices(ctx context.Context, prices []ItemPrice) error
	GetItemByID(ctx context.Context, id int) (*ItemPrice, error)
	SearchItems(ctx context.Context, query string) ([]ItemPrice, error)
	GetAllItems(ctx context.Context) ([]ItemPrice, error)
	SavePriceHistory(ctx context.Context, itemID int, price int, date time.Time) error
	GetPriceHistory(ctx context.Context, itemID int, days int) ([]PriceHistory, error)
}

