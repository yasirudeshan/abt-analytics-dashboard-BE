package api

import (
	"abt-analytics-dashboard/internal/config"
	"abt-analytics-dashboard/internal/processor"
	"compress/gzip"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	server    *http.Server
	processor *processor.Processor
	config    *config.Config
}

// NewServer creates a new HTTP server instance
func NewServer(proc *processor.Processor, cfg *config.Config) *Server {
	s := &Server{
		processor: proc,
		config:    cfg,
	}

	router := s.setupRoutes()

	s.server = &http.Server{
		Addr:         cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return s
}

// setupRoutes configures all API routes
func (s *Server) setupRoutes() http.Handler {
	router := mux.NewRouter()

	// Add middleware
	router.Use(s.loggingMiddleware)
	router.Use(s.corsMiddleware)
	router.Use(s.compressionMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", s.healthCheck).Methods("GET")
	api.HandleFunc("/revenue-by-country", s.getCountryRevenues).Methods("GET")
	api.HandleFunc("/top-products", s.getTopProducts).Methods("GET")
	api.HandleFunc("/sales-by-month", s.getMonthlySales).Methods("GET")
	api.HandleFunc("/top-regions", s.getTopRegions).Methods("GET")
	api.HandleFunc("/dashboard", s.getDashboardData).Methods("GET")
	api.HandleFunc("/summary", s.getDashboardSummary).Methods("GET")

	// Static route for basic info
	router.HandleFunc("/", s.rootHandler).Methods("GET")

	return router
}

// Middleware functions
func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		log.Printf(
			"%s %s %s %v",
			r.Method,
			r.RequestURI,
			r.RemoteAddr,
			time.Since(start),
		)
	})
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) compressionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if client supports gzip
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		// Set headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Vary", "Accept-Encoding")

		// Create gzip writer
		gz := gzip.NewWriter(w)
		defer gz.Close()

		// Wrap response writer
		gzw := &gzipResponseWriter{Writer: gz, ResponseWriter: w}
		next.ServeHTTP(gzw, r)
	})
}

type gzipResponseWriter struct {
	*gzip.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func (w *gzipResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *gzipResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

// Handler functions
func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"service": "ABT Analytics Dashboard API",
		"version": "1.0.0",
		"status":  "running",
		"endpoints": map[string]interface{}{
			"health":             "/api/health",
			"country_revenues":   "/api/revenue-by-country?page=0&limit=1000",
			"top_products":       "/api/top-products",
			"monthly_sales":      "/api/sales-by-month",
			"top_regions":        "/api/top-regions",
			"complete_dashboard": "/api/dashboard",
			"summary":            "/api/summary",
		},
		"pagination": map[string]interface{}{
			"description": "Use ?page=N&limit=M for pagination (max limit: 10000)",
			"example":     "/api/revenue-by-country?page=0&limit=100",
			"response_format": map[string]interface{}{
				"data":       "Array of country revenue records",
				"total":      "Total number of records in entire dataset",
				"page":       "Current page number (0-based)",
				"limit":      "Number of records per page",
				"totalPages": "Total number of pages",
				"hasMore":    "Boolean indicating if there are more pages",
			},
		},
		"summary_endpoint": map[string]interface{}{
			"description": "Returns complete totals and statistics from the entire dataset",
			"url":         "/api/summary",
			"fields": map[string]string{
				"total_revenue":      "Sum of all revenue across all records",
				"total_transactions": "Total number of transactions",
				"total_countries":    "Number of unique countries",
				"total_products":     "Number of unique products",
				"total_records":      "Total number of records processed",
				"unique_products":    "Number of products in top products list",
				"unique_regions":     "Number of regions in top regions list",
				"months_covered":     "Number of months in monthly sales data",
			},
		},
	}
	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	dashboardData := s.processor.GetDashboardData()
	response := map[string]interface{}{
		"status":              "healthy",
		"timestamp":           time.Now(),
		"last_data_update":    dashboardData.LastUpdated,
		"processing_duration": dashboardData.ProcessingDuration.String(),
		"record_count":        dashboardData.RecordCount,
	}
	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) getCountryRevenues(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters for pagination
	page := 0
	limit := 1000 // Default limit to prevent huge responses

	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p >= 0 {
			page = p
		}
	}

	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 10000 {
			limit = l
		}
	}

	// Get paginated data
	data, totalCount, totalPages := s.processor.GetCountryRevenuesPaginated(page, limit)

	response := map[string]interface{}{
		"data":       data,
		"total":      totalCount,
		"page":       page,
		"limit":      limit,
		"totalPages": totalPages,
		"hasMore":    page+1 < totalPages,
		"meta": map[string]interface{}{
			"description": "Country-level revenue data sorted by total revenue (descending)",
			"updated_at":  s.processor.GetDashboardData().LastUpdated,
		},
	}
	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) getTopProducts(w http.ResponseWriter, r *http.Request) {
	data := s.processor.GetTopProducts()
	response := map[string]interface{}{
		"data":  data,
		"count": len(data),
		"meta": map[string]interface{}{
			"description": "Top 20 most frequently purchased products with current stock",
			"updated_at":  s.processor.GetDashboardData().LastUpdated,
		},
	}
	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) getMonthlySales(w http.ResponseWriter, r *http.Request) {
	data := s.processor.GetMonthlySales()
	response := map[string]interface{}{
		"data":  data,
		"count": len(data),
		"meta": map[string]interface{}{
			"description": "Monthly sales volume data highlighting peak sales periods",
			"updated_at":  s.processor.GetDashboardData().LastUpdated,
		},
	}
	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) getTopRegions(w http.ResponseWriter, r *http.Request) {
	data := s.processor.GetTopRegions()
	response := map[string]interface{}{
		"data":  data,
		"count": len(data),
		"meta": map[string]interface{}{
			"description": "Top 30 regions by total revenue and items sold",
			"updated_at":  s.processor.GetDashboardData().LastUpdated,
		},
	}
	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) getDashboardData(w http.ResponseWriter, r *http.Request) {
	data := s.processor.GetDashboardData()
	response := map[string]interface{}{
		"data": data,
		"meta": map[string]interface{}{
			"description": "Complete dashboard data including all metrics",
			"updated_at":  data.LastUpdated,
		},
	}
	s.writeJSONResponse(w, http.StatusOK, response)
}

func (s *Server) getDashboardSummary(w http.ResponseWriter, r *http.Request) {
	// Check if data is loaded
	dashboardData := s.processor.GetDashboardData()
	if dashboardData.RecordCount == 0 {
		s.writeErrorResponse(w, http.StatusServiceUnavailable, "Data not loaded yet")
		return
	}

	// Get summary statistics
	summary := s.processor.GetDashboardSummary()

	response := map[string]interface{}{
		"data": summary,
		"meta": map[string]interface{}{
			"description": "Complete dashboard summary with totals and statistics",
			"updated_at":  dashboardData.LastUpdated,
		},
	}

	s.writeJSONResponse(w, http.StatusOK, response)
	log.Printf("Served dashboard summary: %+v", summary)
}

// Helper functions
func (s *Server) writeJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"error":     true,
		"message":   message,
		"timestamp": time.Now(),
	}
	s.writeJSONResponse(w, statusCode, response)
}

// Server lifecycle methods
func (s *Server) ListenAndServe() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
