package message_queue

import (
	"fmt"
	"github.com/Shopify/sarama"
	"log"
)

// Kafka producer settings
const (
	kafkaBroker = "localhost:9092"
	topicName   = "bidding_decisions"
)

func ProduceMessage(message string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Initialize a Kafka producer
	producer, err := sarama.NewSyncProducer([]string{kafkaBroker}, config)
	if err != nil {
		return fmt.Errorf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Create a message
	msg := &sarama.ProducerMessage{
		Topic: topicName,
		Value: sarama.StringEncoder(message),
	}

	// Send the message
	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent: %s", message)
	return nil
}
