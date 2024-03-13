package types

import "time"

type AdventureRecord struct {
	ID            int
	CampaignID    int
	Coins         *Coins
	Gems          []*Gem
	Jewellery     []*Jewellery
	MagicItems    []*MagicItem
	Monsters      []*MonsterGroup
	AdventureGate time.Time
}

type CreateAdventureRecordRequest struct {
	CampaignID int
	Coins      *Coins
	Gems       []*Gem
	Jewellery  []*Jewellery
	MagicItems []*MagicItem
	Monsters   []*MonsterGroup
}

type UpdateAdventureRecordRequest struct {
	ID         int
	CampaignID int
	Coins      *Coins
	Gems       []*Gem
	Jewellery  []*Jewellery
	MagicItems []*MagicItem
	Monsters   []*MonsterGroup
}

func NewAdventureRecord(id, campId int, c *Coins, g []*Gem, j []*Jewellery, mi []*MagicItem, mo []*MonsterGroup, date time.Time) *AdventureRecord {
	return &AdventureRecord{id, campId, c, g, j, mi, mo, date}
}
