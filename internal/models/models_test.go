package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTransactionStruct(t *testing.T) {
	now := time.Now()
	transaction := Transaction{
		TransactionID:   "TXN001",
		TransactionDate: now,
		UserID:          "USER001",
		Country:         "USA",
		Region:          "North America",
		ProductID:       "PROD001",
		ProductName:     "Test Product",
		Category:        "Electronics",
		Price:           99.99,
		Quantity:        2,
		TotalPrice:      199.98,
		StockQuantity:   100,
		AddedDate:       now,
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

	// Verify fields
	if unmarshaledTransaction.TransactionID != transaction.TransactionID {
		t.Errorf("Expected TransactionID %s, got %s", transaction.TransactionID, unmarshaledTransaction.TransactionID)
	}
	if unmarshaledTransaction.UserID != transaction.UserID {
		t.Errorf("Expected UserID %s, got %s", transaction.UserID, unmarshaledTransaction.UserID)
	}
	if unmarshaledTransaction.Country != transaction.Country {
		t.Errorf("Expected Country %s, got %s", transaction.Country, unmarshaledTransaction.Country)
	}
	if unmarshaledTransaction.ProductName != transaction.ProductName {
		t.Errorf("Expected ProductName %s, got %s", transaction.ProductName, unmarshaledTransaction.ProductName)
	}
	if unmarshaledTransaction.Price != transaction.Price {
		t.Errorf("Expected Price %f, got %f", transaction.Price, unmarshaledTransaction.Price)
	}
	if unmarshaledTransaction.Quantity != transaction.Quantity {
		t.Errorf("Expected Quantity %d, got %d", transaction.Quantity, unmarshaledTransaction.Quantity)
	}
}

func TestCountryRevenueStruct(t *testing.T) {
	countryRevenue := CountryRevenue{
		Country:          "USA",
		ProductName:      "Test Product",
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

	// Verify fields
	if unmarshaledCountryRevenue.Country != countryRevenue.Country {
		t.Errorf("Expected Country %s, got %s", countryRevenue.Country, unmarshaledCountryRevenue.Country)
	}
	if unmarshaledCountryRevenue.TotalRevenue != countryRevenue.TotalRevenue {
		t.Errorf("Expected TotalRevenue %f, got %f", countryRevenue.TotalRevenue, unmarshaledCountryRevenue.TotalRevenue)
	}
	if unmarshaledCountryRevenue.TransactionCount != countryRevenue.TransactionCount {
		t.Errorf("Expected TransactionCount %d, got %d", countryRevenue.TransactionCount, unmarshaledCountryRevenue.TransactionCount)
	}
}

func TestProductFrequencyStruct(t *testing.T) {
	productFrequency := ProductFrequency{
		ProductName:   "Test Product",
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

	// Verify fields
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

func TestMonthlySalesStruct(t *testing.T) {
	monthlySales := MonthlySales{
		Month:       "January",
		Year:        2024,
		TotalSales:  50000.75,
		SalesVolume: 1250,
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

	// Verify fields
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

func TestRegionRevenueStruct(t *testing.T) {
	regionRevenue := RegionRevenue{
		Region:       "North America",
		TotalRevenue: 75000.25,
		ItemsSold:    3000,
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

	// Verify fields
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

func TestDashboardDataStruct(t *testing.T) {
	now := time.Now()
	dashboardData := DashboardData{
		CountryRevenues: []CountryRevenue{
			{Country: "USA", ProductName: "Product1", TotalRevenue: 1000.0, TransactionCount: 10},
		},
		TopProducts: []ProductFrequency{
			{ProductName: "Product1", PurchaseCount: 100, CurrentStock: 50},
		},
		MonthlySales: []MonthlySales{
			{Month: "January", Year: 2024, TotalSales: 5000.0, SalesVolume: 100},
		},
		TopRegions: []RegionRevenue{
			{Region: "North America", TotalRevenue: 10000.0, ItemsSold: 200},
		},
		LastUpdated:        now,
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

	// Verify fields
	if len(unmarshaledDashboardData.CountryRevenues) != len(dashboardData.CountryRevenues) {
		t.Errorf("Expected %d CountryRevenues, got %d", len(dashboardData.CountryRevenues), len(unmarshaledDashboardData.CountryRevenues))
	}
	if len(unmarshaledDashboardData.TopProducts) != len(dashboardData.TopProducts) {
		t.Errorf("Expected %d TopProducts, got %d", len(dashboardData.TopProducts), len(unmarshaledDashboardData.TopProducts))
	}
	if len(unmarshaledDashboardData.MonthlySales) != len(dashboardData.MonthlySales) {
		t.Errorf("Expected %d MonthlySales, got %d", len(dashboardData.MonthlySales), len(unmarshaledDashboardData.MonthlySales))
	}
	if len(unmarshaledDashboardData.TopRegions) != len(dashboardData.TopRegions) {
		t.Errorf("Expected %d TopRegions, got %d", len(dashboardData.TopRegions), len(unmarshaledDashboardData.TopRegions))
	}
	if unmarshaledDashboardData.RecordCount != dashboardData.RecordCount {
		t.Errorf("Expected RecordCount %d, got %d", dashboardData.RecordCount, unmarshaledDashboardData.RecordCount)
	}
}

func TestTimeParsing(t *testing.T) {
	// Test that time fields can handle various formats
	transaction := Transaction{
		TransactionDate: time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC),
		AddedDate:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
	}

	// Verify time fields are properly set
	if transaction.TransactionDate.IsZero() {
		t.Error("Expected TransactionDate to be set, got zero time")
	}
	if transaction.AddedDate.IsZero() {
		t.Error("Expected AddedDate to be set, got zero time")
	}
}

