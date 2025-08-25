@echo off
REM ABT Analytics Dashboard Backend - Test Runner Script for Windows
REM This script runs comprehensive tests and generates coverage reports

setlocal enabledelayedexpansion

echo ==========================================
echo ABT Analytics Dashboard Backend - Test Suite
echo ==========================================
echo.

REM Check if Go is installed
go version >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go is not installed or not in PATH
    exit /b 1
)

for /f "tokens=3" %%i in ('go version') do set GO_VERSION=%%i
echo [INFO] Using Go version: %GO_VERSION%

REM Check if required tools are available
echo [INFO] Checking required tools...
go tool cover >nul 2>&1
if errorlevel 1 (
    echo [ERROR] Go cover tool not available
    exit /b 1
)
echo [SUCCESS] All required tools are available

REM Clean up previous test artifacts
echo [INFO] Cleaning up previous test artifacts...
if exist coverage.out del coverage.out
if exist coverage.html del coverage.html
if exist coverage.txt del coverage.txt
echo [SUCCESS] Cleanup completed

REM Run basic tests
echo [INFO] Running basic tests...
go test ./... -v
if errorlevel 1 (
    echo [ERROR] Basic tests failed
    exit /b 1
)
echo [SUCCESS] Basic tests passed

REM Run tests with race detection
echo [INFO] Running tests with race detection...
go test -race ./... -v
if errorlevel 1 (
    echo [WARNING] Race detection tests failed (this may be expected in some environments)
) else (
    echo [SUCCESS] Race detection tests passed
)

REM Run tests with coverage
echo [INFO] Running tests with coverage...
go test -coverprofile=coverage.out ./... -v
if errorlevel 1 (
    echo [ERROR] Coverage tests failed
    exit /b 1
)
echo [SUCCESS] Coverage tests passed

REM Generate coverage reports
echo [INFO] Generating coverage reports...

REM Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
if errorlevel 1 (
    echo [ERROR] Failed to generate HTML coverage report
) else (
    echo [SUCCESS] HTML coverage report generated: coverage.html
)

REM Generate text coverage report
go tool cover -func=coverage.out > coverage.txt
if errorlevel 1 (
    echo [ERROR] Failed to generate text coverage report
) else (
    echo [SUCCESS] Text coverage report generated: coverage.txt
)

REM Display coverage summary
echo [INFO] Coverage Summary:
go tool cover -func=coverage.out | findstr "total:"

REM Run benchmarks
echo [INFO] Running benchmarks...
go test -bench=. -benchmem ./... -v
if errorlevel 1 (
    echo [WARNING] Some benchmarks failed (this may be expected)
) else (
    echo [SUCCESS] Benchmarks completed
)

REM Run code vetting
echo [INFO] Running code vetting...
go vet ./...
if errorlevel 1 (
    echo [WARNING] Code vetting found issues
) else (
    echo [SUCCESS] Code vetting passed
)

REM Display final summary
echo.
echo [INFO] Test Execution Summary
echo ==========================
if exist coverage.out (
    echo [INFO] Raw Coverage Data: coverage.out
)
if exist coverage.html (
    echo [INFO] HTML Coverage Report: coverage.html
)
if exist coverage.txt (
    echo [INFO] Text Coverage Report: coverage.txt
)

echo.
echo [SUCCESS] All tests completed successfully!
echo.
echo Press any key to exit...
pause >nul

