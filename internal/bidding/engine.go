package bidding

import (
	"campaign-optimization-engine/internal/campaign"
	"campaign-optimization-engine/internal/priorityqueue"
	"log"
	"time"
)

// StartBidEngine runs bid optimization periodically
func StartBidEngine() {
	// Initialize the priority queue once, before starting the optimization loop
	priorityqueue.InitQueue()

	go func() {
		for {
			log.Println("🚀 Running Bid Optimization Engine...")

			// Get all active campaigns
			campaigns := campaign.GetAllCampaigns()

			// Process each campaign, add it to the priority queue with its bid
			for _, c := range campaigns {
				// Optimize bid for the current campaign
				bid := OptimizeBid(c)

				// Enqueue the campaign to the priority queue
				priorityqueue.EnqueueCampaign(&bid)
			}

			// Process the highest-priority (largest bid) campaign
			for pqLen := priorityqueue.GetQueueLength(); pqLen > 0; pqLen = priorityqueue.GetQueueLength() {
				bestBid := priorityqueue.DequeueCampaign()
				if bestBid != nil {
					log.Printf("✅ Placing bid for Campaign: %s, Platform: %s, Amount: $%.2f\n", bestBid.CampaignID, bestBid.Platform, bestBid.BidAmount)
				}
			}

			time.Sleep(10 * time.Second) // Run every 10 seconds
		}
	}()
}
