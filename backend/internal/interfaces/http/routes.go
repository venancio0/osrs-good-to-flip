package http

import (
	"net/http"
	"os"

	"github.com/gabv/osrs-good-to-flip/backend/internal/interfaces/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// SetupRoutes configures all HTTP routes
func SetupRoutes(
	itemsHandler *handlers.ItemsHandler,
	healthHandler *handlers.HealthHandler,
) http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware for frontend
	allowedOrigins := []string{"http://localhost:3000", "http://localhost:3001"}

	// Add Vercel preview/production URLs if provided
	if vercelURL := os.Getenv("VERCEL_URL"); vercelURL != "" {
		allowedOrigins = append(allowedOrigins, "https://"+vercelURL)
	}
	if vercelProdURL := os.Getenv("VERCEL_PRODUCTION_URL"); vercelProdURL != "" {
		allowedOrigins = append(allowedOrigins, "https://"+vercelProdURL)
	}

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", healthHandler.Check)
	r.Route("/items", func(r chi.Router) {
		r.Get("/", itemsHandler.GetItems)
		r.Get("/{id}", itemsHandler.GetItemByID)
		r.Get("/{id}/history", itemsHandler.GetPriceHistory)
	})

	return r
}
