package data_adapters

import (
	"campaign-optimization-engine/internal/models"
)

// CreateCampaign adds a new campaign to the database
func CreateCampaign(campaign *models.Campaign) error {
	if err := DB.Create(&campaign).Error; err != nil {
		return err
	}
	return nil
}

// GetAllCampaigns retrieves all campaigns from the database
func GetAllCampaigns() ([]models.Campaign, error) {
	var campaigns []models.Campaign
	if err := DB.Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}
