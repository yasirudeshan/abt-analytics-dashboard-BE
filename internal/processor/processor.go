package processor

import (
	"abt-analytics-dashboard/internal/models"
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Processor handles data processing and aggregation
type Processor struct {
	dashboardData *models.DashboardData
	mu            sync.RWMutex
}

// New creates a new processor instance
func New() *Processor {
	return &Processor{
		dashboardData: &models.DashboardData{
			CountryRevenues: make([]models.CountryRevenue, 0),
			TopProducts:     make([]models.ProductFrequency, 0),
			MonthlySales:    make([]models.MonthlySales, 0),
			TopRegions:      make([]models.RegionRevenue, 0),
		},
	}
}

// ProcessDataset processes the CSV dataset using concurrent workers
func (p *Processor) ProcessDataset(filePath string) error {
	start := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Create channels for concurrent processing
	transactionCh := make(chan models.Transaction, 1000)
	errorCh := make(chan error, 1)
	done := make(chan struct{})

	// Start aggregation workers
	numWorkers := runtime.NumCPU()
	log.Printf("Starting %d worker goroutines for data processing", numWorkers)

	// Aggregation maps with mutexes for concurrent access
	countryMap := make(map[string]*models.CountryRevenue)
	productMap := make(map[string]*models.ProductFrequency)
	monthMap := make(map[string]*models.MonthlySales)
	regionMap := make(map[string]*models.RegionRevenue)

	var wg sync.WaitGroup
	var mu sync.Mutex

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.aggregateWorker(transactionCh, &mu, countryMap, productMap, monthMap, regionMap)
		}()
	}

	// Start CSV reader goroutine
	go func() {
		defer close(transactionCh)
		if err := p.readCSV(file, transactionCh); err != nil {
			errorCh <- err
			return
		}
	}()

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(done)
	}()

	// Wait for completion or error
	select {
	case err := <-errorCh:
		return fmt.Errorf("error during processing: %w", err)
	case <-done:
		// Processing completed successfully
	}

	// Convert maps to sorted slices and store in dashboard data
	p.mu.Lock()
	p.dashboardData.CountryRevenues = p.sortCountryRevenues(countryMap)
	p.dashboardData.TopProducts = p.sortTopProducts(productMap, 20)
	p.dashboardData.MonthlySales = p.sortMonthlySales(monthMap)
	p.dashboardData.TopRegions = p.sortTopRegions(regionMap, 30)
	p.dashboardData.LastUpdated = time.Now()
	p.dashboardData.ProcessingDuration = time.Since(start)
	p.dashboardData.RecordCount = len(countryMap) // Approximate record count
	p.mu.Unlock()

	log.Printf("Data processing completed in %v", time.Since(start))
	return nil
}

// readCSV reads CSV file and sends transactions to channel
func (p *Processor) readCSV(file *os.File, transactionCh chan<- models.Transaction) error {
	reader := csv.NewReader(bufio.NewReader(file))
	reader.LazyQuotes = true

	// Read header
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Map headers to indices
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[strings.TrimSpace(strings.ToLower(header))] = i
	}

	recordCount := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading record %d: %v", recordCount, err)
			continue
		}

		transaction, err := p.parseTransaction(record, headerMap)
		if err != nil {
			log.Printf("Error parsing record %d: %v", recordCount, err)
			continue
		}

		transactionCh <- transaction
		recordCount++

		// Log progress for large datasets
		if recordCount%100000 == 0 {
			log.Printf("Processed %d records", recordCount)
		}
	}

	log.Printf("Finished reading %d records from CSV", recordCount)
	return nil
}

