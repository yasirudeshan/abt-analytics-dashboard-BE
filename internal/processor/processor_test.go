package processor

import (
	"abt-analytics-dashboard/internal/models"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	processor := New()

	if processor == nil {
		t.Fatal("Expected processor to be created, got nil")
	}

	if processor.dashboardData == nil {
		t.Fatal("Expected dashboardData to be initialized, got nil")
	}

	// Check that slices are initialized
	if len(processor.dashboardData.CountryRevenues) != 0 {
		t.Errorf("Expected empty CountryRevenues slice, got %d items", len(processor.dashboardData.CountryRevenues))
	}
	if len(processor.dashboardData.TopProducts) != 0 {
		t.Errorf("Expected empty TopProducts slice, got %d items", len(processor.dashboardData.TopProducts))
	}
	if len(processor.dashboardData.MonthlySales) != 0 {
		t.Errorf("Expected empty MonthlySales slice, got %d items", len(processor.dashboardData.MonthlySales))
	}
	if len(processor.dashboardData.TopRegions) != 0 {
		t.Errorf("Expected empty TopRegions slice, got %d items", len(processor.dashboardData.TopRegions))
	}
}

func TestParseTransaction(t *testing.T) {
	processor := New()

	// Test header mapping
	headerMap := map[string]int{
		"transaction_id":   0,
		"user_id":          1,
		"country":          2,
		"region":           3,
		"product_id":       4,
		"product_name":     5,
		"category":         6,
		"price":            7,
		"quantity":         8,
		"total_price":      9,
		"stock_quantity":   10,
		"transaction_date": 11,
		"added_date":       12,
	}

	// Test record with all fields
	record := []string{
		"TXN001",        // transaction_id
		"USER001",       // user_id
		"USA",           // country
		"North America", // region
		"PROD001",       // product_id
		"Test Product",  // product_name
		"Electronics",   // category
		"99.99",         // price
		"2",             // quantity
		"199.98",        // total_price
		"100",           // stock_quantity
		"2024-01-15",    // transaction_date
		"2024-01-15",    // added_date
	}

	transaction, err := processor.parseTransaction(record, headerMap)
	if err != nil {
		t.Fatalf("Failed to parse transaction: %v", err)
	}

	// Verify parsed fields
	if transaction.TransactionID != "TXN001" {
		t.Errorf("Expected TransactionID 'TXN001', got '%s'", transaction.TransactionID)
	}
	if transaction.UserID != "USER001" {
		t.Errorf("Expected UserID 'USER001', got '%s'", transaction.UserID)
	}
	if transaction.Country != "USA" {
		t.Errorf("Expected Country 'USA', got '%s'", transaction.Country)
	}
	if transaction.Region != "North America" {
		t.Errorf("Expected Region 'North America', got '%s'", transaction.Region)
	}
	if transaction.ProductID != "PROD001" {
		t.Errorf("Expected ProductID 'PROD001', got '%s'", transaction.ProductID)
	}
	if transaction.ProductName != "Test Product" {
		t.Errorf("Expected ProductName 'Test Product', got '%s'", transaction.ProductName)
	}
	if transaction.Category != "Electronics" {
		t.Errorf("Expected Category 'Electronics', got '%s'", transaction.Category)
	}
	if transaction.Price != 99.99 {
		t.Errorf("Expected Price 99.99, got %f", transaction.Price)
	}
	if transaction.Quantity != 2 {
		t.Errorf("Expected Quantity 2, got %d", transaction.Quantity)
	}
	if transaction.TotalPrice != 199.98 {
		t.Errorf("Expected TotalPrice 199.98, got %f", transaction.TotalPrice)
	}
	if transaction.StockQuantity != 100 {
		t.Errorf("Expected StockQuantity 100, got %d", transaction.StockQuantity)
	}

	// Verify date parsing
	expectedDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	if !transaction.TransactionDate.Equal(expectedDate) {
		t.Errorf("Expected TransactionDate %v, got %v", expectedDate, transaction.TransactionDate)
	}
	if !transaction.AddedDate.Equal(expectedDate) {
		t.Errorf("Expected AddedDate %v, got %v", expectedDate, transaction.AddedDate)
	}
}

