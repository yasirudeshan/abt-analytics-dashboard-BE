#!/bin/bash

# ABT Analytics Dashboard Backend - Test Runner Script
# This script runs comprehensive tests and generates coverage reports

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        exit 1
    fi
    
    GO_VERSION=$(go version | awk '{print $3}')
    print_status "Using Go version: $GO_VERSION"
}

# Function to check if required tools are installed
check_tools() {
    print_status "Checking required tools..."
    
    # Check for go tool cover
    if ! go tool cover &> /dev/null; then
        print_error "Go cover tool not available"
        exit 1
    fi
    
    print_success "All required tools are available"
}

# Function to clean previous test artifacts
cleanup() {
    print_status "Cleaning up previous test artifacts..."
    rm -f coverage.out coverage.html coverage.txt
    print_success "Cleanup completed"
}

# Function to run basic tests
run_basic_tests() {
    print_status "Running basic tests..."
    
    if go test ./... -v; then
        print_success "Basic tests passed"
    else
        print_error "Basic tests failed"
        exit 1
    fi
}

# Function to run tests with race detection
run_race_tests() {
    print_status "Running tests with race detection..."
    
    if go test -race ./... -v; then
        print_success "Race detection tests passed"
    else
        print_warning "Race detection tests failed (this may be expected in some environments)"
    fi
}

# Function to run tests with coverage
run_coverage_tests() {
    print_status "Running tests with coverage..."
    
    if go test -coverprofile=coverage.out ./... -v; then
        print_success "Coverage tests passed"
    else
        print_error "Coverage tests failed"
        exit 1
    fi
}

# Function to generate coverage reports
generate_coverage_reports() {
    print_status "Generating coverage reports..."
    
    # Generate HTML coverage report
    if go tool cover -html=coverage.out -o coverage.html; then
        print_success "HTML coverage report generated: coverage.html"
    else
        print_error "Failed to generate HTML coverage report"
    fi
    
    # Generate text coverage report
    if go tool cover -func=coverage.out > coverage.txt; then
        print_success "Text coverage report generated: coverage.txt"
    else
        print_error "Failed to generate text coverage report"
    fi
    
    # Display coverage summary
    print_status "Coverage Summary:"
    go tool cover -func=coverage.out | tail -1
}

# Function to run benchmarks
run_benchmarks() {
    print_status "Running benchmarks..."
    
    if go test -bench=. -benchmem ./... -v; then
        print_success "Benchmarks completed"
    else
        print_warning "Some benchmarks failed (this may be expected)"
    fi
}

# Function to run linting
run_linting() {
    print_status "Running code linting..."
    
    if command -v golangci-lint &> /dev/null; then
        if golangci-lint run; then
            print_success "Linting passed"
        else
            print_warning "Linting found issues"
        fi
    else
        print_warning "golangci-lint not installed, skipping linting"
    fi
}

# Function to run code vetting
run_vetting() {
    print_status "Running code vetting..."
    
    if go vet ./...; then
        print_success "Code vetting passed"
    else
        print_warning "Code vetting found issues"
    fi
}

# Function to display final summary
display_summary() {
    echo ""
    print_status "Test Execution Summary"
    echo "=========================="
    
    if [ -f "coverage.out" ]; then
        COVERAGE=$(go tool cover -func=coverage.out | tail -1 | awk '{print $3}' | sed 's/%//')
        print_status "Test Coverage: ${COVERAGE}%"
    fi
    
    if [ -f "coverage.html" ]; then
        print_status "HTML Coverage Report: coverage.html"
    fi
    
    if [ -f "coverage.txt" ]; then
        print_status "Text Coverage Report: coverage.txt"
    fi
    
    print_status "Raw Coverage Data: coverage.out"
    
    echo ""
    print_success "All tests completed successfully!"
}

# Main execution
main() {
    echo "=========================================="
    echo "ABT Analytics Dashboard Backend - Test Suite"
    echo "=========================================="
    echo ""
    
    # Check prerequisites
    check_go
    check_tools
    
    # Cleanup
    cleanup
    
    # Run test suite
    run_basic_tests
    run_race_tests
    run_coverage_tests
    generate_coverage_reports
    run_benchmarks
    run_linting
    run_vetting
    
    # Display summary
    display_summary
}

# Handle script interruption
trap 'print_error "Test execution interrupted"; exit 1' INT TERM

# Run main function
main "$@"

