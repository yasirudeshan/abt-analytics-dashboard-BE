@echo off
echo Starting ABT Analytics Dashboard Backend...
echo.

REM Check if GO_test_5m.csv exists in common locations
if exist "data\GO_test_5m.csv" (
    echo Found dataset: data\GO_test_5m.csv
    set DATA_FILE_PATH=data\GO_test_5m.csv
) else if exist "GO_test_5m.csv" (
    echo Found dataset: GO_test_5m.csv
    set DATA_FILE_PATH=GO_test_5m.csv
) else (
    echo No dataset found. Using sample data for demonstration.
    echo.
    echo To use your dataset, place GO_test_5m.csv in:
    echo   - data\GO_test_5m.csv  (recommended)
    echo   - GO_test_5m.csv       (root directory)
    echo.
    echo Or set DATA_FILE_PATH environment variable to your file location.
)

echo.
echo Server will be available at: http://localhost:8080
echo API endpoints:
echo   - Health Check: http://localhost:8080/api/health
echo   - Dashboard Data: http://localhost:8080/api/dashboard
echo   - Country Revenues: http://localhost:8080/api/revenue-by-country
echo   - Top Products: http://localhost:8080/api/top-products
echo   - Monthly Sales: http://localhost:8080/api/sales-by-month
echo   - Top Regions: http://localhost:8080/api/top-regions
echo.
echo Press Ctrl+C to stop the server
echo.

go run main.go
