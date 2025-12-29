package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabv/osrs-good-to-flip/backend/internal/application"
	"github.com/gabv/osrs-good-to-flip/backend/internal/domain"
	"github.com/go-chi/chi/v5"
)

// ItemsHandler handles item-related HTTP requests
type ItemsHandler struct {
	getItemUseCase         *application.GetItemUseCase
	searchItemsUseCase     *application.SearchItemsUseCase
	getPriceHistoryUseCase *application.GetPriceHistoryUseCase
}

// NewItemsHandler creates a new ItemsHandler
func NewItemsHandler(
	getItemUseCase *application.GetItemUseCase,
	searchItemsUseCase *application.SearchItemsUseCase,
	getPriceHistoryUseCase *application.GetPriceHistoryUseCase,
) *ItemsHandler {
	return &ItemsHandler{
		getItemUseCase:         getItemUseCase,
		searchItemsUseCase:     searchItemsUseCase,
		getPriceHistoryUseCase: getPriceHistoryUseCase,
	}
}

// GetItems handles GET /items
func (h *ItemsHandler) GetItems(w http.ResponseWriter, r *http.Request) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	query := r.URL.Query().Get("q")

	// Validate query
	if err := validateQuery(query); err != nil {
		respondWithError(w, http.StatusBadRequest, getSafeErrorMessage(err))
		return
	}

	// Parse and validate pagination params
	page := parseIntQuery(r, "page", 1)
	limit := parseIntQuery(r, "limit", 20)
	
	validatedPage, validatedLimit, err := validatePaginationParams(page, limit)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, getSafeErrorMessage(err))
		return
	}
	
	params := domain.NewPaginationParams(validatedPage, validatedLimit)

	// Use paginated version
	result, err := h.searchItemsUseCase.ExecutePaginated(ctx, query, params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, getSafeErrorMessage(err))
		return
	}

	respondWithJSON(w, http.StatusOK, result)
}

// parseIntQuery parses an integer query parameter with a default value
func parseIntQuery(r *http.Request, key string, defaultValue int) int {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}

	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		return defaultValue
	}
	return result
}

// GetItemByID handles GET /items/{id}
func (h *ItemsHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")

	// Validate item ID
	id, err := validateItemID(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, getSafeErrorMessage(err))
		return
	}

	item, err := h.getItemUseCase.Execute(ctx, fmt.Sprintf("%d", id))
	if err != nil {
		if err.Error() == "item not found" {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, getSafeErrorMessage(err))
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// GetPriceHistory handles GET /items/{id}/history
func (h *ItemsHandler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	idStr := chi.URLParam(r, "id")
	daysStr := r.URL.Query().Get("days")

	// Validate item ID
	id, err := validateItemID(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, getSafeErrorMessage(err))
		return
	}

	// Validate days parameter
	days, err := validateDays(daysStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, getSafeErrorMessage(err))
		return
	}

	history, err := h.getPriceHistoryUseCase.Execute(ctx, fmt.Sprintf("%d", id), fmt.Sprintf("%d", days))
	if err != nil {
		if err.Error() == "item not found" {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, getSafeErrorMessage(err))
		return
	}

	respondWithJSON(w, http.StatusOK, history)
}

// respondWithJSON sends a JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// respondWithError sends an error JSON response
func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	respondWithJSON(w, statusCode, map[string]string{"error": message})
}

// getSafeErrorMessage returns a safe error message that doesn't expose internal details in production
func getSafeErrorMessage(err error) string {
	if err == nil {
		return "An error occurred"
	}
	
	errMsg := err.Error()
	
	// In production, return generic messages for internal errors
	if isProduction() {
		// Only expose validation errors (they're safe)
		if strings.Contains(errMsg, "invalid") || 
		   strings.Contains(errMsg, "too long") || 
		   strings.Contains(errMsg, "out of valid range") ||
		   strings.Contains(errMsg, "required") ||
		   strings.Contains(errMsg, "must be between") {
			return errMsg
		}
		
		// For other errors, return generic message
		return "An error occurred. Please try again later."
	}
	
	// In development, return full error message
	return errMsg
}
