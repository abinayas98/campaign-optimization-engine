package bidding

import (
	"campaign-optimization-engine/internal/data-adapters"
	"campaign-optimization-engine/internal/models"
	"campaign-optimization-engine/internal/predictive_analytics"
	"log"
)

// OptimizeBid determines the best bid for a campaign using both real-time and historical data
func OptimizeBid(c models.Campaign) models.BidResult {
	bestPlatform := ""
	bestBidAmount := 0.0

	// Iterate through platforms and find the best bid
	for _, platform := range c.Platforms {
		// Fetch real-time data from Redis
		realTimeCPC, realTimeCVR, err := data_adapters.GetCachedPrediction(platform)
		if err != nil {
			log.Printf("⚠️ Using fallback prediction for %s: %v", platform, err)
			// Fallback to historical data if real-time data is not available
			realTimeCPC = predictive_analytics.GenerateRandomCPC()
			realTimeCVR = predictive_analytics.GenerateRandomCVR()
		}

		// Get historical data for the platform
		historicalCPC := predictive_analytics.HistoricalCPCData[platform]
		historicalCVR := predictive_analytics.HistoricalCVRData[platform]

		// Predict CPC & CVR using historical data and real-time data
		predictedCPC := predictive_analytics.PredictCPC(historicalCPC, realTimeCPC)
		predictedCVR := predictive_analytics.PredictCVR(historicalCVR, realTimeCVR)

		// Calculate the bid amount
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
