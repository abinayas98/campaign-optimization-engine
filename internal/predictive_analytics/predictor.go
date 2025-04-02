package predictive_analytics

import (
	"github.com/sajari/regression"
	"log"
	"math/rand"
	"time"
)

// Historical data for platforms (extend as necessary)
var HistoricalCPCData = map[string][]float64{
	"google":   {1.5, 1.6, 1.4, 1.8, 1.7}, // Sample historical CPC data
	"facebook": {1.2, 1.3, 1.1, 1.4, 1.5}, // Sample historical CPC data
}

var HistoricalCVRData = map[string][]float64{
	"google":   {3.0, 3.2, 2.9, 3.4, 3.1}, // Sample historical CVR data
	"facebook": {4.0, 4.1, 3.8, 4.3, 4.2}, // Sample historical CVR data
}

// PredictCPC predicts the next CPC value using linear regression with historical and real-time data
func PredictCPC(historicalData []float64, realTimeCPC float64) float64 {
	if len(historicalData) < 2 {
		// Fallback: If not enough data, return random CPC
		return 1.0 + rand.Float64()*2.0
	}

	var r regression.Regression
	r.SetObserved("CPC")
	r.SetVar(0, "Time")

	// Train with historical data
	for i, value := range historicalData {
		r.Train(regression.DataPoint(value, []float64{float64(i)}))
	}
	r.Run() // Train the model

	// Use real-time data as a feature to adjust the prediction
	// Adding real-time data as a feature to shift the trend
	predicted, err := r.Predict([]float64{float64(len(historicalData))})
	if err != nil {
		log.Println("⚠️ Prediction failed, using fallback:", err)
		return 1.0 + rand.Float64()*2.0
	}

	// Adjust the predicted CPC based on real-time data (simple blending)
	return (predicted + realTimeCPC) / 2.0
}

// PredictCVR predicts the next CVR value using linear regression with historical and real-time data
func PredictCVR(historicalData []float64, realTimeCVR float64) float64 {
	if len(historicalData) < 2 {
		// Fallback: If not enough data, return random CVR
		return 2.0 + rand.Float64()*4.0
	}

	var r regression.Regression
	r.SetObserved("CVR")
	r.SetVar(0, "Time")

	// Train with historical data
	for i, value := range historicalData {
		r.Train(regression.DataPoint(value, []float64{float64(i)}))
	}
	r.Run() // Train the model

	// Use real-time data as a feature to adjust the prediction
	// Adding real-time data as a feature to shift the trend
	predicted, err := r.Predict([]float64{float64(len(historicalData))})
	if err != nil {
		log.Println("⚠️ Prediction failed, using fallback:", err)
		return 2.0 + rand.Float64()*4.0
	}

	// Adjust the predicted CVR based on real-time data (simple blending)
	return (predicted + realTimeCVR) / 2.0
}

// GenerateRandomCPC generates a random CPC value as a fallback
func GenerateRandomCPC() float64 {
	// Random CPC between $1.5 and $2.0
	return 1.5 + rand.Float64()*0.5
}

// GenerateRandomCVR generates a random CVR value as a fallback
func GenerateRandomCVR() float64 {
	// Random CVR between 3.0% and 4.0%
	return 3.0 + rand.Float64()*1.0
}

// Initialize random seed (should be done once at the start of the application)
func init() {
	rand.Seed(time.Now().UnixNano()) // Ensure that the random number generator is seeded with current time
}
