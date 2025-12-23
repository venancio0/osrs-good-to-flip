package osrs

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// OsrsWikiClient implements PriceProvider for OSRS Wiki API
type OsrsWikiClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewOsrsWikiClient creates a new OSRS Wiki client
func NewOsrsWikiClient() *OsrsWikiClient {
	return &OsrsWikiClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: "https://prices.runescape.wiki/api/v1/osrs",
	}
}

// FetchLatestPrices fetches the latest prices from OSRS Wiki API
// For MVP, returns mock data. Structure is ready for real API integration.
func (c *OsrsWikiClient) FetchLatestPrices(ctx context.Context) (map[int]int, error) {
	// TODO: Implement real API call when ready
	// For now, return mock data
	return c.getMockPrices(), nil
}

// getMockPrices returns mock price data for MVP
func (c *OsrsWikiClient) getMockPrices() map[int]int {
	// Mock prices for popular OSRS items
	// Format: itemID -> price (in GP)
	return map[int]int{
		1:  15000,  // Rune Scimitar
		2:  60000,  // Dragon Longsword
		3:  3000000, // Abyssal Whip
		4:  500000,  // Dragon Boots
		5:  120000,  // Dragon Platelegs
		6:  80000,   // Dragon Med Helm
		7:  2500000, // Bandos Chestplate
		8:  2000000, // Bandos Tassets
		9:  1500000, // Armadyl Helmet
		10: 5000000, // Armadyl Chestplate
		11: 4500000, // Armadyl Chainskirt
		12: 800000,  // Barrows Gloves
		13: 50000,   // Amulet of Glory
		14: 200000,  // Amulet of Fury
		15: 100000,  // Ring of Wealth
		16: 300000,  // Berserker Ring
		17: 150000,  // Warrior Ring
		18: 250000,  // Seers Ring
		19: 400000,  // Archers Ring
		20: 350000,  // Godsword
		21: 1200000, // Saradomin Sword
		22: 800000,  // Zamorakian Spear
		23: 600000,  // Abyssal Dagger
		24: 400000,  // Dragon Dagger
		25: 200000,  // Rune Full Helm
		26: 180000,  // Rune Platebody
		27: 160000,  // Rune Platelegs
		28: 140000,  // Rune Boots
		29: 100000,  // Rune Gloves
		30: 50000,   // Rune Kiteshield
	}
}

// fetchFromAPI is a placeholder for real API integration
func (c *OsrsWikiClient) fetchFromAPI(ctx context.Context) (map[int]int, error) {
	url := fmt.Sprintf("%s/latest", c.baseURL)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "OSRS-Good-to-Flip/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// TODO: Parse JSON response
	// The API returns data in format: {"data": {itemID: {"high": price, "low": price}}}
	// For now, return empty map
	return make(map[int]int), nil
}

