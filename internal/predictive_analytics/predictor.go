package predictive_analytics

import (
	"log"
	"math/rand"

	"github.com/sajari/regression"
)

// PredictCPC predicts the next CPC value using linear regression
func PredictCPC(historicalData []float64) float64 {
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

	// Predict next CPC
	predicted, err := r.Predict([]float64{float64(len(historicalData))})
	if err != nil {
		log.Println("⚠️ Prediction failed, using fallback:", err)
		return 1.0 + rand.Float64()*2.0
	}

	return predicted
}

// PredictCVR predicts the next CVR value using linear regression
func PredictCVR(historicalData []float64) float64 {
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

	// Predict next CVR
	predicted, err := r.Predict([]float64{float64(len(historicalData))})
	if err != nil {
		log.Println("⚠️ Prediction failed, using fallback:", err)
		return 2.0 + rand.Float64()*4.0
	}

	return predicted
}
