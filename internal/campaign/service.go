package campaign

import (
	"sync"

	"campaign-optimization-engine/internal/models"
)

var (
	campaigns []models.Campaign
	mu        sync.Mutex
)

// GetAllCampaigns returns all active campaigns
func GetAllCampaigns() []models.Campaign {
	mu.Lock()
	defer mu.Unlock()
	return campaigns
}

// AddCampaign adds a new campaign to the list
func AddCampaign(c models.Campaign) {
	mu.Lock()
	defer mu.Unlock()
	campaigns = append(campaigns, c)
}
