package campaign

import (
	"campaign-optimization-engine/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetCampaigns returns all campaigns
func GetCampaigns(c echo.Context) error {
	// Get campaigns from service (which fetches from DB)
	campaigns, err := GetAllCampaigns()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, campaigns)
}

// CreateCampaign adds a new campaign
func CreateCampaign(c echo.Context) error {
	var newCampaign models.Campaign
	if err := c.Bind(&newCampaign); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Add the campaign using the service (which stores it in DB)
	err := AddCampaign(&newCampaign)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, newCampaign)
}
