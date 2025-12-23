package application

import (
	"context"
	"time"

	"github.com/gabv/osrs-good-to-flip/internal/domain"
)

// UpdatePricesUseCase handles updating item prices
type UpdatePricesUseCase struct {
	provider domain.PriceProvider
	repo     domain.ItemRepository
}

// NewUpdatePricesUseCase creates a new UpdatePricesUseCase
func NewUpdatePricesUseCase(provider domain.PriceProvider, repo domain.ItemRepository) *UpdatePricesUseCase {
	return &UpdatePricesUseCase{
		provider: provider,
		repo:     repo,
	}
}

// Execute fetches latest prices and updates the repository
func (uc *UpdatePricesUseCase) Execute(ctx context.Context) error {
	// Fetch latest prices from provider
	prices, err := uc.provider.FetchLatestPrices(ctx)
	if err != nil {
		return err
	}

	// Convert map to ItemPrice slice
	items := make([]domain.ItemPrice, 0, len(prices))
	now := time.Now()

	for itemID, price := range prices {
		// Get existing item to preserve name and averages
		existing, _ := uc.repo.GetItemByID(ctx, itemID)
		
		item := domain.ItemPrice{
			ItemID:    itemID,
			Price:     price,
			UpdatedAt: now,
		}

		if existing != nil {
			item.Name = existing.Name
			item.Avg24h = existing.Avg24h
			item.Avg7d = existing.Avg7d
			// Calculate trend based on current price vs 24h average
			item.Trend = domain.CalculateTrend(price, existing.Avg24h)
		} else {
			// For new items, set defaults
			item.Avg24h = price
			item.Avg7d = price
			item.Trend = domain.TrendFlat
		}

		items = append(items, item)
	}

	return uc.repo.SavePrices(ctx, items)
}

