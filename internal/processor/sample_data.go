package processor

import (
	"abt-analytics-dashboard/internal/models"
	"math/rand"
	"time"
)

// LoadSampleData generates sample data for development and testing
func (p *Processor) LoadSampleData() {
	p.mu.Lock()
	defer p.mu.Unlock()

	start := time.Now()

	// Sample countries and regions
	countries := []string{"USA", "UK", "Germany", "France", "Japan", "Canada", "Australia", "Brazil", "India", "China"}
	regions := []string{"North America", "Europe", "Asia Pacific", "Latin America", "Middle East", "Africa"}
	products := []string{
		"Wireless Headphones", "Smartphone", "Laptop", "Tablet", "Smartwatch",
		"Camera", "Gaming Console", "Keyboard", "Mouse", "Monitor",
		"Speakers", "Microphone", "Webcam", "Router", "Hard Drive",
		"SSD", "Graphics Card", "Processor", "Memory", "Motherboard",
	}

	// Generate sample country revenues
	p.dashboardData.CountryRevenues = make([]models.CountryRevenue, 0)
	for _, country := range countries {
		for i, product := range products {
			if i > 5 && rand.Float32() > 0.7 { // Skip some combinations
				continue
			}
			revenue := models.CountryRevenue{
				Country:          country,
				ProductName:      product,
				TotalRevenue:     rand.Float64()*50000 + 10000, // $10k-$60k
				TransactionCount: rand.Intn(500) + 50,          // 50-550 transactions
			}
			p.dashboardData.CountryRevenues = append(p.dashboardData.CountryRevenues, revenue)
		}
	}

	// Generate sample top products
	p.dashboardData.TopProducts = make([]models.ProductFrequency, len(products))
	for i, product := range products {
		p.dashboardData.TopProducts[i] = models.ProductFrequency{
			ProductName:   product,
			PurchaseCount: rand.Intn(10000) + 1000, // 1000-11000 purchases
			CurrentStock:  rand.Intn(500) + 50,     // 50-550 stock
		}
	}

	// Generate sample monthly sales (last 12 months)
	p.dashboardData.MonthlySales = make([]models.MonthlySales, 12)
	months := []string{
		"January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	currentYear := time.Now().Year()
	for i, month := range months {
		p.dashboardData.MonthlySales[i] = models.MonthlySales{
			Month:       month,
			Year:        currentYear,
			TotalSales:  rand.Float64()*200000 + 100000, // $100k-$300k
			SalesVolume: rand.Intn(5000) + 2000,         // 2000-7000 items
		}
	}

	// Generate sample top regions
	p.dashboardData.TopRegions = make([]models.RegionRevenue, len(regions))
	for i, region := range regions {
		p.dashboardData.TopRegions[i] = models.RegionRevenue{
			Region:       region,
			TotalRevenue: rand.Float64()*500000 + 200000, // $200k-$700k
			ItemsSold:    rand.Intn(20000) + 5000,        // 5000-25000 items
		}
	}

	// Set metadata
	p.dashboardData.LastUpdated = time.Now()
	p.dashboardData.ProcessingDuration = time.Since(start)
	p.dashboardData.RecordCount = len(p.dashboardData.CountryRevenues)
}
