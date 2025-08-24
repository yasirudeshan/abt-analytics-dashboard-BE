# ABT Analytics Dashboard Backend

High-performance Go backend for processing 1M+ records and providing analytics APIs.

## Quick Start

### Development
```bash
# Place your GO_test_5m.csv in the data/ folder
go run main.go
```

### Production
```bash
# Build binary
go build -o app .

# Run binary
./app

# Or use Makefile
make build
make run
```

### Docker
```bash
make docker
docker run -p 8080:8080 abt-analytics-dashboard
```

## API Endpoints

- `GET /api/health` - Server status
- `GET /api/revenue-by-country` - Country revenue table  
- `GET /api/top-products` - Top 20 products
- `GET /api/sales-by-month` - Monthly sales
- `GET /api/top-regions` - Top 30 regions
- `GET /api/dashboard` - All data

## Dataset Format
CSV 
`transaction_id,transaction_date,user_id,country,region,product_id,product_name,category,price,quantity,total_price,stock_quantity,added_date`
