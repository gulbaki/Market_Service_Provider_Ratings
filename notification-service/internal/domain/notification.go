package domain

import "time"

// Notification represents a new rating-based notification
// for a specific provider.
// swagger:model
type Notification struct {
	// ID is an internal identifier
	ID int `json:"id"`
	// ProviderID links to the specific service provider
	ProviderID int `json:"providerId"`
	// Score is the rating score from the Rating Service
	Score int `json:"score"`
	// Comment is an optional user comment
	Comment string `json:"comment"`
	// CreatedAt is the timestamp when rating was created
	CreatedAt time.Time `json:"createdAt"`
}
