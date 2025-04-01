package priorityqueue

import (
	"campaign-optimization-engine/internal/models"
	"container/heap"
	"log"
)

// PriorityQueue represents a heap-based priority queue
type PriorityQueue []*models.BidResult

// Implement the heap.Interface
func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	// Prioritize campaigns with higher bid amount
	return pq[i].BidAmount > pq[j].BidAmount
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

// DequeueCampaign removes and returns the highest-priority campaign (highest bid)
func DequeueCampaign() *models.BidResult {
	if len(pq) == 0 {
		log.Println("⚠️ Priority queue is empty")
		return nil
	}
	// Pop the highest priority bid result
	bestBid := heap.Pop(&pq).(*models.BidResult)
	return bestBid
}

// GetQueueLength returns the current length of the queue
func GetQueueLength() int {
	return len(pq)
}
