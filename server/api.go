package server

import (
	"campaign-optimization-engine/internal/campaign" // Adjust based on your module path
	"github.com/labstack/echo/v4"
	"log"
)

func StartAPI() error {
	e := echo.New()

	// Define routes
	e.GET("/api/campaigns", campaign.GetCampaigns)
	e.POST("/api/campaign", campaign.CreateCampaign)

	// Start server
	err := e.Start(":8080")
	if err != nil {
		log.Fatalf("Error starting API server: %v", err)
	}
	return err
}
