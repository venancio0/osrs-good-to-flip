package application

import (
	"context"

	"github.com/gabv/osrs-good-to-flip/backend/internal/domain"
)

// SearchItemsUseCase handles searching items by name
type SearchItemsUseCase struct {
	repo domain.ItemRepository
}

// NewSearchItemsUseCase creates a new SearchItemsUseCase
func NewSearchItemsUseCase(repo domain.ItemRepository) *SearchItemsUseCase {
	return &SearchItemsUseCase{repo: repo}
}

// Execute searches for items matching the query
func (uc *SearchItemsUseCase) Execute(ctx context.Context, query string) ([]domain.ItemPrice, error) {
	if query == "" {
		// If no query, return all items
		return uc.repo.GetAllItems(ctx)
	}

	items, err := uc.repo.SearchItems(ctx, query)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// ExecutePaginated searches for items with pagination
func (uc *SearchItemsUseCase) ExecutePaginated(ctx context.Context, query string, params domain.PaginationParams) (domain.PaginatedResult[domain.ItemPrice], error) {
	if query == "" {
		return uc.repo.GetAllItemsPaginated(ctx, params)
	}

	return uc.repo.SearchItemsPaginated(ctx, query, params)
}
