package types

import "time"

type AdventureRecord struct {
	ID            int             `json:"id"`
	CampaignID    int             `json:"campaign_id"`
	Coins         *Coins          `json:"coin"`
	Gems          []*Gem          `json:"gems"`
	Jewellery     []*Jewellery    `json:"jewellery"`
	Combat        []*MonsterGroup `json:"combat"`
	MagicItems    []*MagicItem    `json:"magic_items"`
	Name          string          `json:"name"`
	AdventureDate time.Time       `json:"adventure_date"`
}

type CreateAdventureRecordRequest struct {
	CampaignID int `json:"campaign_id"`
}

type UpdateAdventureRecordRequest struct {
	ID         int
	CampaignID int
	Coins      *Coins
	Gems       []*Gem
	Jewellery  []*Jewellery
	Combat     []*MonsterGroup
	MagicItems []*MagicItem
}

func NewAdventureRecord(id, campId int, c *Coins, g []*Gem, j []*Jewellery, mo []*MonsterGroup, mi []*MagicItem, name string, date time.Time) *AdventureRecord {
	return &AdventureRecord{id, campId, c, g, j, mo, mi, name, date}
}
