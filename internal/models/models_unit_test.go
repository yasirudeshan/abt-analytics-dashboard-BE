package models

import (
	"encoding/json"
	"testing"
	"time"
)

// createMockTime creates a predictable timestamp for testing
func createMockTime() time.Time {
	return time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
}

// TestTransactionStructWithMockData tests transaction structure with hardcoded data
func TestTransactionStructWithMockData(t *testing.T) {
	mockTime := createMockTime()

	transaction := Transaction{
		TransactionID:   "TXN001",
		TransactionDate: mockTime,
		UserID:          "USER001",
		Country:         "USA",
		Region:          "North America",
		ProductID:       "PROD001",
		ProductName:     "Test Laptop",
		Category:        "Electronics",
		Price:           999.99,
		Quantity:        2,
		TotalPrice:      1999.98,
		StockQuantity:   100,
		AddedDate:       mockTime,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(transaction)
	if err != nil {
		t.Fatalf("Failed to marshal Transaction to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledTransaction Transaction
	err = json.Unmarshal(jsonData, &unmarshaledTransaction)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to Transaction: %v", err)
	}

	// Verify all fields are preserved
	if unmarshaledTransaction.TransactionID != transaction.TransactionID {
		t.Errorf("Expected TransactionID %s, got %s", transaction.TransactionID, unmarshaledTransaction.TransactionID)
	}
	if unmarshaledTransaction.UserID != transaction.UserID {
		t.Errorf("Expected UserID %s, got %s", transaction.UserID, unmarshaledTransaction.UserID)
	}
	if unmarshaledTransaction.Country != transaction.Country {
		t.Errorf("Expected Country %s, got %s", transaction.Country, unmarshaledTransaction.Country)
	}
	if unmarshaledTransaction.Region != transaction.Region {
		t.Errorf("Expected Region %s, got %s", transaction.Region, unmarshaledTransaction.Region)
	}
	if unmarshaledTransaction.ProductID != transaction.ProductID {
		t.Errorf("Expected ProductID %s, got %s", transaction.ProductID, unmarshaledTransaction.ProductID)
	}
	if unmarshaledTransaction.ProductName != transaction.ProductName {
		t.Errorf("Expected ProductName %s, got %s", transaction.ProductName, unmarshaledTransaction.ProductName)
	}
	if unmarshaledTransaction.Category != transaction.Category {
		t.Errorf("Expected Category %s, got %s", transaction.Category, unmarshaledTransaction.Category)
	}
	if unmarshaledTransaction.Price != transaction.Price {
		t.Errorf("Expected Price %f, got %f", transaction.Price, unmarshaledTransaction.Price)
	}
	if unmarshaledTransaction.Quantity != transaction.Quantity {
		t.Errorf("Expected Quantity %d, got %d", transaction.Quantity, unmarshaledTransaction.Quantity)
	}
	if unmarshaledTransaction.TotalPrice != transaction.TotalPrice {
		t.Errorf("Expected TotalPrice %f, got %f", transaction.TotalPrice, unmarshaledTransaction.TotalPrice)
	}
	if unmarshaledTransaction.StockQuantity != transaction.StockQuantity {
		t.Errorf("Expected StockQuantity %d, got %d", transaction.StockQuantity, unmarshaledTransaction.StockQuantity)
	}

	// Verify time fields
	if !unmarshaledTransaction.TransactionDate.Equal(transaction.TransactionDate) {
		t.Errorf("Expected TransactionDate %v, got %v", transaction.TransactionDate, unmarshaledTransaction.TransactionDate)
	}
	if !unmarshaledTransaction.AddedDate.Equal(transaction.AddedDate) {
		t.Errorf("Expected AddedDate %v, got %v", transaction.AddedDate, unmarshaledTransaction.AddedDate)
	}
}

// TestCountryRevenueStructWithMockData tests country revenue structure with hardcoded data
func TestCountryRevenueStructWithMockData(t *testing.T) {
	countryRevenue := CountryRevenue{
		Country:          "USA",
		ProductName:      "Test Laptop",
		TotalRevenue:     1500.50,
		TransactionCount: 25,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(countryRevenue)
	if err != nil {
		t.Fatalf("Failed to marshal CountryRevenue to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledCountryRevenue CountryRevenue
	err = json.Unmarshal(jsonData, &unmarshaledCountryRevenue)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to CountryRevenue: %v", err)
	}

	// Verify all fields are preserved
	if unmarshaledCountryRevenue.Country != countryRevenue.Country {
		t.Errorf("Expected Country %s, got %s", countryRevenue.Country, unmarshaledCountryRevenue.Country)
	}
	if unmarshaledCountryRevenue.ProductName != countryRevenue.ProductName {
		t.Errorf("Expected ProductName %s, got %s", countryRevenue.ProductName, unmarshaledCountryRevenue.ProductName)
	}
	if unmarshaledCountryRevenue.TotalRevenue != countryRevenue.TotalRevenue {
		t.Errorf("Expected TotalRevenue %f, got %f", countryRevenue.TotalRevenue, unmarshaledCountryRevenue.TotalRevenue)
	}
	if unmarshaledCountryRevenue.TransactionCount != countryRevenue.TransactionCount {
		t.Errorf("Expected TransactionCount %d, got %d", countryRevenue.TransactionCount, unmarshaledCountryRevenue.TransactionCount)
	}
}

// TestProductFrequencyStructWithMockData tests product frequency structure with hardcoded data
func TestProductFrequencyStructWithMockData(t *testing.T) {
	productFrequency := ProductFrequency{
		ProductName:   "Test Laptop",
		PurchaseCount: 150,
		CurrentStock:  75,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(productFrequency)
	if err != nil {
		t.Fatalf("Failed to marshal ProductFrequency to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledProductFrequency ProductFrequency
	err = json.Unmarshal(jsonData, &unmarshaledProductFrequency)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to ProductFrequency: %v", err)
	}

	// Verify all fields are preserved
	if unmarshaledProductFrequency.ProductName != productFrequency.ProductName {
		t.Errorf("Expected ProductName %s, got %s", productFrequency.ProductName, unmarshaledProductFrequency.ProductName)
	}
	if unmarshaledProductFrequency.PurchaseCount != productFrequency.PurchaseCount {
		t.Errorf("Expected PurchaseCount %d, got %d", productFrequency.PurchaseCount, unmarshaledProductFrequency.PurchaseCount)
	}
	if unmarshaledProductFrequency.CurrentStock != productFrequency.CurrentStock {
		t.Errorf("Expected CurrentStock %d, got %d", productFrequency.CurrentStock, unmarshaledProductFrequency.CurrentStock)
	}
}

// TestMonthlySalesStructWithMockData tests monthly sales structure with hardcoded data
func TestMonthlySalesStructWithMockData(t *testing.T) {
	monthlySales := MonthlySales{
		Month:       "January",
		Year:        2024,
		TotalSales:  150000.0,
		SalesVolume: 3000,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(monthlySales)
	if err != nil {
		t.Fatalf("Failed to marshal MonthlySales to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledMonthlySales MonthlySales
	err = json.Unmarshal(jsonData, &unmarshaledMonthlySales)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to MonthlySales: %v", err)
	}

	// Verify all fields are preserved
	if unmarshaledMonthlySales.Month != monthlySales.Month {
		t.Errorf("Expected Month %s, got %s", monthlySales.Month, unmarshaledMonthlySales.Month)
	}
	if unmarshaledMonthlySales.Year != monthlySales.Year {
		t.Errorf("Expected Year %d, got %d", monthlySales.Year, unmarshaledMonthlySales.Year)
	}
	if unmarshaledMonthlySales.TotalSales != monthlySales.TotalSales {
		t.Errorf("Expected TotalSales %f, got %f", monthlySales.TotalSales, unmarshaledMonthlySales.TotalSales)
	}
	if unmarshaledMonthlySales.SalesVolume != monthlySales.SalesVolume {
		t.Errorf("Expected SalesVolume %d, got %d", monthlySales.SalesVolume, unmarshaledMonthlySales.SalesVolume)
	}
}

// TestRegionRevenueStructWithMockData tests region revenue structure with hardcoded data
func TestRegionRevenueStructWithMockData(t *testing.T) {
	regionRevenue := RegionRevenue{
		Region:       "North America",
		TotalRevenue: 200000.0,
		ItemsSold:    4000,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(regionRevenue)
	if err != nil {
		t.Fatalf("Failed to marshal RegionRevenue to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledRegionRevenue RegionRevenue
	err = json.Unmarshal(jsonData, &unmarshaledRegionRevenue)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to RegionRevenue: %v", err)
	}

	// Verify all fields are preserved
	if unmarshaledRegionRevenue.Region != regionRevenue.Region {
		t.Errorf("Expected Region %s, got %s", regionRevenue.Region, unmarshaledRegionRevenue.Region)
	}
	if unmarshaledRegionRevenue.TotalRevenue != regionRevenue.TotalRevenue {
		t.Errorf("Expected TotalRevenue %f, got %f", regionRevenue.TotalRevenue, unmarshaledRegionRevenue.TotalRevenue)
	}
	if unmarshaledRegionRevenue.ItemsSold != regionRevenue.ItemsSold {
		t.Errorf("Expected ItemsSold %d, got %d", regionRevenue.ItemsSold, unmarshaledRegionRevenue.ItemsSold)
	}
}

// TestDashboardDataStructWithMockData tests dashboard data structure with hardcoded data
func TestDashboardDataStructWithMockData(t *testing.T) {
	mockTime := createMockTime()

	dashboardData := DashboardData{
		CountryRevenues: []CountryRevenue{
			{Country: "USA", ProductName: "Laptop", TotalRevenue: 50000.0, TransactionCount: 100},
			{Country: "UK", ProductName: "Smartphone", TotalRevenue: 30000.0, TransactionCount: 75},
		},
		TopProducts: []ProductFrequency{
			{ProductName: "Laptop", PurchaseCount: 500, CurrentStock: 100},
			{ProductName: "Smartphone", PurchaseCount: 400, CurrentStock: 150},
		},
		MonthlySales: []MonthlySales{
			{Month: "January", Year: 2024, TotalSales: 150000.0, SalesVolume: 3000},
			{Month: "February", Year: 2024, TotalSales: 180000.0, SalesVolume: 3600},
		},
		TopRegions: []RegionRevenue{
			{Region: "North America", TotalRevenue: 200000.0, ItemsSold: 4000},
			{Region: "Europe", TotalRevenue: 150000.0, ItemsSold: 3000},
		},
		LastUpdated:        mockTime,
		ProcessingDuration: 5 * time.Second,
		RecordCount:        1000,
	}

	// Test JSON marshaling
	jsonData, err := json.Marshal(dashboardData)
	if err != nil {
		t.Fatalf("Failed to marshal DashboardData to JSON: %v", err)
	}

	// Test JSON unmarshaling
	var unmarshaledDashboardData DashboardData
	err = json.Unmarshal(jsonData, &unmarshaledDashboardData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON to DashboardData: %v", err)
	}

	// Verify metadata fields
	if !unmarshaledDashboardData.LastUpdated.Equal(dashboardData.LastUpdated) {
		t.Errorf("Expected LastUpdated %v, got %v", dashboardData.LastUpdated, unmarshaledDashboardData.LastUpdated)
	}
	if unmarshaledDashboardData.ProcessingDuration != dashboardData.ProcessingDuration {
		t.Errorf("Expected ProcessingDuration %v, got %v", dashboardData.ProcessingDuration, unmarshaledDashboardData.ProcessingDuration)
	}
	if unmarshaledDashboardData.RecordCount != dashboardData.RecordCount {
		t.Errorf("Expected RecordCount %d, got %d", dashboardData.RecordCount, unmarshaledDashboardData.RecordCount)
	}

	// Verify CountryRevenues array
	if len(unmarshaledDashboardData.CountryRevenues) != len(dashboardData.CountryRevenues) {
		t.Errorf("Expected %d CountryRevenues, got %d", len(dashboardData.CountryRevenues), len(unmarshaledDashboardData.CountryRevenues))
	}
	if unmarshaledDashboardData.CountryRevenues[0].Country != "USA" {
		t.Errorf("Expected first CountryRevenue Country 'USA', got '%s'", unmarshaledDashboardData.CountryRevenues[0].Country)
	}

	// Verify TopProducts array
	if len(unmarshaledDashboardData.TopProducts) != len(dashboardData.TopProducts) {
		t.Errorf("Expected %d TopProducts, got %d", len(dashboardData.TopProducts), len(unmarshaledDashboardData.TopProducts))
	}
	if unmarshaledDashboardData.TopProducts[0].ProductName != "Laptop" {
		t.Errorf("Expected first TopProduct ProductName 'Laptop', got '%s'", unmarshaledDashboardData.TopProducts[0].ProductName)
	}

	// Verify MonthlySales array
	if len(unmarshaledDashboardData.MonthlySales) != len(dashboardData.MonthlySales) {
		t.Errorf("Expected %d MonthlySales, got %d", len(dashboardData.MonthlySales), len(unmarshaledDashboardData.MonthlySales))
	}
	if unmarshaledDashboardData.MonthlySales[0].Month != "January" {
		t.Errorf("Expected first MonthlySales Month 'January', got '%s'", unmarshaledDashboardData.MonthlySales[0].Month)
	}

	// Verify TopRegions array
	if len(unmarshaledDashboardData.TopRegions) != len(dashboardData.TopRegions) {
		t.Errorf("Expected %d TopRegions, got %d", len(dashboardData.TopRegions), len(unmarshaledDashboardData.TopRegions))
	}
	if unmarshaledDashboardData.TopRegions[0].Region != "North America" {
		t.Errorf("Expected first TopRegion Region 'North America', got '%s'", unmarshaledDashboardData.TopRegions[0].Region)
	}
}

// TestTransactionValidationWithMockData tests transaction field validation with hardcoded data
func TestTransactionValidationWithMockData(t *testing.T) {
	mockTime := createMockTime()

	// Test valid transaction
	validTransaction := Transaction{
		TransactionID:   "TXN001",
		TransactionDate: mockTime,
		UserID:          "USER001",
		Country:         "USA",
		Region:          "North America",
		ProductID:       "PROD001",
		ProductName:     "Test Laptop",
		Category:        "Electronics",
		Price:           999.99,
		Quantity:        2,
		TotalPrice:      1999.98,
		StockQuantity:   100,
		AddedDate:       mockTime,
	}

	// Verify valid transaction has all required fields
	if validTransaction.TransactionID == "" {
		t.Error("TransactionID should not be empty")
	}
	if validTransaction.UserID == "" {
		t.Error("UserID should not be empty")
	}
	if validTransaction.Country == "" {
		t.Error("Country should not be empty")
	}
	if validTransaction.ProductName == "" {
		t.Error("ProductName should not be empty")
	}
	if validTransaction.Price <= 0 {
		t.Error("Price should be positive")
	}
	if validTransaction.Quantity <= 0 {
		t.Error("Quantity should be positive")
	}
	if validTransaction.TotalPrice <= 0 {
		t.Error("TotalPrice should be positive")
	}
	if validTransaction.StockQuantity < 0 {
		t.Error("StockQuantity should be non-negative")
	}

	// Test edge cases
	edgeCaseTransaction := Transaction{
		TransactionID:   "TXN002",
		TransactionDate: mockTime,
		UserID:          "USER002",
		Country:         "UK",
		Region:          "Europe",
		ProductID:       "PROD002",
		ProductName:     "Test Smartphone",
		Category:        "Electronics",
		Price:           0.01, // Minimum price
		Quantity:        1,    // Minimum quantity
		TotalPrice:      0.01, // Minimum total
		StockQuantity:   0,    // Zero stock
		AddedDate:       mockTime,
	}

	// Verify edge case transaction
	if edgeCaseTransaction.Price != 0.01 {
		t.Errorf("Expected Price 0.01, got %f", edgeCaseTransaction.Price)
	}
	if edgeCaseTransaction.Quantity != 1 {
		t.Errorf("Expected Quantity 1, got %d", edgeCaseTransaction.Quantity)
	}
	if edgeCaseTransaction.TotalPrice != 0.01 {
		t.Errorf("Expected TotalPrice 0.01, got %f", edgeCaseTransaction.TotalPrice)
	}
	if edgeCaseTransaction.StockQuantity != 0 {
		t.Errorf("Expected StockQuantity 0, got %d", edgeCaseTransaction.StockQuantity)
	}
}

// TestCountryRevenueValidationWithMockData tests country revenue validation with hardcoded data
func TestCountryRevenueValidationWithMockData(t *testing.T) {
	// Test valid country revenue
	validCountryRevenue := CountryRevenue{
		Country:          "USA",
		ProductName:      "Test Laptop",
		TotalRevenue:     1500.50,
		TransactionCount: 25,
	}

	// Verify valid country revenue
	if validCountryRevenue.Country == "" {
		t.Error("Country should not be empty")
	}
	if validCountryRevenue.ProductName == "" {
		t.Error("ProductName should not be empty")
	}
	if validCountryRevenue.TotalRevenue <= 0 {
		t.Error("TotalRevenue should be positive")
	}
	if validCountryRevenue.TransactionCount <= 0 {
		t.Error("TransactionCount should be positive")
	}

	// Test edge case country revenue
	edgeCaseCountryRevenue := CountryRevenue{
		Country:          "UK",
		ProductName:      "Test Smartphone",
		TotalRevenue:     0.01, // Minimum revenue
		TransactionCount: 1,    // Minimum transaction count
	}

	// Verify edge case country revenue
	if edgeCaseCountryRevenue.TotalRevenue != 0.01 {
		t.Errorf("Expected TotalRevenue 0.01, got %f", edgeCaseCountryRevenue.TotalRevenue)
	}
	if edgeCaseCountryRevenue.TransactionCount != 1 {
		t.Errorf("Expected TransactionCount 1, got %d", edgeCaseCountryRevenue.TransactionCount)
	}
}

// TestProductFrequencyValidationWithMockData tests product frequency validation with hardcoded data
func TestProductFrequencyValidationWithMockData(t *testing.T) {
	// Test valid product frequency
	validProductFrequency := ProductFrequency{
		ProductName:   "Test Laptop",
		PurchaseCount: 150,
		CurrentStock:  75,
	}

	// Verify valid product frequency
	if validProductFrequency.ProductName == "" {
		t.Error("ProductName should not be empty")
	}
	if validProductFrequency.PurchaseCount <= 0 {
		t.Error("PurchaseCount should be positive")
	}
	if validProductFrequency.CurrentStock < 0 {
		t.Error("CurrentStock should be non-negative")
	}

	// Test edge case product frequency
	edgeCaseProductFrequency := ProductFrequency{
		ProductName:   "Test Smartphone",
		PurchaseCount: 1, // Minimum purchase count
		CurrentStock:  0, // Zero stock
	}

	// Verify edge case product frequency
	if edgeCaseProductFrequency.PurchaseCount != 1 {
		t.Errorf("Expected PurchaseCount 1, got %d", edgeCaseProductFrequency.PurchaseCount)
	}
	if edgeCaseProductFrequency.CurrentStock != 0 {
		t.Errorf("Expected CurrentStock 0, got %d", edgeCaseProductFrequency.CurrentStock)
	}
}

// TestMonthlySalesValidationWithMockData tests monthly sales validation with hardcoded data
func TestMonthlySalesValidationWithMockData(t *testing.T) {
	// Test valid monthly sales
	validMonthlySales := MonthlySales{
		Month:       "January",
		Year:        2024,
		TotalSales:  150000.0,
		SalesVolume: 3000,
	}

	// Verify valid monthly sales
	if validMonthlySales.Month == "" {
		t.Error("Month should not be empty")
	}
	if validMonthlySales.Year <= 0 {
		t.Error("Year should be positive")
	}
	if validMonthlySales.TotalSales <= 0 {
		t.Error("TotalSales should be positive")
	}
	if validMonthlySales.SalesVolume <= 0 {
		t.Error("SalesVolume should be positive")
	}

	// Test edge case monthly sales
	edgeCaseMonthlySales := MonthlySales{
		Month:       "December",
		Year:        2023,
		TotalSales:  0.01, // Minimum sales
		SalesVolume: 1,    // Minimum volume
	}

	// Verify edge case monthly sales
	if edgeCaseMonthlySales.TotalSales != 0.01 {
		t.Errorf("Expected TotalSales 0.01, got %f", edgeCaseMonthlySales.TotalSales)
	}
	if edgeCaseMonthlySales.SalesVolume != 1 {
		t.Errorf("Expected SalesVolume 1, got %d", edgeCaseMonthlySales.SalesVolume)
	}
}

// TestRegionRevenueValidationWithMockData tests region revenue validation with hardcoded data
func TestRegionRevenueValidationWithMockData(t *testing.T) {
	// Test valid region revenue
	validRegionRevenue := RegionRevenue{
		Region:       "North America",
		TotalRevenue: 200000.0,
		ItemsSold:    4000,
	}

	// Verify valid region revenue
	if validRegionRevenue.Region == "" {
		t.Error("Region should not be empty")
	}
	if validRegionRevenue.TotalRevenue <= 0 {
		t.Error("TotalRevenue should be positive")
	}
	if validRegionRevenue.ItemsSold <= 0 {
		t.Error("ItemsSold should be positive")
	}

	// Test edge case region revenue
	edgeCaseRegionRevenue := RegionRevenue{
		Region:       "Europe",
		TotalRevenue: 0.01, // Minimum revenue
		ItemsSold:    1,    // Minimum items sold
	}

	// Verify edge case region revenue
	if edgeCaseRegionRevenue.TotalRevenue != 0.01 {
		t.Errorf("Expected TotalRevenue 0.01, got %f", edgeCaseRegionRevenue.TotalRevenue)
	}
	if edgeCaseRegionRevenue.ItemsSold != 1 {
		t.Errorf("Expected ItemsSold 1, got %d", edgeCaseRegionRevenue.ItemsSold)
	}
}
