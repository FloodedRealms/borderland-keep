package types

import "time"

type AdventureRecord struct {
	ID            int            `json:"id"`
	CampaignID    int            `json:"campaign_id"`
	Copper        int            `json:"copper"`
	Silver        int            `json:"silver"`
	Electrum      int            `json:"electrum"`
	Gold          int            `json:"gold"`
	Platinum      int            `json:"platinum"`
	Gems          []Gem          `json:"gems"`
	Jewellery     []Jewellery    `json:"jewellery"`
	Combat        []MonsterGroup `json:"combat"`
	MagicItems    []MagicItem    `json:"magic_items"`
	Characters    []Character    `json:"characters"`
	Name          string         `json:"name"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	AdventureDate time.Time      `json:"adventure_date"`
}

type CreateAdventureRecordRequest struct {
	CampaignID int `json:"campaign_id"`
	//Copper        int             `json:"copper"`
	//Silver        int             `json:"silver"`
	//Electrum      int             `json:"electrum"`
	//Gold          int             `json:"gold"`
	//Platinum      int             `json:"platinum"`
	//Gems          []*Gem          `json:"gems"`
	//Jewellery     []*Jewellery    `json:"jewellery"`
	//Combat        []*MonsterGroup `json:"combat"`
	//MagicItems    []*MagicItem    `json:"magic_items"`
	//Characters    []*Character    `json:"characters"`
	Name          string    `json:"name"`
	AdventureDate time.Time `json:"adventure_date"`
}

type UpdateAdventureRecordRequest struct {
	ID         int
	CampaignID int
	Coins      *Coins
	Gems       []Gem
	Jewellery  []Jewellery
	Combat     []MonsterGroup
	MagicItems []MagicItem
}

func NewAdventureRecord(id, campId int, g []Gem, j []Jewellery, mo []MonsterGroup, mi []MagicItem, ch []Character, name string, cdate, udate, adate time.Time) *AdventureRecord {
	return &AdventureRecord{id, campId, 0, 0, 0, 0, 0, g, j, mo, mi, ch, name, cdate, udate, adate}
}

func NewAdventureRecordById(id int) *AdventureRecord {
	return NewAdventureRecord(id, -1, nil, nil, nil, nil, nil, "", time.Now(), time.Now(), time.Now())
}
