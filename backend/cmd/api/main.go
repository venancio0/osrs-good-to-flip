package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gabv/osrs-good-to-flip/backend/internal/application"
	osrsclient "github.com/gabv/osrs-good-to-flip/backend/internal/infrastructure/osrs"
	"github.com/gabv/osrs-good-to-flip/backend/internal/infrastructure/repository"
	"github.com/gabv/osrs-good-to-flip/backend/internal/infrastructure/worker"
	httpInterface "github.com/gabv/osrs-good-to-flip/backend/internal/interfaces/http"
	"github.com/gabv/osrs-good-to-flip/backend/internal/interfaces/http/handlers"
)

func main() {
	// Initialize infrastructure
	repo := repository.NewInMemoryRepository()
	osrsClient := osrsclient.NewOsrsWikiClient()

	// Initialize use cases
	getItemUseCase := application.NewGetItemUseCase(repo)
	searchItemsUseCase := application.NewSearchItemsUseCase(repo)
	updatePricesUseCase := application.NewUpdatePricesUseCase(osrsClient, repo)
	getPriceHistoryUseCase := application.NewGetPriceHistoryUseCase(repo)

	// Initialize handlers
	itemsHandler := handlers.NewItemsHandler(getItemUseCase, searchItemsUseCase, getPriceHistoryUseCase)
	healthHandler := handlers.NewHealthHandler()

	// Setup routes
	router := httpInterface.SetupRoutes(itemsHandler, healthHandler)

	// Create HTTP server
	port := getPort()
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Run initial price update
	ctx := context.Background()
	if err := updatePricesUseCase.Execute(ctx); err != nil {
		log.Printf("Warning: Failed to update prices: %v", err)
	}

	// Start price updater worker (updates every 5 minutes)
	updateInterval := 5 * time.Minute
	if envInterval := os.Getenv("PRICE_UPDATE_INTERVAL_MIN"); envInterval != "" {
		if minutes, err := strconv.Atoi(envInterval); err == nil && minutes > 0 {
			updateInterval = time.Duration(minutes) * time.Minute
		}
	}

	priceWorker := worker.NewPriceUpdaterWorker(updatePricesUseCase, updateInterval)
	priceWorker.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Stop price worker first
	priceWorker.Stop()

	// Shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
