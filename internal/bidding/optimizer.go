package bidding

import (
	"campaign-optimization-engine/internal/models"
	"campaign-optimization-engine/internal/predictive_analytics"
	"log"
)

// Mock historical CPC and CVR data
var historicalCPC = map[string][]float64{
	"google":   {1.5, 1.6, 1.4, 1.8, 1.7},
	"facebook": {1.2, 1.3, 1.1, 1.4, 1.5},
}

var historicalCVR = map[string][]float64{
	"google":   {3.0, 3.2, 2.9, 3.4, 3.1},
	"facebook": {4.0, 4.1, 3.8, 4.3, 4.2},
}

// OptimizeBid determines the best bid for a campaign
func OptimizeBid(c models.Campaign) models.BidResult {
	bestPlatform := ""
	bestBidAmount := 0.0

	// Iterate through platforms and find the best bid
	for _, platform := range c.Platforms {
		// Get historical data
		cpcHistory := historicalCPC[platform]
		cvrHistory := historicalCVR[platform]

		// Predict CPC & CVR using ML
		predictedCPC := predictive_analytics.PredictCPC(cpcHistory)
		predictedCVR := predictive_analytics.PredictCVR(cvrHistory)

		// Calculate bid amount
		bidAmount := predictedCPC * (1 + (predictedCVR / 10))

		log.Printf("📊 Campaign: %s, Platform: %s, Predicted CPC: $%.2f, Predicted CVR: %.2f%%, Bid: $%.2f\n",
			c.ID, platform, predictedCPC, predictedCVR, bidAmount)

		// Choose the best platform with the highest bid
		if bidAmount > bestBidAmount {
			bestPlatform = platform
			bestBidAmount = bidAmount
		}
	}

	log.Printf("✅ Best Bid - Campaign: %s, Platform: %s, Amount: $%.2f\n", c.ID, bestPlatform, bestBidAmount)

	return models.BidResult{
		CampaignID: c.ID,
		Platform:   bestPlatform,
		BidAmount:  bestBidAmount,
	}
}
