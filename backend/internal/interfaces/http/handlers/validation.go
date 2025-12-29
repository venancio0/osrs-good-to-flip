package handlers

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	maxQueryLength    = 100
	maxItemID         = 10000000
	minItemID         = 1
	maxDays           = 30
	minDays           = 1
	maxPage           = 10000
	maxLimit          = 100
)

// validateItemID validates and parses an item ID from URL parameter
func validateItemID(idStr string) (int, error) {
	if idStr == "" {
		return 0, fmt.Errorf("item ID is required")
	}
	
	// Limit length to prevent DoS
	if len(idStr) > 10 {
		return 0, fmt.Errorf("item ID too long")
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid item ID format")
	}
	
	if id < minItemID || id > maxItemID {
		return 0, fmt.Errorf("item ID out of valid range")
	}
	
	return id, nil
}

// validateQuery validates search query string
func validateQuery(query string) error {
	if len(query) > maxQueryLength {
		return fmt.Errorf("query string too long (max %d characters)", maxQueryLength)
	}
	
	// Check for potentially dangerous characters
	if strings.ContainsAny(query, "<>\"'&") {
		return fmt.Errorf("query contains invalid characters")
	}
	
	return nil
}

// validatePaginationParams validates and normalizes pagination parameters
func validatePaginationParams(page, limit int) (int, int, error) {
	if page < 1 {
		page = 1
	}
	if page > maxPage {
		return 0, 0, fmt.Errorf("page number too large (max %d)", maxPage)
	}
	
	if limit < 1 {
		limit = 20
	}
	if limit > maxLimit {
		return 0, 0, fmt.Errorf("limit too large (max %d)", maxLimit)
	}
	
	return page, limit, nil
}

// validateDays validates the days parameter for history queries
func validateDays(daysStr string) (int, error) {
	if daysStr == "" {
		return 7, nil // Default
	}
	
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return 0, fmt.Errorf("invalid days format")
	}
	
	if days < minDays || days > maxDays {
		return 0, fmt.Errorf("days must be between %d and %d", minDays, maxDays)
	}
	
	return days, nil
}

// isProduction checks if the application is running in production mode
func isProduction() bool {
	env := os.Getenv("ENV")
	return env == "production" || env == "prod"
}

