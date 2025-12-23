package domain

import "time"

// PriceHistory represents a historical price point for an item
type PriceHistory struct {
	ItemID    int       `json:"item_id"`
	Price     int       `json:"price"`
	Date      time.Time `json:"date"`
	CreatedAt time.Time `json:"created_at"`
}

// PriceHistoryEntry represents a single entry in the price history
// Used for API responses with formatted date
type PriceHistoryEntry struct {
	Date  string `json:"date"`  // ISO 8601 format
	Price int    `json:"price"`
}

