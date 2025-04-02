// main.go

package main

import (
	"log"
	"sync"

	"campaign-optimization-engine/internal/bidding"
	"campaign-optimization-engine/internal/data-adapters"
	"campaign-optimization-engine/server"
)

func main() {
	// Initialize Redis (if using it)
	data_adapters.InitRedis()

	// Initialize the PostgreSQL database connection
	data_adapters.InitDB()
	defer data_adapters.CloseDB()

	// Start the Bid Optimization Engine
	bidding.StartBidEngine()

	// Use WaitGroup to wait for both server and worker to run concurrently
	var wg sync.WaitGroup

	// Start the API server in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.StartAPI(); err != nil {
			log.Fatalf("Error starting API server: %v", err)
		}
	}()

	// Start the worker to consume Kafka messages in another goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.StartWorker(); err != nil {
			log.Fatalf("Error starting worker: %v", err)
		}
	}()

	// Wait for both server and worker to finish
	wg.Wait()
}
