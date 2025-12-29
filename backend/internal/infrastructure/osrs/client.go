package osrs

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gabv/osrs-good-to-flip/backend/internal/domain"
)

// OsrsWikiClient implements PriceProvider for OSRS Wiki API.
type OsrsWikiClient struct {
	httpClient *http.Client
	baseURL    string
	userAgent  string

	// simple in-memory cache to respect the API and reduce calls
	cacheMu      sync.RWMutex
	cachedLatest map[int]domain.PriceSnapshot
	cachedAt     time.Time
	cacheTTL     time.Duration
}

// latestResponse mirrors the /latest payload.
type latestResponse struct {
	Data map[string]struct {
		High int `json:"high"`
		Low  int `json:"low"`
		// highTime, lowTime omitted for now
	} `json:"data"`
}

// NewOsrsWikiClient creates a new OSRS Wiki client with sane defaults.
func NewOsrsWikiClient() *OsrsWikiClient {
	timeout := 10 * time.Second
	if v := os.Getenv("OSRS_WIKI_TIMEOUT_MS"); v != "" {
		if ms, err := strconv.Atoi(v); err == nil && ms > 0 {
			timeout = time.Duration(ms) * time.Millisecond
		}
	}

	baseURL := os.Getenv("OSRS_WIKI_BASE_URL")
	if baseURL == "" {
		baseURL = "https://prices.runescape.wiki/api/v1/osrs"
	}

	userAgent := os.Getenv("OSRS_WIKI_USER_AGENT")
	if userAgent == "" {
		userAgent = "OSRS-GoodToFlip/1.0 (contact: your-email@example.com)"
	}

	cacheTTL := 60 * time.Second
	if v := os.Getenv("OSRS_WIKI_CACHE_TTL_SEC"); v != "" {
		if s, err := strconv.Atoi(v); err == nil && s > 0 {
			cacheTTL = time.Duration(s) * time.Second
		}
	}

	return &OsrsWikiClient{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
		userAgent:  userAgent,
		cacheTTL:   cacheTTL,
	}
}

// FetchLatestPrices fetches latest high/low prices with a small TTL cache.
func (c *OsrsWikiClient) FetchLatestPrices(ctx context.Context) (map[int]domain.PriceSnapshot, error) {
	// serve from cache when fresh
	if data, ok := c.getCached(); ok {
		return data, nil
	}

	url := fmt.Sprintf("%s/latest", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("osrs wiki: status %d", resp.StatusCode)
	}

	var payload latestResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}

	result := make(map[int]domain.PriceSnapshot, len(payload.Data))
	for idStr, v := range payload.Data {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		result[id] = domain.PriceSnapshot{High: v.High, Low: v.Low}
	}

	c.setCache(result)
	return result, nil
}

func (c *OsrsWikiClient) getCached() (map[int]domain.PriceSnapshot, bool) {
	c.cacheMu.RLock()
	defer c.cacheMu.RUnlock()
	if c.cachedLatest == nil {
		return nil, false
	}
	if time.Since(c.cachedAt) > c.cacheTTL {
		return nil, false
	}
	// return a copy to avoid external mutation
	out := make(map[int]domain.PriceSnapshot, len(c.cachedLatest))
	for k, v := range c.cachedLatest {
		out[k] = v
	}
	return out, true
}

func (c *OsrsWikiClient) setCache(data map[int]domain.PriceSnapshot) {
	c.cacheMu.Lock()
	defer c.cacheMu.Unlock()
	c.cachedLatest = data
	c.cachedAt = time.Now()
}
