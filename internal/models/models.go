package models

import "time"

// Campaign struct represents an ad campaign.
type Campaign struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Budget      float64   `json:"budget"`
	TargetReach int       `json:"target_reach"`
	Platforms   []string  `json:"platforms"` // ["google", "facebook"]
	CPC         float64   `json:"cpc"`       // Cost per Click
	CVR         float64   `json:"cvr"`       // Conversion Rate
	CreatedAt   time.Time `json:"created_at"`
}

// BidResult stores bid decision
type BidResult struct {
	CampaignID string
	Platform   string
	BidAmount  float64
}
