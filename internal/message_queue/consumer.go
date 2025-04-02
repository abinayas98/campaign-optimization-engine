package message_queue

import (
	"campaign-optimization-engine/internal/models"
	"campaign-optimization-engine/internal/priorityqueue"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client

// ConsumeMessages consumes messages from Kafka topic with retries and exponential backoff
func ConsumeMessages() error {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin

	consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "bidding_group", config)
	if err != nil {
		return fmt.Errorf("Error creating consumer group: %v", err)
	}
	defer consumerGroup.Close()

	for {
		err := consumerGroup.Consume(context.Background(), []string{"bidding_decisions"}, &ConsumerHandler{})
		if err != nil {
			log.Printf("Error consuming messages: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

type ConsumerHandler struct{}

func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Consumed message: %s", string(message.Value))

		// Parse the Kafka message (BidResult)
		var bid models.BidResult
		err := json.Unmarshal(message.Value, &bid)
		if err != nil {
			log.Printf("Error parsing bid result: %v", err)
			continue
		}

		// Enqueue the bid into the priority queue
		priorityqueue.EnqueueCampaign(&bid)
		log.Printf("Bid enqueued: Campaign %s, Bid Amount: %.2f", bid.CampaignID, bid.BidAmount)

		session.MarkMessage(message, "")
	}
	return nil
}
