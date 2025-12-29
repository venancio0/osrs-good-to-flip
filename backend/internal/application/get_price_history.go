package application

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/gabv/osrs-good-to-flip/internal/domain"
)

// GetPriceHistoryUseCase handles retrieving price history for an item
type GetPriceHistoryUseCase struct {
	repo domain.ItemRepository
}

// NewGetPriceHistoryUseCase creates a new GetPriceHistoryUseCase
func NewGetPriceHistoryUseCase(repo domain.ItemRepository) *GetPriceHistoryUseCase {
	return &GetPriceHistoryUseCase{repo: repo}
}

// Execute retrieves price history for an item
// days defaults to 7 if not provided or invalid
func (uc *GetPriceHistoryUseCase) Execute(ctx context.Context, idStr string, daysStr string) ([]domain.PriceHistoryEntry, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	// Verify item exists
	_, err = uc.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, errors.New("item not found")
	}

	// Parse days, default to 7
	days := 7
	if daysStr != "" {
		parsedDays, err := strconv.Atoi(daysStr)
		if err == nil && parsedDays > 0 && parsedDays <= 30 {
			days = parsedDays
		}
	}

	// Get history from repository
	history, err := uc.repo.GetPriceHistory(ctx, id, days)
	if err != nil {
		return nil, err
	}

	// Convert to API response format
	entries := make([]domain.PriceHistoryEntry, len(history))
	for i, h := range history {
		entries[i] = domain.PriceHistoryEntry{
			Date:  h.Date.Format(time.RFC3339),
			Price: h.Price,
		}
	}

	return entries, nil
}
