package types

import (
	"time"
)

type CampaignRecord struct {
	Id            int               `json:"id"`
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
	Adventures    []AdventureRecord `json:"adventures"`
}

func ValidateCampaign(u *CampaignRecord) bool { return true }

func NewCampaign(id int) *CampaignRecord {
	characters := make([]CharacterRecord, 0)
	adventures := make([]AdventureRecord, 0)
	return &CampaignRecord{id, "", "", false, "", "", "", time.Now(), time.Now(), time.Now(), characters, adventures}
}

func NewCampaignFull(id int, name, clientId, judge, timekeeping, cadence string, recruit bool, cdate, udate, ldate time.Time, characters []CharacterRecord, a []AdventureRecord) *CampaignRecord {
	return &CampaignRecord{id, name, clientId, recruit, judge, timekeeping, cadence, time.Now(), time.Now(), time.Now(), characters, a}
}

func NewCampaignInsertion(name, clientId, judge, timekeeping, cadence string, recruit bool, cdate, udate, ldate time.Time, characters []CharacterRecord) *CampaignRecord {
	return &CampaignRecord{-1, clientId, name, recruit, judge, timekeeping, cadence, time.Now(), time.Now(), time.Now(), characters, []AdventureRecord{}}
}

func (c CampaignRecord) NumberOfCharacters() int {
	return len(c.Characters)
}

func (c CampaignRecord) NumberOfAdventures() int {
	return len(c.Adventures)
}

type CampaignClassOption struct {
	ClassId   int
	ClassName string
}
