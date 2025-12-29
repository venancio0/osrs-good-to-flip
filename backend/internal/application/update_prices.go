package application

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gabv/osrs-good-to-flip/backend/internal/domain"
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
	snapshots, err := uc.provider.FetchLatestPrices(ctx)
	if err != nil {
		return err
	}

	// Fetch item names mapping (only if we have new items)
	itemNames, err := uc.provider.FetchItemNames(ctx)
	if err != nil {
		log.Printf("Warning: Failed to fetch item names: %v", err)
		// Continue without names - will use existing names or fallback
		itemNames = make(map[int]string)
	}

	// Convert map to ItemPrice slice
	items := make([]domain.ItemPrice, 0, len(snapshots))
	now := time.Now()

	for itemID, snap := range snapshots {
		// choose a representative price; we use High as current price
		price := snap.High
		if price == 0 {
			price = snap.Low
		}

		// Skip entries with no price info
		if price == 0 {
			continue
		}

		// Get existing item to preserve name and averages
		existing, _ := uc.repo.GetItemByID(ctx, itemID)

		item := domain.ItemPrice{
			ItemID:    itemID,
			Price:     price,
			High:      snap.High,
			Low:       snap.Low,
			Volume:    snap.Volume,
			UpdatedAt: now,
		}

		if existing != nil {
			// Preserve existing name and averages
			item.Name = existing.Name
			item.Avg24h = existing.Avg24h
			item.Avg7d = existing.Avg7d
			// Calculate trend based on current price vs 24h average
			item.Trend = domain.CalculateTrend(price, existing.Avg24h)
		} else {
			// For new items, try to get name from mapping
			if name, ok := itemNames[itemID]; ok && name != "" {
				item.Name = name
			} else {
				// Fallback: use placeholder name
				item.Name = fmt.Sprintf("Item %d", itemID)
			}
			// Set defaults for new items
			item.Avg24h = price
			item.Avg7d = price
			item.Trend = domain.TrendFlat
		}

		items = append(items, item)
	}

	return uc.repo.SavePrices(ctx, items)
}
