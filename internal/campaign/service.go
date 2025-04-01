package campaign

import (
	"campaign-optimization-engine/internal/data-adapters"
	"campaign-optimization-engine/internal/models"
)

// GetAllCampaigns returns all active campaigns from the database
func GetAllCampaigns() ([]models.Campaign, error) {
	return data_adapters.GetAllCampaigns()
}

// AddCampaign adds a new campaign to the database
func AddCampaign(c *models.Campaign) error {
	return data_adapters.CreateCampaign(c)
}
