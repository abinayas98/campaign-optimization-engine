package bidding

import (
	"campaign-optimization-engine/internal/data-adapters"
	"campaign-optimization-engine/internal/models"
	"campaign-optimization-engine/internal/predictive_analytics"
)

// OptimizeBid determines the best bid for a campaign using both real-time and historical data
func OptimizeBid(c models.Campaign) models.BidResult {
	bestPlatform := ""
	bestBidAmount := 0.0
	bestParetoScore := 0.0

	for _, platform := range c.Platforms {
		// Fetch real-time & historical data
		realTimeCPC, realTimeCVR, err := data_adapters.GetCachedPrediction(platform)
		if err != nil {
			realTimeCPC = predictive_analytics.GenerateRandomCPC()
			realTimeCVR = predictive_analytics.GenerateRandomCVR()
		}

		historicalCPC := predictive_analytics.HistoricalCPCData[platform]
		historicalCVR := predictive_analytics.HistoricalCVRData[platform]

		// Predict CPC & CVR using AI model
		predictedCPC := predictive_analytics.PredictCPC(historicalCPC, realTimeCPC)
		predictedCVR := predictive_analytics.PredictCVR(historicalCVR, realTimeCVR)

		// Calculate bid amount
		bidAmount := predictedCPC * (1 + (predictedCVR / 10))

		// Compute Pareto Efficiency Score (maximize CVR, minimize CPC)
		paretoScore := (predictedCVR / predictedCPC) * c.Budget

		// Select the platform with the highest Pareto score
		if paretoScore > bestParetoScore {
			bestParetoScore = paretoScore
			bestPlatform = platform
			bestBidAmount = bidAmount
		}
	}

	return models.BidResult{
		CampaignID: c.ID,
		Platform:   bestPlatform,
		BidAmount:  bestBidAmount,
	}
}
