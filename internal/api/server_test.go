package api

import (
	"abt-analytics-dashboard/internal/config"
	"abt-analytics-dashboard/internal/processor"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()

	server := NewServer(proc, cfg)

	if server == nil {
		t.Fatal("Expected server to be created, got nil")
	}

	if server.processor != proc {
		t.Error("Expected processor to be set correctly")
	}

	if server.config != cfg {
		t.Error("Expected config to be set correctly")
	}

	if server.server == nil {
		t.Fatal("Expected HTTP server to be initialized")
	}

	if server.server.Addr != ":8080" {
		t.Errorf("Expected server address to be ':8080', got '%s'", server.server.Addr)
	}
}

func TestRootHandler(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.rootHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Check response content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response structure
	if response["service"] != "ABT Analytics Dashboard API" {
		t.Errorf("Expected service name 'ABT Analytics Dashboard API', got '%v'", response["service"])
	}
	if response["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%v'", response["version"])
	}
	if response["status"] != "running" {
		t.Errorf("Expected status 'running', got '%v'", response["status"])
	}

	// Check endpoints
	endpoints, ok := response["endpoints"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected endpoints to be a map")
	}

	expectedEndpoints := []string{"health", "country_revenues", "top_products", "monthly_sales", "top_regions", "complete_dashboard"}
	for _, endpoint := range expectedEndpoints {
		if _, exists := endpoints[endpoint]; !exists {
			t.Errorf("Expected endpoint '%s' to be present", endpoint)
		}
	}
}

func TestHealthCheck(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	proc.LoadSampleData() // Load sample data for health check
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.healthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response structure
	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", response["status"])
	}

	if _, exists := response["timestamp"]; !exists {
		t.Error("Expected timestamp to be present")
	}

	if _, exists := response["last_data_update"]; !exists {
		t.Error("Expected last_data_update to be present")
	}

	if _, exists := response["processing_duration"]; !exists {
		t.Error("Expected processing_duration to be present")
	}

	if _, exists := response["record_count"]; !exists {
		t.Error("Expected record_count to be present")
	}
}

func TestGetCountryRevenues(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	proc.LoadSampleData()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/api/revenue-by-country", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getCountryRevenues)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}

	if _, exists := response["count"]; !exists {
		t.Error("Expected count field to be present")
	}

	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Country-level revenue data sorted by total revenue (descending)" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}

	if _, exists := meta["updated_at"]; !exists {
		t.Error("Expected updated_at to be present in meta")
	}
}

func TestGetTopProducts(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	proc.LoadSampleData()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/api/top-products", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getTopProducts)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}

	if _, exists := response["count"]; !exists {
		t.Error("Expected count field to be present")
	}

	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Top 20 most frequently purchased products with current stock" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}
}

func TestGetMonthlySales(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	proc.LoadSampleData()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/api/sales-by-month", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getMonthlySales)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}

	if _, exists := response["count"]; !exists {
		t.Error("Expected count field to be present")
	}

	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Monthly sales volume data highlighting peak sales periods" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}
}

func TestGetTopRegions(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	proc.LoadSampleData()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/api/top-regions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getTopRegions)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}

	if _, exists := response["count"]; !exists {
		t.Error("Expected count field to be present")
	}

	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Top 30 regions by total revenue and items sold" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}
}

func TestGetDashboardData(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	proc.LoadSampleData()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/api/dashboard", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.getDashboardData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify response structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}

	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Complete dashboard data including all metrics" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}
}

func TestCorsMiddleware(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("OPTIONS", "/api/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := server.corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d for OPTIONS request, got %d", http.StatusOK, status)
	}

	// Check CORS headers
	corsOrigin := rr.Header().Get("Access-Control-Allow-Origin")
	if corsOrigin != "*" {
		t.Errorf("Expected Access-Control-Allow-Origin '*', got '%s'", corsOrigin)
	}

	corsMethods := rr.Header().Get("Access-Control-Allow-Methods")
	if corsMethods != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("Expected Access-Control-Allow-Methods to include all methods, got '%s'", corsMethods)
	}
}

func TestLoggingMiddleware(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	server := NewServer(proc, cfg)

	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := server.loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		time.Sleep(10 * time.Millisecond) // Simulate some processing time
	}))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}
}

func TestWriteJSONResponse(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	server := NewServer(proc, cfg)

	rr := httptest.NewRecorder()
	testData := map[string]string{"test": "data"}

	server.writeJSONResponse(rr, http.StatusCreated, testData)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Verify response body
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if response["test"] != "data" {
		t.Errorf("Expected response data to match, got '%v'", response)
	}
}

func TestWriteErrorResponse(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	server := NewServer(proc, cfg)

	rr := httptest.NewRecorder()

	server.writeErrorResponse(rr, http.StatusBadRequest, "Invalid request")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	// Parse response body
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	if response["error"] != true {
		t.Error("Expected error field to be true")
	}

	if response["message"] != "Invalid request" {
		t.Errorf("Expected message 'Invalid request', got '%v'", response["message"])
	}

	if _, exists := response["timestamp"]; !exists {
		t.Error("Expected timestamp to be present")
	}
}

func TestServerShutdown(t *testing.T) {
	cfg := &config.Config{Port: ":0"} // Use port 0 for testing
	proc := processor.New()
	server := NewServer(proc, cfg)

	// Test shutdown with context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		t.Errorf("Expected no error on shutdown, got %v", err)
	}
}

func TestSetupRoutes(t *testing.T) {
	cfg := &config.Config{Port: ":8080"}
	proc := processor.New()
	server := NewServer(proc, cfg)

	router := server.setupRoutes()
	if router == nil {
		t.Fatal("Expected router to be created")
	}

	// Test that routes are properly configured by making requests
	testRoutes := []string{
		"/",
		"/api/health",
		"/api/revenue-by-country",
		"/api/top-products",
		"/api/sales-by-month",
		"/api/top-regions",
		"/api/dashboard",
	}

	for _, route := range testRoutes {
		req, err := http.NewRequest("GET", route, nil)
		if err != nil {
			t.Fatalf("Failed to create request for %s: %v", route, err)
		}

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		// All routes should return some response (not 404)
		if rr.Code == http.StatusNotFound {
			t.Errorf("Route %s returned 404, expected valid response", route)
		}
	}
}
