package bidding

import (
	"campaign-optimization-engine/internal/campaign"
	"campaign-optimization-engine/internal/message_queue"
	"campaign-optimization-engine/internal/models"
	"campaign-optimization-engine/internal/priorityqueue"
	"encoding/json"
	"log"
	"sync"
	"time"
)

// Worker function to process campaigns and generate bids
func bidWorker(campaignChan <-chan models.Campaign, bidResultChan chan<- *models.BidResult, wg *sync.WaitGroup) {
	defer wg.Done()

	for c := range campaignChan {
		bid := OptimizeBid(c)
		bidResultChan <- &bid
	}
}

// StartBidEngine runs bid optimization concurrently
func StartBidEngine() {
	priorityqueue.InitQueue()
	go priorityqueue.ProcessBids()
	
	campaignChan := make(chan models.Campaign, 10)
	bidResultChan := make(chan *models.BidResult, 10)
	var wg sync.WaitGroup

	// Start worker goroutines
	numWorkers := 4
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go bidWorker(campaignChan, bidResultChan, &wg)
	}

	// Periodically fetch campaigns
	go func() {
		for {
			log.Println("Running Bid Optimization Engine...")

			campaigns, _ := campaign.GetAllCampaigns()
			for _, c := range campaigns {
				campaignChan <- c
			}

			time.Sleep(10 * time.Second)
		}
	}()

	// Process bid results and send to Kafka
	go func() {
		for bid := range bidResultChan {
			// Convert bid result to JSON
			bidJSON, err := json.Marshal(bid)
			if err != nil {
				log.Printf("Error marshalling bid result: %v", err)
				continue
			}

			// Send bid to Kafka
			err = message_queue.ProduceMessage(string(bidJSON))
			if err != nil {
				log.Printf("Failed to send bid to Kafka: %v", err)
			} else {
				log.Printf("Bid decision sent to Kafka: %s", string(bidJSON))
			}

			// Enqueue the bid in the priority queue
			priorityqueue.EnqueueCampaign(bid)
		}
	}()
}
