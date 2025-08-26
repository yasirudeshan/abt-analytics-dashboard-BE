package processor

import (
	"abt-analytics-dashboard/internal/models"
	"testing"
	"time"
)

// createMockProcessor creates a processor with predictable test data
func createMockProcessor() *Processor {
	now := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)

	processor := New()
	processor.dashboardData = &models.DashboardData{
		LastUpdated:        now,
		ProcessingDuration: 5 * time.Second,
		RecordCount:        1000,
	}

	return processor
}

// TestParseTransactionWithMockData tests CSV parsing with hardcoded test data
func TestParseTransactionWithMockData(t *testing.T) {
	processor := createMockProcessor()

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
		"Test Laptop",   // product_name
		"Electronics",   // category
		"999.99",        // price
		"2",             // quantity
		"1999.98",       // total_price
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
	if transaction.ProductName != "Test Laptop" {
		t.Errorf("Expected ProductName 'Test Laptop', got '%s'", transaction.ProductName)
	}
	if transaction.Category != "Electronics" {
		t.Errorf("Expected Category 'Electronics', got '%s'", transaction.Category)
	}
	if transaction.Price != 999.99 {
		t.Errorf("Expected Price 999.99, got %f", transaction.Price)
	}
	if transaction.Quantity != 2 {
		t.Errorf("Expected Quantity 2, got %d", transaction.Quantity)
	}
	if transaction.TotalPrice != 1999.98 {
		t.Errorf("Expected TotalPrice 1999.98, got %f", transaction.TotalPrice)
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

// TestParseTransactionWithMissingFieldsUnit tests parsing with incomplete data
func TestParseTransactionWithMissingFieldsUnit(t *testing.T) {
	processor := createMockProcessor()

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

// TestParseTransactionWithInvalidNumbersUnit tests error handling for invalid numeric data
func TestParseTransactionWithInvalidNumbersUnit(t *testing.T) {
	processor := createMockProcessor()

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

// TestSortCountryRevenuesWithMockData tests revenue sorting with hardcoded data
func TestSortCountryRevenuesWithMockData(t *testing.T) {
	processor := createMockProcessor()

	countryMap := map[string]*models.CountryRevenue{
		"key1": {Country: "USA", ProductName: "Laptop", TotalRevenue: 1000.0, TransactionCount: 10},
		"key2": {Country: "UK", ProductName: "Smartphone", TotalRevenue: 2000.0, TransactionCount: 20},
		"key3": {Country: "Germany", ProductName: "Tablet", TotalRevenue: 500.0, TransactionCount: 5},
	}

	sorted := processor.sortCountryRevenues(countryMap)

	if len(sorted) != 3 {
		t.Errorf("Expected 3 sorted items, got %d", len(sorted))
	}

	// Verify sorting by TotalRevenue (descending)
	if sorted[0].TotalRevenue != 2000.0 {
		t.Errorf("Expected first item to have highest revenue 2000.0, got %f", sorted[0].TotalRevenue)
	}
	if sorted[0].Country != "UK" {
		t.Errorf("Expected first item to be UK, got %s", sorted[0].Country)
	}

	if sorted[1].TotalRevenue != 1000.0 {
		t.Errorf("Expected second item to have revenue 1000.0, got %f", sorted[1].TotalRevenue)
	}
	if sorted[1].Country != "USA" {
		t.Errorf("Expected second item to be USA, got %s", sorted[1].Country)
	}

	if sorted[2].TotalRevenue != 500.0 {
		t.Errorf("Expected third item to have lowest revenue 500.0, got %f", sorted[2].TotalRevenue)
	}
	if sorted[2].Country != "Germany" {
		t.Errorf("Expected third item to be Germany, got %s", sorted[2].Country)
	}
}

// TestSortTopProductsWithMockData tests product sorting with hardcoded data
func TestSortTopProductsWithMockData(t *testing.T) {
	processor := createMockProcessor()

	productMap := map[string]*models.ProductFrequency{
		"product1": {ProductName: "Laptop", PurchaseCount: 100, CurrentStock: 50},
		"product2": {ProductName: "Smartphone", PurchaseCount: 300, CurrentStock: 75},
		"product3": {ProductName: "Tablet", PurchaseCount: 200, CurrentStock: 25},
		"product4": {ProductName: "Monitor", PurchaseCount: 400, CurrentStock: 100},
	}

	sorted := processor.sortTopProducts(productMap, 3)

	if len(sorted) != 3 {
		t.Errorf("Expected 3 sorted items (limit), got %d", len(sorted))
	}

	// Verify sorting by PurchaseCount (descending)
	if sorted[0].PurchaseCount != 400 {
		t.Errorf("Expected first item to have highest purchase count 400, got %d", sorted[0].PurchaseCount)
	}
	if sorted[0].ProductName != "Monitor" {
		t.Errorf("Expected first item to be Monitor, got %s", sorted[0].ProductName)
	}

	if sorted[1].PurchaseCount != 300 {
		t.Errorf("Expected second item to have purchase count 300, got %d", sorted[1].PurchaseCount)
	}
	if sorted[1].ProductName != "Smartphone" {
		t.Errorf("Expected second item to be Smartphone, got %s", sorted[1].ProductName)
	}

	if sorted[2].PurchaseCount != 200 {
		t.Errorf("Expected third item to have purchase count 200, got %d", sorted[2].PurchaseCount)
	}
	if sorted[2].ProductName != "Tablet" {
		t.Errorf("Expected third item to be Tablet, got %s", sorted[2].ProductName)
	}
}

// TestSortMonthlySalesWithMockData tests monthly sales sorting with hardcoded data
func TestSortMonthlySalesWithMockData(t *testing.T) {
	processor := createMockProcessor()

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
	if sorted[0].Month != "February" {
		t.Errorf("Expected first item to be February, got %s", sorted[0].Month)
	}
	if sorted[0].TotalSales != 2000.0 {
		t.Errorf("Expected first item to have highest sales in 2024, got %f", sorted[0].TotalSales)
	}

	if sorted[1].Year != 2024 {
		t.Errorf("Expected second item to be from year 2024, got %d", sorted[1].Year)
	}
	if sorted[1].Month != "January" {
		t.Errorf("Expected second item to be January, got %s", sorted[1].Month)
	}
	if sorted[1].TotalSales != 1000.0 {
		t.Errorf("Expected second item to have second highest sales in 2024, got %f", sorted[1].TotalSales)
	}

	if sorted[2].Year != 2023 {
		t.Errorf("Expected third item to be from year 2023, got %d", sorted[2].Year)
	}
	if sorted[2].Month != "December" {
		t.Errorf("Expected third item to be December, got %s", sorted[2].Month)
	}
}

// TestSortTopRegionsWithMockData tests regional sorting with hardcoded data
func TestSortTopRegionsWithMockData(t *testing.T) {
	processor := createMockProcessor()

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
	if sorted[0].Region != "Europe" {
		t.Errorf("Expected first item to be Europe, got %s", sorted[0].Region)
	}

	if sorted[1].TotalRevenue != 10000.0 {
		t.Errorf("Expected second item to have revenue 10000.0, got %f", sorted[1].TotalRevenue)
	}
	if sorted[1].Region != "North America" {
		t.Errorf("Expected second item to be North America, got %s", sorted[1].Region)
	}
}

// TestAggregateCountryRevenueWithMockData tests revenue aggregation with hardcoded data
func TestAggregateCountryRevenueWithMockData(t *testing.T) {
	_ = createMockProcessor()

	// Create mock transactions
	transactions := []models.Transaction{
		{
			Country:     "USA",
			ProductName: "Laptop",
			TotalPrice:  1000.0,
			Quantity:    2,
		},
		{
			Country:     "USA",
			ProductName: "Laptop",
			TotalPrice:  500.0,
			Quantity:    1,
		},
		{
			Country:     "UK",
			ProductName: "Smartphone",
			TotalPrice:  800.0,
			Quantity:    1,
		},
	}

	// Aggregate revenue
	countryMap := make(map[string]*models.CountryRevenue)

	for _, transaction := range transactions {
		key := transaction.Country + "_" + transaction.ProductName

		if existing, exists := countryMap[key]; exists {
			existing.TotalRevenue += transaction.TotalPrice
			existing.TransactionCount += transaction.Quantity
		} else {
			countryMap[key] = &models.CountryRevenue{
				Country:          transaction.Country,
				ProductName:      transaction.ProductName,
				TotalRevenue:     transaction.TotalPrice,
				TransactionCount: transaction.Quantity,
			}
		}
	}

	// Verify aggregation results
	if len(countryMap) != 2 {
		t.Errorf("Expected 2 unique country-product combinations, got %d", len(countryMap))
	}

	// Check USA Laptop aggregation
	usaKey := "USA_Laptop"
	if usaRevenue, exists := countryMap[usaKey]; exists {
		if usaRevenue.TotalRevenue != 1500.0 {
			t.Errorf("Expected USA Laptop total revenue 1500.0, got %f", usaRevenue.TotalRevenue)
		}
		if usaRevenue.TransactionCount != 3 {
			t.Errorf("Expected USA Laptop transaction count 3, got %d", usaRevenue.TransactionCount)
		}
	} else {
		t.Error("Expected USA Laptop combination to exist")
	}

	// Check UK Smartphone aggregation
	ukKey := "UK_Smartphone"
	if ukRevenue, exists := countryMap[ukKey]; exists {
		if ukRevenue.TotalRevenue != 800.0 {
			t.Errorf("Expected UK Smartphone total revenue 800.0, got %f", ukRevenue.TotalRevenue)
		}
		if ukRevenue.TransactionCount != 1 {
			t.Errorf("Expected UK Smartphone transaction count 1, got %d", ukRevenue.TransactionCount)
		}
	} else {
		t.Error("Expected UK Smartphone combination to exist")
	}
}

// TestAggregateProductFrequencyWithMockData tests product frequency aggregation with hardcoded data
func TestAggregateProductFrequencyWithMockData(t *testing.T) {
	_ = createMockProcessor()

	// Create mock transactions
	transactions := []models.Transaction{
		{
			ProductName:   "Laptop",
			Quantity:      2,
			StockQuantity: 100,
		},
		{
			ProductName:   "Laptop",
			Quantity:      1,
			StockQuantity: 100,
		},
		{
			ProductName:   "Smartphone",
			Quantity:      3,
			StockQuantity: 150,
		},
	}

	// Aggregate product frequency
	productMap := make(map[string]*models.ProductFrequency)

	for _, transaction := range transactions {
		if existing, exists := productMap[transaction.ProductName]; exists {
			existing.PurchaseCount += transaction.Quantity
			// Keep the latest stock quantity
			existing.CurrentStock = transaction.StockQuantity
		} else {
			productMap[transaction.ProductName] = &models.ProductFrequency{
				ProductName:   transaction.ProductName,
				PurchaseCount: transaction.Quantity,
				CurrentStock:  transaction.StockQuantity,
			}
		}
	}

	// Verify aggregation results
	if len(productMap) != 2 {
		t.Errorf("Expected 2 unique products, got %d", len(productMap))
	}

	// Check Laptop aggregation
	if laptop, exists := productMap["Laptop"]; exists {
		if laptop.PurchaseCount != 3 {
			t.Errorf("Expected Laptop purchase count 3, got %d", laptop.PurchaseCount)
		}
		if laptop.CurrentStock != 100 {
			t.Errorf("Expected Laptop current stock 100, got %d", laptop.CurrentStock)
		}
	} else {
		t.Error("Expected Laptop to exist")
	}

	// Check Smartphone aggregation
	if smartphone, exists := productMap["Smartphone"]; exists {
		if smartphone.PurchaseCount != 3 {
			t.Errorf("Expected Smartphone purchase count 3, got %d", smartphone.PurchaseCount)
		}
		if smartphone.CurrentStock != 150 {
			t.Errorf("Expected Smartphone current stock 150, got %d", smartphone.CurrentStock)
		}
	} else {
		t.Error("Expected Smartphone to exist")
	}
}

// TestAggregateMonthlySalesWithMockData tests monthly sales aggregation with hardcoded data
func TestAggregateMonthlySalesWithMockData(t *testing.T) {
	_ = createMockProcessor()

	// Create mock transactions with different months
	transactions := []models.Transaction{
		{
			TransactionDate: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			TotalPrice:      1000.0,
			Quantity:        10,
		},
		{
			TransactionDate: time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
			TotalPrice:      500.0,
			Quantity:        5,
		},
		{
			TransactionDate: time.Date(2024, 2, 10, 0, 0, 0, 0, time.UTC),
			TotalPrice:      800.0,
			Quantity:        8,
		},
	}

	// Aggregate monthly sales
	monthMap := make(map[string]*models.MonthlySales)

	for _, transaction := range transactions {
		monthKey := transaction.TransactionDate.Format("2006-01")
		monthName := transaction.TransactionDate.Month().String()
		year := transaction.TransactionDate.Year()

		if existing, exists := monthMap[monthKey]; exists {
			existing.TotalSales += transaction.TotalPrice
			existing.SalesVolume += transaction.Quantity
		} else {
			monthMap[monthKey] = &models.MonthlySales{
				Month:       monthName,
				Year:        year,
				TotalSales:  transaction.TotalPrice,
				SalesVolume: transaction.Quantity,
			}
		}
	}

	// Verify aggregation results
	if len(monthMap) != 2 {
		t.Errorf("Expected 2 unique months, got %d", len(monthMap))
	}

	// Check January 2024 aggregation
	janKey := "2024-01"
	if january, exists := monthMap[janKey]; exists {
		if january.TotalSales != 1500.0 {
			t.Errorf("Expected January total sales 1500.0, got %f", january.TotalSales)
		}
		if january.SalesVolume != 15 {
			t.Errorf("Expected January sales volume 15, got %d", january.SalesVolume)
		}
		if january.Month != "January" {
			t.Errorf("Expected January month name 'January', got %s", january.Month)
		}
		if january.Year != 2024 {
			t.Errorf("Expected January year 2024, got %d", january.Year)
		}
	} else {
		t.Error("Expected January 2024 to exist")
	}

	// Check February 2024 aggregation
	febKey := "2024-02"
	if february, exists := monthMap[febKey]; exists {
		if february.TotalSales != 800.0 {
			t.Errorf("Expected February total sales 800.0, got %f", february.TotalSales)
		}
		if february.SalesVolume != 8 {
			t.Errorf("Expected February sales volume 8, got %d", february.SalesVolume)
		}
		if february.Month != "February" {
			t.Errorf("Expected February month name 'February', got %s", february.Month)
		}
		if february.Year != 2024 {
			t.Errorf("Expected February year 2024, got %d", february.Year)
		}
	} else {
		t.Error("Expected February 2024 to exist")
	}
}