func TestParseTransactionWithMissingFields(t *testing.T) {
	processor := New()

	// Test header mapping with missing fields
	headerMap := map[string]int{
		"transaction_id": 0,
		"user_id":        1,
		"country":        2,
	}

	// Test record with missing fields
	record := []string{
		"TXN001",  // transaction_id
		"USER001", // user_id
		"USA",     // country
	}

	transaction, err := processor.parseTransaction(record, headerMap)
	if err != nil {
		t.Fatalf("Failed to parse transaction with missing fields: %v", err)
	}

	// Verify only present fields are parsed
	if transaction.TransactionID != "TXN001" {
		t.Errorf("Expected TransactionID 'TXN001', got '%s'", transaction.TransactionID)
	}
	if transaction.UserID != "USER001" {
		t.Errorf("Expected UserID 'USER001', got '%s'", transaction.UserID)
	}
	if transaction.Country != "USA" {
		t.Errorf("Expected Country 'USA', got '%s'", transaction.Country)
	}

	// Verify missing fields are zero values
	if transaction.Region != "" {
		t.Errorf("Expected Region to be empty, got '%s'", transaction.Region)
	}
	if transaction.Price != 0 {
		t.Errorf("Expected Price to be 0, got %f", transaction.Price)
	}
	if transaction.Quantity != 0 {
		t.Errorf("Expected Quantity to be 0, got %d", transaction.Quantity)
	}
}

func TestParseTransactionWithInvalidNumbers(t *testing.T) {
	processor := New()

	headerMap := map[string]int{
		"price":          0,
		"quantity":       1,
		"total_price":    2,
		"stock_quantity": 3,
	}

	// Test record with invalid numbers
	record := []string{
		"invalid_price", // price
		"invalid_qty",   // quantity
		"invalid_total", // total_price
		"invalid_stock", // stock_quantity
	}

	transaction, err := processor.parseTransaction(record, headerMap)
	if err != nil {
		t.Fatalf("Failed to parse transaction with invalid numbers: %v", err)
	}

	// Verify invalid numbers result in zero values
	if transaction.Price != 0 {
		t.Errorf("Expected Price to be 0 for invalid input, got %f", transaction.Price)
	}
	if transaction.Quantity != 0 {
		t.Errorf("Expected Quantity to be 0 for invalid input, got %d", transaction.Quantity)
	}
	if transaction.TotalPrice != 0 {
		t.Errorf("Expected TotalPrice to be 0 for invalid input, got %f", transaction.TotalPrice)
	}
	if transaction.StockQuantity != 0 {
		t.Errorf("Expected StockQuantity to be 0 for invalid input, got %d", transaction.StockQuantity)
	}
}

func TestSortCountryRevenues(t *testing.T) {
	processor := New()

	countryMap := map[string]*models.CountryRevenue{
		"key1": {Country: "USA", ProductName: "Product1", TotalRevenue: 1000.0, TransactionCount: 10},
		"key2": {Country: "UK", ProductName: "Product2", TotalRevenue: 2000.0, TransactionCount: 20},
		"key3": {Country: "Germany", ProductName: "Product3", TotalRevenue: 500.0, TransactionCount: 5},
	}

	sorted := processor.sortCountryRevenues(countryMap)

	if len(sorted) != 3 {
		t.Errorf("Expected 3 sorted items, got %d", len(sorted))
	}

	// Verify sorting by TotalRevenue (descending)
	if sorted[0].TotalRevenue != 2000.0 {
		t.Errorf("Expected first item to have highest revenue 2000.0, got %f", sorted[0].TotalRevenue)
	}
	if sorted[1].TotalRevenue != 1000.0 {
		t.Errorf("Expected second item to have revenue 1000.0, got %f", sorted[1].TotalRevenue)
	}
	if sorted[2].TotalRevenue != 500.0 {
		t.Errorf("Expected third item to have lowest revenue 500.0, got %f", sorted[2].TotalRevenue)
	}
}

