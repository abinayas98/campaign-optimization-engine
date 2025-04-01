package campaign

import (
	"campaign-optimization-engine/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCampaigns returns all campaigns
func GetCampaigns(c echo.Context) error {
	// Call the service to get campaigns
	campaigns := GetAllCampaigns() // Get campaigns from service
	return c.JSON(http.StatusOK, campaigns)
}

// CreateCampaign adds a new campaign
func CreateCampaign(c echo.Context) error {
	var newCampaign Campaign
	if err := c.Bind(&newCampaign); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Call the service to add the new campaign
	AddCampaign(models.Campaign(newCampaign)) // Add campaign to service

	return c.JSON(http.StatusCreated, newCampaign)
}