// parseTransaction parses a CSV record into a Transaction struct
func (p *Processor) parseTransaction(record []string, headerMap map[string]int) (models.Transaction, error) {
	var transaction models.Transaction

	if idx, ok := headerMap["transaction_id"]; ok && idx < len(record) {
		transaction.TransactionID = strings.TrimSpace(record[idx])
	}
	if idx, ok := headerMap["user_id"]; ok && idx < len(record) {
		transaction.UserID = strings.TrimSpace(record[idx])
	}
	if idx, ok := headerMap["product_id"]; ok && idx < len(record) {
		transaction.ProductID = strings.TrimSpace(record[idx])
	}
	if idx, ok := headerMap["product_name"]; ok && idx < len(record) {
		transaction.ProductName = strings.TrimSpace(record[idx])
	}
	if idx, ok := headerMap["category"]; ok && idx < len(record) {
		transaction.Category = strings.TrimSpace(record[idx])
	}
	if idx, ok := headerMap["country"]; ok && idx < len(record) {
		transaction.Country = strings.TrimSpace(record[idx])
	}
	if idx, ok := headerMap["region"]; ok && idx < len(record) {
		transaction.Region = strings.TrimSpace(record[idx])
	}

	if idx, ok := headerMap["price"]; ok && idx < len(record) {
		if price, err := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); err == nil {
			transaction.Price = price
		}
	}

	if idx, ok := headerMap["total_price"]; ok && idx < len(record) {
		if totalPrice, err := strconv.ParseFloat(strings.TrimSpace(record[idx]), 64); err == nil {
			transaction.TotalPrice = totalPrice
		}
	}

	if idx, ok := headerMap["quantity"]; ok && idx < len(record) {
		if quantity, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
			transaction.Quantity = quantity
		}
	}

	if idx, ok := headerMap["stock_quantity"]; ok && idx < len(record) {
		if stock, err := strconv.Atoi(strings.TrimSpace(record[idx])); err == nil {
			transaction.StockQuantity = stock
		}
	}

	// Parse transaction_date
	if idx, ok := headerMap["transaction_date"]; ok && idx < len(record) {
		dateStr := strings.TrimSpace(record[idx])
		// Try multiple date formats
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			"01/02/2006",
			"01-02-2006",
			"2006/01/02",
		}
		for _, format := range formats {
			if date, err := time.Parse(format, dateStr); err == nil {
				transaction.TransactionDate = date
				break
			}
		}
	}

	// Parse added_date
	if idx, ok := headerMap["added_date"]; ok && idx < len(record) {
		dateStr := strings.TrimSpace(record[idx])
		// Try multiple date formats
		formats := []string{
			"2006-01-02",
			"2006-01-02 15:04:05",
			"01/02/2006",
			"01-02-2006",
			"2006/01/02",
		}
		for _, format := range formats {
			if date, err := time.Parse(format, dateStr); err == nil {
				transaction.AddedDate = date
				break
			}
		}
	}

	return transaction, nil
}

// aggregateWorker processes transactions and updates aggregation maps
func (p *Processor) aggregateWorker(
	transactionCh <-chan models.Transaction,
	mu *sync.Mutex,
	countryMap map[string]*models.CountryRevenue,
	productMap map[string]*models.ProductFrequency,
	monthMap map[string]*models.MonthlySales,
	regionMap map[string]*models.RegionRevenue,
) {
	for transaction := range transactionCh {
		mu.Lock()

		// Aggregate country revenue
		countryKey := fmt.Sprintf("%s-%s", transaction.Country, transaction.ProductName)
		if countryRev, exists := countryMap[countryKey]; exists {
			countryRev.TotalRevenue += transaction.TotalPrice
			countryRev.TransactionCount++
		} else {
			countryMap[countryKey] = &models.CountryRevenue{
				Country:          transaction.Country,
				ProductName:      transaction.ProductName,
				TotalRevenue:     transaction.TotalPrice,
				TransactionCount: 1,
			}
		}

		// Aggregate product frequency
		if product, exists := productMap[transaction.ProductName]; exists {
			product.PurchaseCount++
			if transaction.StockQuantity > 0 {
				product.CurrentStock = transaction.StockQuantity // Keep latest stock value
			}
		} else {
			productMap[transaction.ProductName] = &models.ProductFrequency{
				ProductName:   transaction.ProductName,
				PurchaseCount: 1,
				CurrentStock:  transaction.StockQuantity,
			}
		}

		// Aggregate monthly sales (use transaction_date)
		monthKey := fmt.Sprintf("%d-%02d", transaction.TransactionDate.Year(), transaction.TransactionDate.Month())
		if monthlySales, exists := monthMap[monthKey]; exists {
			monthlySales.TotalSales += transaction.TotalPrice
			monthlySales.SalesVolume += transaction.Quantity
		} else {
			monthMap[monthKey] = &models.MonthlySales{
				Month:       transaction.TransactionDate.Format("January"),
				Year:        transaction.TransactionDate.Year(),
				TotalSales:  transaction.TotalPrice,
				SalesVolume: transaction.Quantity,
			}
		}

		// Aggregate region revenue
		if region, exists := regionMap[transaction.Region]; exists {
			region.TotalRevenue += transaction.TotalPrice
			region.ItemsSold += transaction.Quantity
		} else {
			regionMap[transaction.Region] = &models.RegionRevenue{
				Region:       transaction.Region,
				TotalRevenue: transaction.TotalPrice,
				ItemsSold:    transaction.Quantity,
			}
		}

		mu.Unlock()
	}
}

