package priorityqueue

import (
	c "campaign-optimization-engine/internal/campaign"
	"campaign-optimization-engine/internal/models"
	"container/heap"
	"log"
)

// PriorityQueue represents a heap-based priority queue
type PriorityQueue []*models.BidResult

// Implement the heap.Interface
func (pq PriorityQueue) Len() int { return len(pq) }

// Modify Less function to use urgency-based sorting
func (pq PriorityQueue) Less(i, j int) bool {
	campaignI, errI := c.GetCampaignByID(pq[i].CampaignID)
	campaignJ, errJ := c.GetCampaignByID(pq[j].CampaignID)

	if errI != nil || errJ != nil {
		// If campaign not found, fallback to simple bid comparison
		return pq[i].BidAmount > pq[j].BidAmount
	}

	// Define urgency score
	urgencyI := campaignI.Budget / float64(campaignI.TargetReach)
	urgencyJ := campaignJ.Budget / float64(campaignJ.TargetReach)

	// Compute priority score
	scoreI := pq[i].BidAmount * urgencyI
	scoreJ := pq[j].BidAmount * urgencyJ

	return scoreI > scoreJ
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(*models.BidResult))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

// Global priority queue
var pq PriorityQueue

// InitQueue initializes the priority queue
func InitQueue() {
	pq = make(PriorityQueue, 0)
	heap.Init(&pq)
}

// EnqueueCampaign adds a bid result to the priority queue
func EnqueueCampaign(bid *models.BidResult) {
	heap.Push(&pq, bid)
	log.Printf("Enqueued Campaign: %s, Bid Amount: $%.2f", bid.CampaignID, bid.BidAmount)
}

// DequeueCampaign removes and returns the highest-priority campaign (highest urgency-adjusted bid)
func DequeueCampaign() *models.BidResult {
	for len(pq) > 0 {
		bestBid := heap.Pop(&pq).(*models.BidResult)

		// Fetch the campaign details
		campaign, err := c.GetCampaignByID(bestBid.CampaignID)
		if err != nil {
			log.Printf("⚠️ Skipping Campaign %s: Not Found", bestBid.CampaignID)
			continue
		}

		// Ensure the campaign has enough budget left
		if campaign.Budget >= bestBid.BidAmount {
			// Deduct the bid amount from the campaign budget
			campaign.Budget -= bestBid.BidAmount

			// Update the campaign budget in the database
			err = c.UpdateCampaignBudget(campaign.ID, campaign.Budget)
			if err != nil {
				log.Printf("⚠️ Failed to update budget for Campaign %s", campaign.ID)
				continue
			}

			log.Printf("✅ Processed Campaign: %s, Remaining Budget: $%.2f", campaign.ID, campaign.Budget)
			return bestBid
		}

		log.Printf("⚠️ Skipping Campaign %s due to insufficient budget", bestBid.CampaignID)
	}
	return nil
}

// GetQueueLength returns the current length of the queue
func GetQueueLength() int {
	return len(pq)
}
