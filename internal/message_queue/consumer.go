package message_queue

import (
	"github.com/Shopify/sarama"
	"log"
)

func ConsumeMessages() error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "bidding_group", config)
	if err != nil {
		log.Fatalf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	// Consume messages from the Kafka topic
	for {
		if err := consumerGroup.Consume(nil, []string{"bidding_decisions"}, &ConsumerHandler{}); err != nil {
			log.Fatalf("Error consuming messages: %v", err)
			return err
		}
	}
}

type ConsumerHandler struct{}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a Kafka topic partition
func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Consumed message: %s", string(message.Value))
		session.MarkMessage(message, "") // Mark the message as processed
	}
	return nil
}
func (h *ConsumerHandler) ConsumeMessage(msg *sarama.ConsumerMessage) error {
	log.Printf("Received message: %s", msg.Value)
	// Process the bidding task
	return nil
}