func TestSortTopProducts(t *testing.T) {
	processor := New()

	productMap := map[string]*models.ProductFrequency{
		"product1": {ProductName: "Product1", PurchaseCount: 100, CurrentStock: 50},
		"product2": {ProductName: "Product2", PurchaseCount: 300, CurrentStock: 75},
		"product3": {ProductName: "Product3", PurchaseCount: 200, CurrentStock: 25},
		"product4": {ProductName: "Product4", PurchaseCount: 400, CurrentStock: 100},
	}

	sorted := processor.sortTopProducts(productMap, 3)

	if len(sorted) != 3 {
		t.Errorf("Expected 3 sorted items (limit), got %d", len(sorted))
	}

	// Verify sorting by PurchaseCount (descending)
	if sorted[0].PurchaseCount != 400 {
		t.Errorf("Expected first item to have highest purchase count 400, got %d", sorted[0].PurchaseCount)
	}
	if sorted[1].PurchaseCount != 300 {
		t.Errorf("Expected second item to have purchase count 300, got %d", sorted[1].PurchaseCount)
	}
	if sorted[2].PurchaseCount != 200 {
		t.Errorf("Expected third item to have purchase count 200, got %d", sorted[2].PurchaseCount)
	}
}

func TestSortMonthlySales(t *testing.T) {
	processor := New()

	monthMap := map[string]*models.MonthlySales{
		"2024-01": {Month: "January", Year: 2024, TotalSales: 1000.0, SalesVolume: 100},
		"2024-02": {Month: "February", Year: 2024, TotalSales: 2000.0, SalesVolume: 200},
		"2023-12": {Month: "December", Year: 2023, TotalSales: 3000.0, SalesVolume: 300},
	}

	sorted := processor.sortMonthlySales(monthMap)

	if len(sorted) != 3 {
		t.Errorf("Expected 3 sorted items, got %d", len(sorted))
	}

	// Verify sorting by Year (descending) then TotalSales (descending)
	if sorted[0].Year != 2024 {
		t.Errorf("Expected first item to be from year 2024, got %d", sorted[0].Year)
	}
	if sorted[0].TotalSales != 2000.0 {
		t.Errorf("Expected first item to have highest sales in 2024, got %f", sorted[0].TotalSales)
	}
	if sorted[1].Year != 2024 {
		t.Errorf("Expected second item to be from year 2024, got %d", sorted[1].Year)
	}
	if sorted[1].TotalSales != 1000.0 {
		t.Errorf("Expected second item to have second highest sales in 2024, got %f", sorted[1].TotalSales)
	}
	if sorted[2].Year != 2023 {
		t.Errorf("Expected third item to be from year 2023, got %d", sorted[2].Year)
	}
}

func TestSortTopRegions(t *testing.T) {
	processor := New()

	regionMap := map[string]*models.RegionRevenue{
		"region1": {Region: "North America", TotalRevenue: 10000.0, ItemsSold: 1000},
		"region2": {Region: "Europe", TotalRevenue: 15000.0, ItemsSold: 1500},
		"region3": {Region: "Asia", TotalRevenue: 5000.0, ItemsSold: 500},
	}

	sorted := processor.sortTopRegions(regionMap, 2)

	if len(sorted) != 2 {
		t.Errorf("Expected 2 sorted items (limit), got %d", len(sorted))
	}

	// Verify sorting by TotalRevenue (descending)
	if sorted[0].TotalRevenue != 15000.0 {
		t.Errorf("Expected first item to have highest revenue 15000.0, got %f", sorted[0].TotalRevenue)
	}
	if sorted[1].TotalRevenue != 10000.0 {
		t.Errorf("Expected second item to have revenue 10000.0, got %f", sorted[1].TotalRevenue)
	}
}

func TestGetDashboardData(t *testing.T) {
	processor := New()

	// Load sample data
	processor.LoadSampleData()

	data := processor.GetDashboardData()
	if data == nil {
		t.Fatal("Expected dashboard data, got nil")
	}

	// Verify data is populated
	if len(data.CountryRevenues) == 0 {
		t.Error("Expected CountryRevenues to be populated")
	}
	if len(data.TopProducts) == 0 {
		t.Error("Expected TopProducts to be populated")
	}
	if len(data.MonthlySales) == 0 {
		t.Error("Expected MonthlySales to be populated")
	}
	if len(data.TopRegions) == 0 {
		t.Error("Expected TopRegions to be populated")
	}

	// Verify metadata
	if data.LastUpdated.IsZero() {
		t.Error("Expected LastUpdated to be set")
	}
	if data.ProcessingDuration == 0 {
		t.Error("Expected ProcessingDuration to be set")
	}
	if data.RecordCount == 0 {
		t.Error("Expected RecordCount to be set")
	}
}

