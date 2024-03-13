package types

import "time"

type Campaign struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Recruitment   bool      `json:"recruitment"`
	Judge         string    `json:"judge"`
	Timekeeping   string    `json:"timekeeping"`
	Cadence       string    `json:"cadence"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	LastAdventure time.Time `json:"last_adventure,omitempty"`
}

type CreateCampaignRequest struct {
	Name          string    `json:"name"`
	Recruitment   bool      `json:"recruitment"`
	Judge         string    `json:"judge"`
	Timekeeping   string    `json:"timekeeping"`
	Cadence       string    `json:"cadence"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	LastAdventure time.Time `json:"last_adventure,omitempty"`
}
type UpdateCampaignRequest struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Recruitment   bool      `json:"recruitment"`
	Judge         string    `json:"judge"`
	Timekeeping   string    `json:"timekeeping"`
	Cadence       string    `json:"cadence"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	LastAdventure time.Time `json:"last_adventure,omitempty"`
}

func ValidateCampaign(u *Campaign) bool { return true }
