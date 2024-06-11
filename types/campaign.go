package types

import "time"

type Campaign struct {
	ID            int         `json:"id"`
	Name          string      `json:"name"`
	Recruitment   bool        `json:"recruitment"`
	Judge         string      `json:"judge"`
	Timekeeping   string      `json:"timekeeping"`
	Cadence       string      `json:"cadence"`
	CreatedAt     time.Time   `json:"created_at,omitempty"`
	UpdatedAt     time.Time   `json:"updated_at,omitempty"`
	LastAdventure time.Time   `json:"last_adventure,omitempty"`
	Characters    []Character `json:"characters"`
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

func NewCampaign(id int) *Campaign {
	characters := make([]Character, 0)
	return &Campaign{id, "", false, "", "", "", time.Now(), time.Now(), time.Now(), characters}
}

func NewCampaignFull(id int, name, judge, timekeeping, cadence string, recruit bool, cdate, udate, ldate time.Time, characters []Character) *Campaign {
	return &Campaign{id, name, recruit, judge, timekeeping, cadence, time.Now(), time.Now(), time.Now(), characters}
}
