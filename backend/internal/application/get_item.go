package application

import (
	"context"
	"errors"
	"strconv"

	"github.com/gabv/osrs-good-to-flip/internal/domain"
)

// GetItemUseCase handles retrieving an item by ID
type GetItemUseCase struct {
	repo domain.ItemRepository
}

// NewGetItemUseCase creates a new GetItemUseCase
func NewGetItemUseCase(repo domain.ItemRepository) *GetItemUseCase {
	return &GetItemUseCase{repo: repo}
}

// Execute retrieves an item by its ID
func (uc *GetItemUseCase) Execute(ctx context.Context, idStr string) (*domain.ItemPrice, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, errors.New("invalid item ID")
	}

	item, err := uc.repo.GetItemByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return item, nil
}

