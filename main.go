// main.go

package main

import (
	"campaign-optimization-engine/internal/bidding"
	data_adapters "campaign-optimization-engine/internal/data-adapters"
	"campaign-optimization-engine/server"
	"log"
	"sync"
)

func main() {

	// Initialize Redis (if using it)
	data_adapters.InitRedis()

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
