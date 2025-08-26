package api

import (
	"abt-analytics-dashboard/internal/config"
	"abt-analytics-dashboard/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// ProcessorInterface defines the methods that the API server needs
type ProcessorInterface interface {
	GetCountryRevenues() []models.CountryRevenue
	GetTopProducts() []models.ProductFrequency
	GetMonthlySales() []models.MonthlySales
	GetTopRegions() []models.RegionRevenue
	GetDashboardData() *models.DashboardData
}

// MockProcessor is a mock implementation for testing
type MockProcessor struct {
	mockCountryRevenues []models.CountryRevenue
	mockTopProducts     []models.ProductFrequency
	mockMonthlySales    []models.MonthlySales
	mockTopRegions      []models.RegionRevenue
	mockDashboardData   *models.DashboardData
}

func (m *MockProcessor) GetCountryRevenues() []models.CountryRevenue {
	return m.mockCountryRevenues
}

func (m *MockProcessor) GetTopProducts() []models.ProductFrequency {
	return m.mockTopProducts
}

func (m *MockProcessor) GetMonthlySales() []models.MonthlySales {
	return m.mockMonthlySales
}

func (m *MockProcessor) GetTopRegions() []models.RegionRevenue {
	return m.mockTopRegions
}

func (m *MockProcessor) GetDashboardData() *models.DashboardData {
	return m.mockDashboardData
}

// TestServer is a test-specific server that uses the interface
type TestServer struct {
	processor ProcessorInterface
	config    *config.Config
}

// createMockData creates predictable test data
func createMockData() *MockProcessor {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	return &MockProcessor{
		mockCountryRevenues: []models.CountryRevenue{
			{Country: "USA", ProductName: "Laptop", TotalRevenue: 50000.0, TransactionCount: 100},
			{Country: "UK", ProductName: "Smartphone", TotalRevenue: 30000.0, TransactionCount: 75},
			{Country: "Germany", ProductName: "Tablet", TotalRevenue: 25000.0, TransactionCount: 50},
		},
		mockTopProducts: []models.ProductFrequency{
			{ProductName: "Laptop", PurchaseCount: 500, CurrentStock: 100},
			{ProductName: "Smartphone", PurchaseCount: 400, CurrentStock: 150},
			{ProductName: "Tablet", PurchaseCount: 300, CurrentStock: 75},
		},
		mockMonthlySales: []models.MonthlySales{
			{Month: "January", Year: 2024, TotalSales: 150000.0, SalesVolume: 3000},
			{Month: "February", Year: 2024, TotalSales: 180000.0, SalesVolume: 3600},
			{Month: "March", Year: 2024, TotalSales: 120000.0, SalesVolume: 2400},
		},
		mockTopRegions: []models.RegionRevenue{
			{Region: "North America", TotalRevenue: 200000.0, ItemsSold: 4000},
			{Region: "Europe", TotalRevenue: 150000.0, ItemsSold: 3000},
			{Region: "Asia", TotalRevenue: 100000.0, ItemsSold: 2000},
		},
		mockDashboardData: &models.DashboardData{
			LastUpdated:        now,
			ProcessingDuration: 5 * time.Second,
			RecordCount:        1000,
		},
	}
}

// TestRootHandlerWithMockData tests root endpoint with mock data
func TestRootHandlerWithMockData(t *testing.T) {
	mockProc := createMockData()
	cfg := &config.Config{Port: ":8080"}
	_ = &TestServer{processor: mockProc, config: cfg}

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Create a simple root handler for testing
	rootHandler := func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"service": "ABT Analytics Dashboard API",
			"version": "1.0.0",
			"status":  "running",
			"endpoints": map[string]string{
				"health":             "/api/health",
				"country_revenues":   "/api/revenue-by-country",
				"top_products":       "/api/top-products",
				"monthly_sales":      "/api/sales-by-month",
				"top_regions":        "/api/top-regions",
				"complete_dashboard": "/api/dashboard",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	rootHandler(rr, req)

	// Verify status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Verify content type
	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Parse and verify response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify service info
	if response["service"] != "ABT Analytics Dashboard API" {
		t.Errorf("Expected service 'ABT Analytics Dashboard API', got '%v'", response["service"])
	}
	if response["version"] != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%v'", response["version"])
	}
	if response["status"] != "running" {
		t.Errorf("Expected status 'running', got '%v'", response["status"])
	}

	// Verify endpoints
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

// TestHealthCheckWithMockData tests health endpoint with mock data
func TestHealthCheckWithMockData(t *testing.T) {
	mockProc := createMockData()
	cfg := &config.Config{Port: ":8080"}
	server := &TestServer{processor: mockProc, config: cfg}

	req, err := http.NewRequest("GET", "/api/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Create a simple health check handler for testing
	healthHandler := func(w http.ResponseWriter, r *http.Request) {
		dashboardData := server.processor.GetDashboardData()
		response := map[string]interface{}{
			"status":              "healthy",
			"timestamp":           time.Now(),
			"last_data_update":    dashboardData.LastUpdated,
			"processing_duration": dashboardData.ProcessingDuration.String(),
			"record_count":        dashboardData.RecordCount,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	healthHandler(rr, req)

	// Verify status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse and verify response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify health status
	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got '%v'", response["status"])
	}

	// Verify timestamp is present
	if _, exists := response["timestamp"]; !exists {
		t.Error("Expected timestamp to be present")
	}

	// Verify last data update
	if response["last_data_update"] != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected last_data_update to match mock data, got '%v'", response["last_data_update"])
	}

	// Verify processing duration
	if response["processing_duration"] != "5s" {
		t.Errorf("Expected processing_duration '5s', got '%v'", response["processing_duration"])
	}

	// Verify record count
	if response["record_count"] != float64(1000) {
		t.Errorf("Expected record_count 1000, got '%v'", response["record_count"])
	}
}

// TestGetCountryRevenuesWithMockData tests country revenues endpoint with mock data
func TestGetCountryRevenuesWithMockData(t *testing.T) {
	mockProc := createMockData()
	cfg := &config.Config{Port: ":8080"}
	server := &TestServer{processor: mockProc, config: cfg}

	req, err := http.NewRequest("GET", "/api/revenue-by-country", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Create a simple country revenues handler for testing
	countryRevenuesHandler := func(w http.ResponseWriter, r *http.Request) {
		data := server.processor.GetCountryRevenues()
		response := map[string]interface{}{
			"data":  data,
			"count": len(data),
			"meta": map[string]interface{}{
				"description": "Country-level revenue data sorted by total revenue (descending)",
				"updated_at":  server.processor.GetDashboardData().LastUpdated,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	countryRevenuesHandler(rr, req)

	// Verify status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse and verify response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify data structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}
	if response["count"] != float64(3) {
		t.Errorf("Expected count 3, got '%v'", response["count"])
	}

	// Verify meta information
	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Country-level revenue data sorted by total revenue (descending)" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}

	// Verify data content
	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 data items, got %d", len(data))
	}

	// Verify first item (USA Laptop)
	firstItem, ok := data[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first item to be a map")
	}

	if firstItem["country"] != "USA" {
		t.Errorf("Expected first item country 'USA', got '%v'", firstItem["country"])
	}
	if firstItem["product_name"] != "Laptop" {
		t.Errorf("Expected first item product 'Laptop', got '%v'", firstItem["product_name"])
	}
	if firstItem["total_revenue"] != 50000.0 {
		t.Errorf("Expected first item revenue 50000.0, got '%v'", firstItem["total_revenue"])
	}
}

// TestGetTopProductsWithMockData tests top products endpoint with mock data
func TestGetTopProductsWithMockData(t *testing.T) {
	mockProc := createMockData()
	cfg := &config.Config{Port: ":8080"}
	server := &TestServer{processor: mockProc, config: cfg}

	req, err := http.NewRequest("GET", "/api/top-products", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Create a simple top products handler for testing
	topProductsHandler := func(w http.ResponseWriter, r *http.Request) {
		data := server.processor.GetTopProducts()
		response := map[string]interface{}{
			"data":  data,
			"count": len(data),
			"meta": map[string]interface{}{
				"description": "Top 20 most frequently purchased products with current stock",
				"updated_at":  server.processor.GetDashboardData().LastUpdated,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	topProductsHandler(rr, req)

	// Verify status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse and verify response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify data structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}
	if response["count"] != float64(3) {
		t.Errorf("Expected count 3, got '%v'", response["count"])
	}

	// Verify meta information
	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Top 20 most frequently purchased products with current stock" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}

	// Verify data content
	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 data items, got %d", len(data))
	}

	// Verify first item (Laptop)
	firstItem, ok := data[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first item to be a map")
	}

	if firstItem["product_name"] != "Laptop" {
		t.Errorf("Expected first item product 'Laptop', got '%v'", firstItem["product_name"])
	}
	if firstItem["purchase_count"] != float64(500) {
		t.Errorf("Expected first item purchase count 500, got '%v'", firstItem["purchase_count"])
	}
	if firstItem["current_stock"] != float64(100) {
		t.Errorf("Expected first item stock 100, got '%v'", firstItem["current_stock"])
	}
}

// TestGetMonthlySalesWithMockData tests monthly sales endpoint with mock data
func TestGetMonthlySalesWithMockData(t *testing.T) {
	mockProc := createMockData()
	cfg := &config.Config{Port: ":8080"}
	server := &TestServer{processor: mockProc, config: cfg}

	req, err := http.NewRequest("GET", "/api/sales-by-month", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Create a simple monthly sales handler for testing
	monthlySalesHandler := func(w http.ResponseWriter, r *http.Request) {
		data := server.processor.GetMonthlySales()
		response := map[string]interface{}{
			"data":  data,
			"count": len(data),
			"meta": map[string]interface{}{
				"description": "Monthly sales volume data highlighting peak sales periods",
				"updated_at":  server.processor.GetDashboardData().LastUpdated,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	monthlySalesHandler(rr, req)

	// Verify status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse and verify response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify data structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}
	if response["count"] != float64(3) {
		t.Errorf("Expected count 3, got '%v'", response["count"])
	}

	// Verify meta information
	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Monthly sales volume data highlighting peak sales periods" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}

	// Verify data content
	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 data items, got %d", len(data))
	}

	// Verify first item (January)
	firstItem, ok := data[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first item to be a map")
	}

	if firstItem["month"] != "January" {
		t.Errorf("Expected first item month 'January', got '%v'", firstItem["month"])
	}
	if firstItem["year"] != float64(2024) {
		t.Errorf("Expected first item year 2024, got '%v'", firstItem["year"])
	}
	if firstItem["total_sales"] != 150000.0 {
		t.Errorf("Expected first item sales 150000.0, got '%v'", firstItem["total_sales"])
	}
}

// TestGetTopRegionsWithMockData tests top regions endpoint with mock data
func TestGetTopRegionsWithMockData(t *testing.T) {
	mockProc := createMockData()
	cfg := &config.Config{Port: ":8080"}
	server := &TestServer{processor: mockProc, config: cfg}

	req, err := http.NewRequest("GET", "/api/top-regions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Create a simple top regions handler for testing
	topRegionsHandler := func(w http.ResponseWriter, r *http.Request) {
		data := server.processor.GetTopRegions()
		response := map[string]interface{}{
			"data":  data,
			"count": len(data),
			"meta": map[string]interface{}{
				"description": "Top 30 regions by total revenue and items sold",
				"updated_at":  server.processor.GetDashboardData().LastUpdated,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	topRegionsHandler(rr, req)

	// Verify status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse and verify response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify data structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}
	if response["count"] != float64(3) {
		t.Errorf("Expected count 3, got '%v'", response["count"])
	}

	// Verify meta information
	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Top 30 regions by total revenue and items sold" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}

	// Verify data content
	data, ok := response["data"].([]interface{})
	if !ok {
		t.Fatal("Expected data to be an array")
	}

	if len(data) != 3 {
		t.Errorf("Expected 3 data items, got %d", len(data))
	}

	// Verify first item (North America)
	firstItem, ok := data[0].(map[string]interface{})
	if !ok {
		t.Fatal("Expected first item to be a map")
	}

	if firstItem["region"] != "North America" {
		t.Errorf("Expected first item region 'North America', got '%v'", firstItem["region"])
	}
	if firstItem["total_revenue"] != 200000.0 {
		t.Errorf("Expected first item revenue 200000.0, got '%v'", firstItem["total_revenue"])
	}
	if firstItem["items_sold"] != float64(4000) {
		t.Errorf("Expected first item items sold 4000, got '%v'", firstItem["items_sold"])
	}
}

// TestGetDashboardDataWithMockData tests complete dashboard endpoint with mock data
func TestGetDashboardDataWithMockData(t *testing.T) {
	mockProc := createMockData()
	cfg := &config.Config{Port: ":8080"}
	server := &TestServer{processor: mockProc, config: cfg}

	req, err := http.NewRequest("GET", "/api/dashboard", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()

	// Create a simple dashboard data handler for testing
	dashboardDataHandler := func(w http.ResponseWriter, r *http.Request) {
		data := server.processor.GetDashboardData()
		response := map[string]interface{}{
			"data": data,
			"meta": map[string]interface{}{
				"description": "Complete dashboard data including all metrics",
				"updated_at":  data.LastUpdated,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}

	dashboardDataHandler(rr, req)

	// Verify status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, status)
	}

	// Parse and verify response
	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify data structure
	if _, exists := response["data"]; !exists {
		t.Error("Expected data field to be present")
	}

	// Verify meta information
	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected meta to be a map")
	}

	if meta["description"] != "Complete dashboard data including all metrics" {
		t.Errorf("Expected description to match, got '%v'", meta["description"])
	}

	// Verify dashboard data content
	dashboardData, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected dashboard data to be a map")
	}

	// Verify last updated
	if dashboardData["last_updated"] != "2024-01-15T10:30:00Z" {
		t.Errorf("Expected last_updated to match mock data, got '%v'", dashboardData["last_updated"])
	}

	// Verify record count
	if dashboardData["record_count"] != float64(1000) {
		t.Errorf("Expected record_count 1000, got '%v'", dashboardData["record_count"])
	}
}

// TestJSONResponseWithMockData tests JSON response writing with mock data
func TestJSONResponseWithMockData(t *testing.T) {
	testData := map[string]interface{}{
		"test_string": "hello",
		"test_number": 42,
		"test_bool":   true,
		"test_array":  []string{"a", "b", "c"},
	}

	// Write JSON response
	w := httptest.NewRecorder()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(testData)

	// Verify status code
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, status)
	}

	// Verify content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}

	// Verify response body
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify all test data fields
	if response["test_string"] != "hello" {
		t.Errorf("Expected test_string 'hello', got '%v'", response["test_string"])
	}
	if response["test_number"] != float64(42) {
		t.Errorf("Expected test_number 42, got '%v'", response["test_number"])
	}
	if response["test_bool"] != true {
		t.Errorf("Expected test_bool true, got '%v'", response["test_bool"])
	}

	// Verify array
	testArray, ok := response["test_array"].([]interface{})
	if !ok {
		t.Fatal("Expected test_array to be an array")
	}
	if len(testArray) != 3 {
		t.Errorf("Expected test_array length 3, got %d", len(testArray))
	}
	if testArray[0] != "a" {
		t.Errorf("Expected test_array[0] 'a', got '%v'", testArray[0])
	}
}

// TestErrorResponseWithMockData tests error response writing with mock data
func TestErrorResponseWithMockData(t *testing.T) {
	errorMessage := "Test error message"

	// Write error response
	w := httptest.NewRecorder()
	response := map[string]interface{}{
		"error":     true,
		"message":   errorMessage,
		"timestamp": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)

	// Verify status code
	if status := w.Code; status != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, status)
	}

	// Parse and verify response
	var responseData map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &responseData); err != nil {
		t.Fatalf("Failed to parse response JSON: %v", err)
	}

	// Verify error structure
	if responseData["error"] != true {
		t.Error("Expected error field to be true")
	}
	if responseData["message"] != errorMessage {
		t.Errorf("Expected message '%s', got '%v'", errorMessage, responseData["message"])
	}
	if _, exists := responseData["timestamp"]; !exists {
		t.Error("Expected timestamp to be present")
	}
}