func TestGetCountryRevenues(t *testing.T) {
	processor := New()
	processor.LoadSampleData()

	revenues := processor.GetCountryRevenues()
	if len(revenues) == 0 {
		t.Error("Expected CountryRevenues to be populated")
	}

	// Verify structure of returned data
	for _, revenue := range revenues {
		if revenue.Country == "" {
			t.Error("Expected Country to be set")
		}
		if revenue.ProductName == "" {
			t.Error("Expected ProductName to be set")
		}
		if revenue.TotalRevenue <= 0 {
			t.Error("Expected TotalRevenue to be positive")
		}
		if revenue.TransactionCount <= 0 {
			t.Error("Expected TransactionCount to be positive")
		}
	}
}

func TestGetTopProducts(t *testing.T) {
	processor := New()
	processor.LoadSampleData()

	products := processor.GetTopProducts()
	if len(products) == 0 {
		t.Error("Expected TopProducts to be populated")
	}

	// Verify structure of returned data
	for _, product := range products {
		if product.ProductName == "" {
			t.Error("Expected ProductName to be set")
		}
		if product.PurchaseCount <= 0 {
			t.Error("Expected PurchaseCount to be positive")
		}
		if product.CurrentStock < 0 {
			t.Error("Expected CurrentStock to be non-negative")
		}
	}
}

func TestGetMonthlySales(t *testing.T) {
	processor := New()
	processor.LoadSampleData()

	sales := processor.GetMonthlySales()
	if len(sales) == 0 {
		t.Error("Expected MonthlySales to be populated")
	}

	// Verify structure of returned data
	for _, sale := range sales {
		if sale.Month == "" {
			t.Error("Expected Month to be set")
		}
		if sale.Year <= 0 {
			t.Error("Expected Year to be positive")
		}
		if sale.TotalSales <= 0 {
			t.Error("Expected TotalSales to be positive")
		}
		if sale.SalesVolume <= 0 {
			t.Error("Expected SalesVolume to be positive")
		}
	}
}

func TestGetTopRegions(t *testing.T) {
	processor := New()
	processor.LoadSampleData()

	regions := processor.GetTopRegions()
	if len(regions) == 0 {
		t.Error("Expected TopRegions to be populated")
	}

	// Verify structure of returned data
	for _, region := range regions {
		if region.Region == "" {
			t.Error("Expected Region to be set")
		}
		if region.TotalRevenue <= 0 {
			t.Error("Expected TotalRevenue to be positive")
		}
		if region.ItemsSold <= 0 {
			t.Error("Expected ItemsSold to be positive")
		}
	}
}

func TestLoadSampleData(t *testing.T) {
	processor := New()

	// Initially empty
	if len(processor.dashboardData.CountryRevenues) != 0 {
		t.Error("Expected initial CountryRevenues to be empty")
	}

	processor.LoadSampleData()

	// Should be populated after loading
	if len(processor.dashboardData.CountryRevenues) == 0 {
		t.Error("Expected CountryRevenues to be populated after loading sample data")
	}
	if len(processor.dashboardData.TopProducts) == 0 {
		t.Error("Expected TopProducts to be populated after loading sample data")
	}
	if len(processor.dashboardData.MonthlySales) == 0 {
		t.Error("Expected MonthlySales to be populated after loading sample data")
	}
	if len(processor.dashboardData.TopRegions) == 0 {
		t.Error("Expected TopRegions to be populated after loading sample data")
	}

	// Verify metadata is set
	if processor.dashboardData.LastUpdated.IsZero() {
		t.Error("Expected LastUpdated to be set after loading sample data")
	}
	// Note: ProcessingDuration might be very small for sample data
	if processor.dashboardData.RecordCount == 0 {
		t.Error("Expected RecordCount to be set after loading sample data")
	}
}