// Sorting functions
func (p *Processor) sortCountryRevenues(countryMap map[string]*models.CountryRevenue) []models.CountryRevenue {
	revenues := make([]models.CountryRevenue, 0, len(countryMap))
	for _, rev := range countryMap {
		revenues = append(revenues, *rev)
	}

	sort.Slice(revenues, func(i, j int) bool {
		return revenues[i].TotalRevenue > revenues[j].TotalRevenue
	})

	return revenues
}

func (p *Processor) sortTopProducts(productMap map[string]*models.ProductFrequency, limit int) []models.ProductFrequency {
	products := make([]models.ProductFrequency, 0, len(productMap))
	for _, product := range productMap {
		products = append(products, *product)
	}

	sort.Slice(products, func(i, j int) bool {
		return products[i].PurchaseCount > products[j].PurchaseCount
	})

	if len(products) > limit {
		products = products[:limit]
	}

	return products
}

func (p *Processor) sortMonthlySales(monthMap map[string]*models.MonthlySales) []models.MonthlySales {
	sales := make([]models.MonthlySales, 0, len(monthMap))
	for _, sale := range monthMap {
		sales = append(sales, *sale)
	}

	sort.Slice(sales, func(i, j int) bool {
		if sales[i].Year != sales[j].Year {
			return sales[i].Year > sales[j].Year
		}
		return sales[i].TotalSales > sales[j].TotalSales
	})

	return sales
}

func (p *Processor) sortTopRegions(regionMap map[string]*models.RegionRevenue, limit int) []models.RegionRevenue {
	regions := make([]models.RegionRevenue, 0, len(regionMap))
	for _, region := range regionMap {
		regions = append(regions, *region)
	}

	sort.Slice(regions, func(i, j int) bool {
		return regions[i].TotalRevenue > regions[j].TotalRevenue
	})

	if len(regions) > limit {
		regions = regions[:limit]
	}

	return regions
}

// GetDashboardData returns the current dashboard data (thread-safe)
func (p *Processor) GetDashboardData() *models.DashboardData {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dashboardData
}

// GetCountryRevenues returns country revenue data
func (p *Processor) GetCountryRevenues() []models.CountryRevenue {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dashboardData.CountryRevenues
}

// GetTopProducts returns top products data
func (p *Processor) GetTopProducts() []models.ProductFrequency {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dashboardData.TopProducts
}

// GetMonthlySales returns monthly sales data
func (p *Processor) GetMonthlySales() []models.MonthlySales {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dashboardData.MonthlySales
}

// GetTopRegions returns top regions data
func (p *Processor) GetTopRegions() []models.RegionRevenue {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.dashboardData.TopRegions
}
