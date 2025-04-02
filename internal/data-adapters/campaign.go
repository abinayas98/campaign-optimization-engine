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

// GetCampaignByID retrieves a campaign by its ID
func GetCampaignByID(campaignID string) (*models.Campaign, error) {
	var campaign models.Campaign
	if err := DB.Where("id = ?", campaignID).First(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

// UpdateCampaignBudget updates the budget of a specific campaign
func UpdateCampaignBudget(campaignID string, newBudget float64) error {
	return DB.Model(&models.Campaign{}).Where("id = ?", campaignID).Update("budget", newBudget).Error
}
