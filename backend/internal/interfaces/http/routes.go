package http

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gabv/osrs-good-to-flip/backend/internal/interfaces/http/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
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
	
	// Security headers middleware
	r.Use(securityHeadersMiddleware)
	
	// Rate limiting middleware
	r.Use(rateLimitMiddleware())

	// CORS middleware for frontend
	allowedOrigins := []string{
		"http://localhost:3000",
		"http://localhost:3001",
		"https://osrs-good-to-flip.vercel.app", // Production Vercel URL
	}

	// Add custom allowed origins from environment variable
	if customOrigins := os.Getenv("ALLOWED_ORIGINS"); customOrigins != "" {
		// Support comma-separated list
		origins := strings.Split(customOrigins, ",")
		for _, origin := range origins {
			origin = strings.TrimSpace(origin)
			if origin != "" {
				allowedOrigins = append(allowedOrigins, origin)
			}
		}
	}

	// Custom CORS handler that allows *.vercel.app domains
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			// Check if origin is in allowed list
			isAllowed := false
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					isAllowed = true
					break
				}
			}

			// Also allow any *.vercel.app domain (preview deployments)
			if !isAllowed && origin != "" && strings.HasSuffix(origin, ".vercel.app") {
				isAllowed = true
			}

			if isAllowed {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
				w.Header().Set("Access-Control-Expose-Headers", "Link")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "300")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Routes
	r.Get("/health", healthHandler.Check)
	r.Route("/items", func(r chi.Router) {
		r.Get("/", itemsHandler.GetItems)
		r.Get("/{id}", itemsHandler.GetItemByID)
		r.Get("/{id}/history", itemsHandler.GetPriceHistory)
	})

	return r
}

// securityHeadersMiddleware adds security headers to all responses
func securityHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// HSTS - Force HTTPS
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		
		// Prevent MIME type sniffing
		w.Header().Set("X-Content-Type-Options", "nosniff")
		
		// Prevent clickjacking
		w.Header().Set("X-Frame-Options", "DENY")
		
		// XSS Protection
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		
		// Referrer Policy
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy (relaxed for API)
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		
		// Permissions Policy
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		next.ServeHTTP(w, r)
	})
}

// rateLimitMiddleware configures rate limiting based on environment variables
func rateLimitMiddleware() func(http.Handler) http.Handler {
	// Default: 100 requests per minute per IP
	requestsPerMinute := 100
	if envRate := os.Getenv("RATE_LIMIT_REQUESTS_PER_MINUTE"); envRate != "" {
		if rate, err := strconv.Atoi(envRate); err == nil && rate > 0 {
			requestsPerMinute = rate
		}
	}
	
	return httprate.Limit(
		requestsPerMinute,
		1*time.Minute,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint),
	)
}
