package message_queue

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

var rdb *redis.Client

// Initialize Redis client

// ConsumeMessages consumes messages from Kafka topic with retries and exponential backoff
func ConsumeMessages() error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	// Set Kafka consumer group
	consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "bidding_group", config)
	if err != nil {
		return fmt.Errorf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Consume messages from the Kafka topic
	for {
		// Retry consuming messages with exponential backoff
		err := consumerGroup.Consume(nil, []string{"bidding_decisions"}, &ConsumerHandler{})
		if err != nil {
			log.Printf("Error consuming messages: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second) // Retry with a delay
			continue                    // Retry consuming
		}
	}
}

type ConsumerHandler struct{}

// Setup initializes resources needed for consumption
func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Initialize any necessary resources here
	return nil
}

// Cleanup cleans up resources after consumption
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	// Clean up resources here
	return nil
}

// ConsumeClaim processes messages from a Kafka topic partition with retry logic and tracking
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		messageID := fmt.Sprintf("%d-%s", message.Partition, message.Offset)
		log.Printf("Consumed message: %s", string(message.Value))

		// Check if this message has already been processed or retried
		retryCount, err := rdb.Get(context.Background(), messageID).Int()
		if err == nil && retryCount >= 3 { // Stop retrying if failed 3 times
			log.Printf("Message %s has failed 3 times, skipping.", messageID)
			session.MarkMessage(message, "") // Mark message as processed (failed)
			continue
		}

		// Process the message (e.g., business logic, bidding engine)
		if err := processMessage(message); err != nil {
			// If processing fails, increment retry count
			rdb.Incr(context.Background(), messageID)
			session.MarkMessage(message, "") // Mark message as processed
			log.Printf("Message %s failed processing. Retry attempt %d", messageID, retryCount+1)
			return err // Exit so the message will be retried
		}

		// Mark the message as processed successfully
		rdb.Del(context.Background(), messageID) // Reset retry count
		session.MarkMessage(message, "")         // Mark message as processed
	}
	return nil
}

// processMessage simulates message processing logic
func processMessage(message *sarama.ConsumerMessage) error {
	// Implement the actual business logic here (e.g., bidding decision)
	log.Printf("Processing message: %s", string(message.Value))
	return nil // Simulate success
}
