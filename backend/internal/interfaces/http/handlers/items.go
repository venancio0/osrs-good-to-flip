package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/gabv/osrs-good-to-flip/internal/application"
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
	query := r.URL.Query().Get("q")

	items, err := h.searchItemsUseCase.Execute(r.Context(), query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to fetch items")
		return
	}

	respondWithJSON(w, http.StatusOK, items)
}

// GetItemByID handles GET /items/{id}
func (h *ItemsHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	item, err := h.getItemUseCase.Execute(r.Context(), id)
	if err != nil {
		if err.Error() == "item not found" {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, item)
}

// GetPriceHistory handles GET /items/{id}/history
func (h *ItemsHandler) GetPriceHistory(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	days := r.URL.Query().Get("days") // Optional query param for number of days

	history, err := h.getPriceHistoryUseCase.Execute(r.Context(), id, days)
	if err != nil {
		if err.Error() == "item not found" {
			respondWithError(w, http.StatusNotFound, "Item not found")
			return
		}
		respondWithError(w, http.StatusBadRequest, err.Error())
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

