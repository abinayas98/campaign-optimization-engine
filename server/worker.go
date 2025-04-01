package server

import (
	"campaign-optimization-engine/internal/message_queue"
	"fmt"
	"log"
)

func StartWorker() error {
	log.Println("Starting Kafka consumer worker...")

	// Start consuming messages from Kafka
	if err := message_queue.ConsumeMessages(); err != nil {
		return fmt.Errorf("Error consuming messages: %v", err)
	}

	return nil
}
