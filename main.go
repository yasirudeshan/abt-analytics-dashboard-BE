package main

import (
	"abt-analytics-dashboard/internal/api"
	"abt-analytics-dashboard/internal/config"
	"abt-analytics-dashboard/internal/processor"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Error loading .env file: %v, using system environment variables", err)
	} else {
		log.Println("Successfully loaded .env file")
	}

	// Load configuration
	cfg := config.Load()

	// Initialize data processor
	dataProcessor := processor.New()

	// Process the dataset file if provided
	if cfg.DataFilePath != "" {
		log.Printf("Processing dataset from: %s", cfg.DataFilePath)
		start := time.Now()

		if err := dataProcessor.ProcessDataset(cfg.DataFilePath); err != nil {
			log.Fatalf("Failed to process dataset: %v", err)
		}

		duration := time.Since(start)
		log.Printf("Dataset processed successfully in %v", duration)
	} else {
		log.Println("No dataset file provided. Using sample data for development.")
		dataProcessor.LoadSampleData()
	}

	// Initialize API server
	server := api.NewServer(dataProcessor, cfg)

	// Setup graceful shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Trigger graceful shutdown
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Printf("Starting server on port %s", cfg.Port)
	log.Printf("Server running at http://localhost%s", cfg.Port)

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	fmt.Println("Server stopped gracefully")
}
