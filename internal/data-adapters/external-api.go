package data_adapters

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// APIResponse represents the structure of CPC & CVR data from an API
type APIResponse struct {
	Platform string  `json:"platform"`
	CPC      float64 `json:"cpc"`
	CVR      float64 `json:"cvr"`
}

// FetchFromAPI fetches real-time CPC & CVR from an external API
func FetchFromAPI(platform string) (float64, float64, error) {
	url := "https://external-ad-platform.com/api/cpc_cvr?platform=" + platform
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch API data for %s: %v", platform, err)
		return GenerateRandomCPC(), GenerateRandomCVR(), nil // Fallback
	}
	defer resp.Body.Close()

	var data APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Failed to parse API response for %s: %v", platform, err)
		return GenerateRandomCPC(), GenerateRandomCVR(), nil // Fallback
	}

	return data.CPC, data.CVR, nil
}

// GenerateRandomCPC generates a random CPC as a fallback
func GenerateRandomCPC() float64 {
	return 1.5 + rand.Float64()*0.5
}

// GenerateRandomCVR generates a random CVR as a fallback
func GenerateRandomCVR() float64 {
	return 3.0 + rand.Float64()*1.0
}

// UpdateCPC_CVR fetches and updates Redis with the latest CPC & CVR values
func UpdateCPC_CVR() {
	platforms := []string{"google", "facebook"} // Extend as needed

	for _, platform := range platforms {
		cpc, cvr, err := FetchFromAPI(platform)
		if err != nil {
			log.Printf("Using fallback CPC & CVR for %s", platform)
		}

		UpdateCPC(platform, cpc)
		UpdateCVR(platform, cvr)

		log.Printf("Updated %s - CPC: $%.2f, CVR: %.2f%%", platform, cpc, cvr)
	}
}

// StartAPIUpdater runs a background job to fetch and update CPC & CVR every 10 seconds
func CPCCVRUpdater() {
	go func() {
		for {
			UpdateCPC_CVR()
			time.Sleep(10 * time.Second)
		}
	}()
}
