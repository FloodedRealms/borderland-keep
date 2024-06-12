package types

import (
	"time"
)

type CampaignRecord struct {
	ID            int               `json:"id"`
	ClientId      string            `json:"client_id"`
	Name          string            `json:"name"`
	Recruitment   bool              `json:"recruitment"`
	Judge         string            `json:"judge"`
	Timekeeping   string            `json:"timekeeping"`
	Cadence       string            `json:"cadence"`
	CreatedAt     time.Time         `json:"created_at,omitempty"`
	UpdatedAt     time.Time         `json:"updated_at,omitempty"`
	LastAdventure time.Time         `json:"last_adventure,omitempty"`
	Characters    []CharacterRecord `json:"characters"`
}

type CreateCampaignRecordRequest struct {
	Name        string `json:"name"`
	Recruitment bool   `json:"recruitment"`
	Judge       string `json:"judge"`
	Timekeeping string `json:"timekeeping"`
	Cadence     string `json:"cadence"`
}
type UpdateCampaignRecordRequest struct {
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

func ValidateCampaign(u *CampaignRecord) bool { return true }

func NewCampaign(id int) *CampaignRecord {
	characters := make([]CharacterRecord, 0)
	return &CampaignRecord{id, "", "", false, "", "", "", time.Now(), time.Now(), time.Now(), characters}
}

func NewCampaignFull(id int, name, clientId, judge, timekeeping, cadence string, recruit bool, cdate, udate, ldate time.Time, characters []CharacterRecord) *CampaignRecord {
	return &CampaignRecord{id, name, clientId, recruit, judge, timekeeping, cadence, time.Now(), time.Now(), time.Now(), characters}
}

func NewCampaignInsertion(name, clientId, judge, timekeeping, cadence string, recruit bool, cdate, udate, ldate time.Time, characters []CharacterRecord) *CampaignRecord {
	return &CampaignRecord{-1, clientId, name, recruit, judge, timekeeping, cadence, time.Now(), time.Now(), time.Now(), characters}
}
