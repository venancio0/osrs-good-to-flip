package worker

import (
	"context"
	"log"
	"time"

	"github.com/gabv/osrs-good-to-flip/internal/application"
)

// PriceUpdaterWorker handles periodic price updates
type PriceUpdaterWorker struct {
	updateUseCase *application.UpdatePricesUseCase
	interval      time.Duration
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewPriceUpdaterWorker creates a new price updater worker
func NewPriceUpdaterWorker(updateUseCase *application.UpdatePricesUseCase, interval time.Duration) *PriceUpdaterWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &PriceUpdaterWorker{
		updateUseCase: updateUseCase,
		interval:      interval,
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Start begins the periodic price update worker
func (w *PriceUpdaterWorker) Start() {
	log.Printf("Price updater worker started with interval: %v", w.interval)

	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		// Run immediately on start
		w.updatePrices()

		for {
			select {
			case <-ticker.C:
				w.updatePrices()
			case <-w.ctx.Done():
				log.Println("Price updater worker stopped")
				return
			}
		}
	}()
}

// Stop stops the worker
func (w *PriceUpdaterWorker) Stop() {
	log.Println("Stopping price updater worker...")
	w.cancel()
}

// updatePrices performs a single price update
func (w *PriceUpdaterWorker) updatePrices() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	log.Println("Updating prices from OSRS Wiki API...")
	if err := w.updateUseCase.Execute(ctx); err != nil {
		log.Printf("Error updating prices: %v", err)
		return
	}
	log.Println("Prices updated successfully")
}
