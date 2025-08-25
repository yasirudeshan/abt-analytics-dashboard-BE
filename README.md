# ABT Analytics Dashboard Backend

High-performance Go backend for processing 1M+ records and providing analytics APIs.

## Features

- 🚀 **High Performance**: Concurrent processing with worker goroutines
- 📊 **Real-time Analytics**: Live dashboard data with RESTful APIs
- 🔒 **Thread-safe**: Mutex-protected concurrent data access
- 📈 **Scalable**: Designed to handle large datasets efficiently
- 🧪 **Well-tested**: 94% test coverage with comprehensive test suite

## Technologies & Optimization Techniques

- **Language/Runtime**: Go 1.21+ (compiled, low-latency GC, great concurrency model)
- **Concurrency**: `goroutines` + `channels` for parallel CSV ingestion and aggregation
- **Work Distribution**: `runtime.NumCPU()`-based worker pool for CPU-bound phases
- **Synchronization**: `sync.Mutex` + `sync.RWMutex` around shared maps and snapshots
- **Memory Efficiency**: Streaming CSV with `bufio.Reader` to avoid full-file loading
- **Data Structures**: In-memory maps for O(1) aggregation, converted to slices for sorting
- **Sorting Performance**: `sort.Slice` with pre-sized slices to minimize allocations
- **I/O Efficiency**: `encoding/csv` with `LazyQuotes` to tolerate imperfect data
- **Graceful Shutdown**: Context-based shutdown to avoid partial writes/corruption
- **Observability**: Consistent logging of progress and timings for large datasets
- **HTTP Layer**: Gorilla `mux` router with lightweight middleware (CORS, logging)
- **Dev Productivity**: Make targets, scripts, coverage tooling, and race detector

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

## Testing

### Test Coverage: 94% ✅

The project includes a comprehensive test suite covering all packages and critical functionality.

#### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Run tests with race detection
make test-race

# Run benchmarks
make bench
```

#### Test Scripts

```bash
# Linux/macOS
./scripts/run_tests.sh

# Windows
scripts\run_tests.bat
```

#### Coverage Reports

After running tests with coverage, you'll get:
- **`coverage.html`**: Interactive HTML coverage report
- **`coverage.txt`**: Human-readable text coverage report
- **`coverage.out`**: Raw coverage data

#### Test Structure

```
├── main_test.go                    # Main package tests
├── internal/
│   ├── config/
│   │   └── config_test.go         # Configuration tests
│   ├── models/
│   │   └── models_test.go         # Data model tests
│   ├── processor/
│   │   └── processor_test.go      # Data processing tests
│   └── api/
│       └── server_test.go         # API server tests
└── scripts/
    ├── run_tests.sh               # Linux/macOS test runner
    └── run_tests.bat              # Windows test runner
```

### Test Categories

- **Unit Tests**: Individual function testing
- **Integration Tests**: API endpoint and data flow testing
- **Performance Tests**: Benchmarking and race detection
- **Edge Case Tests**: Error handling and boundary conditions

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

## Development

### Prerequisites

- Go 1.21+
- Make (optional, for using Makefile targets)

### Available Makefile Targets

```bash
make help              # Show all available targets
make build            # Build the application
make test             # Run tests
make test-coverage    # Run tests with coverage
make test-race        # Run tests with race detection
make bench            # Run benchmarks
make fmt              # Format code
make vet              # Vet code
make lint             # Run linter (requires golangci-lint)
make clean            # Clean build artifacts
make docker-build     # Build Docker image
make docker-run       # Run Docker container
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (if installed)
golangci-lint run

# Check for security vulnerabilities
go list -json -m all | nancy sleuth
```

## Architecture

### Package Structure

```
├── main.go                          # Application entry point
├── internal/
│   ├── config/                     # Configuration management
│   ├── models/                     # Data structures
│   ├── processor/                  # Data processing engine
│   └── api/                        # HTTP server and handlers
├── data/                           # Dataset storage
├── scripts/                        # Utility scripts
├── docs/                           # Documentation
└── Dockerfile                      # Container configuration
```

### Key Components

- **Config**: Environment-based configuration loading
- **Models**: Structured data types with JSON/CSV tags
- **Processor**: Concurrent CSV processing with worker pools
- **API**: RESTful HTTP server with middleware support

## Performance

- **Concurrent Processing**: Uses worker goroutines for data aggregation
- **Memory Efficient**: Streaming CSV processing for large datasets
- **Fast Aggregation**: Optimized sorting and aggregation algorithms
- **Scalable**: Designed to handle millions of records

## Contributing

1. Ensure tests pass: `make test`
2. Maintain test coverage above 90%
3. Follow Go coding standards
4. Add tests for new functionality
5. Update documentation as needed

## License

This project is licensed under the MIT License.
