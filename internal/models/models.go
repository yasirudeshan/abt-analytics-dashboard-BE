package models

import "time"

// Transaction represents a single transaction record
type Transaction struct {
	TransactionID   string    `json:"transaction_id" csv:"transaction_id"`
	TransactionDate time.Time `json:"transaction_date" csv:"transaction_date"`
	UserID          string    `json:"user_id" csv:"user_id"`
	Country         string    `json:"country" csv:"country"`
	Region          string    `json:"region" csv:"region"`
	ProductID       string    `json:"product_id" csv:"product_id"`
	ProductName     string    `json:"product_name" csv:"product_name"`
	Category        string    `json:"category" csv:"category"`
	Price           float64   `json:"price" csv:"price"`
	Quantity        int       `json:"quantity" csv:"quantity"`
	TotalPrice      float64   `json:"total_price" csv:"total_price"`
	StockQuantity   int       `json:"stock_quantity" csv:"stock_quantity"`
	AddedDate       time.Time `json:"added_date" csv:"added_date"`
}

// CountryRevenue represents country-level revenue data
type CountryRevenue struct {
	Country          string  `json:"country"`
	ProductName      string  `json:"product_name"`
	TotalRevenue     float64 `json:"total_revenue"`
	TransactionCount int     `json:"transaction_count"`
}

// ProductFrequency represents product purchase frequency data
type ProductFrequency struct {
	ProductName   string `json:"product_name"`
	PurchaseCount int    `json:"purchase_count"`
	CurrentStock  int    `json:"current_stock"`
}

// MonthlySales represents monthly sales volume data
type MonthlySales struct {
	Month       string  `json:"month"`
	Year        int     `json:"year"`
	TotalSales  float64 `json:"total_sales"`
	SalesVolume int     `json:"sales_volume"`
}

// RegionRevenue represents region-level revenue data
type RegionRevenue struct {
	Region       string  `json:"region"`
	TotalRevenue float64 `json:"total_revenue"`
	ItemsSold    int     `json:"items_sold"`
}

// DashboardData contains all pre-aggregated dashboard data
type DashboardData struct {
	CountryRevenues    []CountryRevenue   `json:"country_revenues"`
	TopProducts        []ProductFrequency `json:"top_products"`
	MonthlySales       []MonthlySales     `json:"monthly_sales"`
	TopRegions         []RegionRevenue    `json:"top_regions"`
	LastUpdated        time.Time          `json:"last_updated"`
	ProcessingDuration time.Duration      `json:"processing_duration"`
	RecordCount        int                `json:"record_count"`
}
