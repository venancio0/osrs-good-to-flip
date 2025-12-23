package domain

import "time"

// ItemPrice represents the price information for an OSRS item
type ItemPrice struct {
	ItemID    int       `json:"item_id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`      // Current average price
	High      int       `json:"high"`       // Best buy price
	Low       int       `json:"low"`        // Best sell price
	Volume    int       `json:"volume"`     // Trading volume
	Avg24h    int       `json:"avg_24h"`
	Avg7d     int       `json:"avg_7d"`
	Trend     TrendType `json:"trend"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TrendType represents the price trend direction
type TrendType string

const (
	TrendUp   TrendType = "UP"
	TrendDown TrendType = "DOWN"
	TrendFlat TrendType = "FLAT"
)

// CalculateTrend calculates the trend based on current price vs average
// Simple logic: compare current price with 24h average
func CalculateTrend(currentPrice, avg24h int) TrendType {
	if avg24h == 0 {
		return TrendFlat
	}

	diff := float64(currentPrice-avg24h) / float64(avg24h) * 100

	// Consider trend up if price is 2% higher than average
	if diff > 2.0 {
		return TrendUp
	}
	// Consider trend down if price is 2% lower than average
	if diff < -2.0 {
		return TrendDown
	}
	return TrendFlat
}

// CalculateMargin calculates the profit margin percentage
func CalculateMargin(buyPrice, sellPrice int) float64 {
	if buyPrice == 0 {
		return 0
	}
	return float64(sellPrice-buyPrice) / float64(buyPrice) * 100
}

// CalculateGETax calculates the Grand Exchange tax (1% of sell price)
func CalculateGETax(sellPrice int) int {
	return int(float64(sellPrice) * 0.01)
}

// CalculateExpectedProfit calculates expected profit after GE tax
func CalculateExpectedProfit(buyPrice, sellPrice int) int {
	tax := CalculateGETax(sellPrice)
	return sellPrice - buyPrice - tax
}

