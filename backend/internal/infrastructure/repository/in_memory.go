package repository

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gabv/osrs-good-to-flip/backend/internal/domain"
)

// InMemoryRepository implements ItemRepository using in-memory storage
type InMemoryRepository struct {
	mu      sync.RWMutex
	items   map[int]*domain.ItemPrice
	history map[int][]domain.PriceHistory // itemID -> []PriceHistory
}

// NewInMemoryRepository creates a new in-memory repository with mock data
func NewInMemoryRepository() *InMemoryRepository {
	repo := &InMemoryRepository{
		items:   make(map[int]*domain.ItemPrice),
		history: make(map[int][]domain.PriceHistory),
	}

	// Initialize with mock data
	repo.initializeMockData()
	repo.initializeMockHistory()

	return repo
}

// SavePrices saves or updates item prices
func (r *InMemoryRepository) SavePrices(ctx context.Context, prices []domain.ItemPrice) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, price := range prices {
		// Create a copy to avoid pointer issues
		item := price
		r.items[price.ItemID] = &item
	}

	return nil
}

// GetItemByID retrieves an item by its ID
func (r *InMemoryRepository) GetItemByID(ctx context.Context, id int) (*domain.ItemPrice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, errors.New("item not found")
	}

	// Return a copy to avoid race conditions
	itemCopy := *item
	return &itemCopy, nil
}

// SearchItems searches for items by name (case-insensitive)
func (r *InMemoryRepository) SearchItems(ctx context.Context, query string) ([]domain.ItemPrice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	queryLower := strings.ToLower(query)
	results := make([]domain.ItemPrice, 0)

	for _, item := range r.items {
		if strings.Contains(strings.ToLower(item.Name), queryLower) {
			// Create a copy
			itemCopy := *item
			results = append(results, itemCopy)
		}
	}

	return results, nil
}

// GetAllItems returns all items in the repository
func (r *InMemoryRepository) GetAllItems(ctx context.Context) ([]domain.ItemPrice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make([]domain.ItemPrice, 0, len(r.items))
	for _, item := range r.items {
		itemCopy := *item
		results = append(results, itemCopy)
	}

	return results, nil
}

// initializeMockData populates the repository with mock OSRS items
func (r *InMemoryRepository) initializeMockData() {
	now := time.Now()

	// Helper function to generate realistic high/low prices and volume
	generateItem := func(itemID int, name string, price, avg24h, avg7d int, trend domain.TrendType) domain.ItemPrice {
		// High is typically 1-3% above average, Low is 1-3% below
		high := int(float64(price) * (1.0 + 0.01*float64((itemID%3)+1)))
		low := int(float64(price) * (1.0 - 0.01*float64((itemID%3)+1)))
		// Volume varies based on price (cheaper items trade more)
		volume := 0
		if price < 100000 {
			volume = 5000 + (itemID%1000)*10
		} else if price < 1000000 {
			volume = 500 + (itemID%100)*5
		} else {
			volume = 50 + (itemID%20)*2
		}
		return domain.ItemPrice{
			ItemID: itemID, Name: name, Price: price, High: high, Low: low,
			Volume: volume, Avg24h: avg24h, Avg7d: avg7d, Trend: trend, UpdatedAt: now,
		}
	}

	mockItems := []domain.ItemPrice{
		generateItem(1, "Rune Scimitar", 15000, 14800, 15000, domain.TrendUp),
		generateItem(2, "Dragon Longsword", 60000, 61000, 60000, domain.TrendDown),
		generateItem(3, "Abyssal Whip", 3000000, 2950000, 3000000, domain.TrendUp),
		generateItem(4, "Dragon Boots", 500000, 500000, 500000, domain.TrendFlat),
		generateItem(5, "Dragon Platelegs", 120000, 125000, 120000, domain.TrendDown),
		generateItem(6, "Dragon Med Helm", 80000, 82000, 80000, domain.TrendDown),
		generateItem(7, "Bandos Chestplate", 2500000, 2480000, 2500000, domain.TrendUp),
		generateItem(8, "Bandos Tassets", 2000000, 2000000, 2000000, domain.TrendFlat),
		generateItem(9, "Armadyl Helmet", 1500000, 1520000, 1500000, domain.TrendDown),
		generateItem(10, "Armadyl Chestplate", 5000000, 4950000, 5000000, domain.TrendUp),
		generateItem(11, "Armadyl Chainskirt", 4500000, 4500000, 4500000, domain.TrendFlat),
		generateItem(12, "Barrows Gloves", 800000, 810000, 800000, domain.TrendDown),
		generateItem(13, "Amulet of Glory", 50000, 49000, 50000, domain.TrendUp),
		generateItem(14, "Amulet of Fury", 200000, 205000, 200000, domain.TrendDown),
		generateItem(15, "Ring of Wealth", 100000, 100000, 100000, domain.TrendFlat),
		generateItem(16, "Berserker Ring", 300000, 295000, 300000, domain.TrendUp),
		generateItem(17, "Warrior Ring", 150000, 152000, 150000, domain.TrendDown),
		generateItem(18, "Seers Ring", 250000, 250000, 250000, domain.TrendFlat),
		generateItem(19, "Archers Ring", 400000, 395000, 400000, domain.TrendUp),
		generateItem(20, "Godsword", 350000, 360000, 350000, domain.TrendDown),
		generateItem(21, "Saradomin Sword", 1200000, 1180000, 1200000, domain.TrendUp),
		generateItem(22, "Zamorakian Spear", 800000, 800000, 800000, domain.TrendFlat),
		generateItem(23, "Abyssal Dagger", 600000, 610000, 600000, domain.TrendDown),
		generateItem(24, "Dragon Dagger", 400000, 395000, 400000, domain.TrendUp),
		generateItem(25, "Rune Full Helm", 200000, 200000, 200000, domain.TrendFlat),
		generateItem(26, "Rune Platebody", 180000, 182000, 180000, domain.TrendDown),
		generateItem(27, "Rune Platelegs", 160000, 158000, 160000, domain.TrendUp),
		generateItem(28, "Rune Boots", 140000, 140000, 140000, domain.TrendFlat),
		generateItem(29, "Rune Gloves", 100000, 102000, 100000, domain.TrendDown),
		generateItem(30, "Rune Kiteshield", 50000, 49000, 50000, domain.TrendUp),
	}

	for i := range mockItems {
		item := mockItems[i]
		r.items[item.ItemID] = &item
	}
}

