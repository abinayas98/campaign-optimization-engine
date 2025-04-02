package priorityqueue

import (
	c "campaign-optimization-engine/internal/campaign"
	data_adapters "campaign-optimization-engine/internal/data-adapters"
	"campaign-optimization-engine/internal/models"
	"container/heap"
	"log"
	"sync"
)

// PriorityQueue represents a heap-based priority queue
type PriorityQueue struct {
	queue PriorityQueueHeap
	mu    sync.Mutex
}

// PriorityQueueHeap implements heap.Interface
type PriorityQueueHeap []*models.BidResult

func (pq PriorityQueueHeap) Len() int { return len(pq) }
func (pq PriorityQueueHeap) Less(i, j int) bool {
	return pq[i].BidAmount > pq[j].BidAmount // Higher bid has higher priority
}
func (pq PriorityQueueHeap) Swap(i, j int) { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueueHeap) Push(x interface{}) {
	*pq = append(*pq, x.(*models.BidResult))
}
func (pq *PriorityQueueHeap) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

var pq PriorityQueue

// InitQueue initializes the priority queue
func InitQueue() {
	pq.queue = make(PriorityQueueHeap, 0)
	heap.Init(&pq.queue)
}

// EnqueueCampaign adds a bid result to the priority queue
func EnqueueCampaign(bid *models.BidResult) {
	pq.mu.Lock()
	defer pq.mu.Unlock()
	heap.Push(&pq.queue, bid)
	log.Printf("✅ Enqueued Campaign: %s, Bid Amount: $%.2f", bid.CampaignID, bid.BidAmount)
}

func DequeueCampaign() *models.BidResult {
	pq.mu.Lock()
	defer pq.mu.Unlock() // Ensure unlock even if there's an error

	if pq.queue.Len() == 0 {
		return nil
	}
	return heap.Pop(&pq.queue).(*models.BidResult)
}

// ProcessBids continuously dequeues and processes the highest-priority bid
func ProcessBids() {
	for {
		pq.mu.Lock()
		if pq.queue.Len() == 0 {
			pq.mu.Unlock()
			continue
		}
		bestBid := DequeueCampaign()
		pq.mu.Unlock()

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
			err = data_adapters.UpdateCampaignBudget(campaign.ID, campaign.Budget)
			if err != nil {
				log.Printf("⚠️ Failed to update budget for Campaign %s", campaign.ID)
				continue
			}

			log.Printf("✅ Processed Campaign: %s, Remaining Budget: $%.2f", campaign.ID, campaign.Budget)
		} else {
			log.Printf("⚠️ Skipping Campaign %s due to insufficient budget", bestBid.CampaignID)
		}
	}
}