// SavePriceHistory saves a price history entry for an item
func (r *InMemoryRepository) SavePriceHistory(ctx context.Context, itemID int, price int, date time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entry := domain.PriceHistory{
		ItemID:    itemID,
		Price:     price,
		Date:      date,
		CreatedAt: time.Now(),
	}

	r.history[itemID] = append(r.history[itemID], entry)
	return nil
}

// GetPriceHistory retrieves price history for an item for the last N days
func (r *InMemoryRepository) GetPriceHistory(ctx context.Context, itemID int, days int) ([]domain.PriceHistory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	history, exists := r.history[itemID]
	if !exists {
		return []domain.PriceHistory{}, nil
	}

	// Filter by date range (last N days)
	cutoffDate := time.Now().AddDate(0, 0, -days)
	result := make([]domain.PriceHistory, 0)

	for _, entry := range history {
		if entry.Date.After(cutoffDate) || entry.Date.Equal(cutoffDate) {
			entryCopy := entry
			result = append(result, entryCopy)
		}
	}

	// Sort by date ascending
	for i := 0; i < len(result)-1; i++ {
		for j := i + 1; j < len(result); j++ {
			if result[i].Date.After(result[j].Date) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// initializeMockHistory generates mock price history for the last 7 days
func (r *InMemoryRepository) initializeMockHistory() {
	now := time.Now()

	// Generate history for each item
	for itemID, item := range r.items {
		basePrice := item.Price
		history := make([]domain.PriceHistory, 0)

		// Generate one entry per day for the last 7 days
		// Create a more realistic price progression
		for i := 6; i >= 0; i-- {
			date := now.AddDate(0, 0, -i)

			// Create variation: start from a base and trend towards current price
			// This creates a more realistic progression
			daysAgo := float64(6 - i)
			progress := daysAgo / 6.0 // 0.0 (7 days ago) to 1.0 (today)

			// Start price is slightly different from current (simulating past prices)
			startPrice := float64(basePrice) * 0.92
			endPrice := float64(basePrice)

			// Linear interpolation with small random variation
			price := startPrice + (endPrice-startPrice)*progress

			// Add small random variation (±2%)
			variation := 1.0 + (float64((itemID+i)%5)-2.0)*0.01
			price = price * variation

			entry := domain.PriceHistory{
				ItemID:    itemID,
				Price:     int(price),
				Date:      date,
				CreatedAt: now,
			}
			history = append(history, entry)
		}

		r.history[itemID] = history
	}
}
