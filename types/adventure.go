package types

import (
	"math"
	"time"
)

type Adventure struct {
	ID             int                  `json:"id"`
	CampaignId     int                  `json:"campaign_id"`
	Name           string               `json:"name"`
	TotalXPAmount  int                  `json:"total_xp"`
	NumberOfShares float64              `json:"total_shares"`
	FullShareXP    int                  `json:"full_share"`
	HalfShareXP    int                  `json:"half_share"`
	AdventureDate  time.Time            `json:"adventure_date"`
	Coins          Coins                `json:"coins"`
	Gems           []Gem                `json:"gems"`
	Jewellery      []Jewellery          `json:"jewellery"`
	Combat         []MonsterGroup       `json:"combat"`
	MagicItems     []MagicItem          `json:"magic_items"`
	Characters     []AdventureCharacter `json:"characters"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
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
	Coins      Coins
	Gems       []Gem
	Jewellery  []Jewellery
	Combat     []MonsterGroup
	MagicItems []MagicItem
}

func NewAdventureRecord(id, campId int, c Coins, g []Gem, j []Jewellery, mo []MonsterGroup, mi []MagicItem, ch []AdventureCharacter, name string, cdate, udate, adate time.Time) *Adventure {
	newAdventure := Adventure{
		ID:            id,
		CampaignId:    campId,
		Coins:         c,
		Gems:          g,
		Jewellery:     j,
		MagicItems:    mi,
		Combat:        mo,
		Characters:    ch,
		CreatedAt:     cdate,
		UpdatedAt:     udate,
		AdventureDate: adate,
	}
	newAdventure.NumberOfShares = newAdventure.calculateNumberOfShares()
	newAdventure.TotalXPAmount = newAdventure.calculateTotalXP()
	newAdventure.FullShareXP, newAdventure.HalfShareXP = newAdventure.calculateXPShares()
	return &newAdventure
}

func NewAdventureRecordById(id int) *Adventure {
	gems := make([]Gem, 0)
	jewellery := make([]Jewellery, 0)
	magicItems := make([]MagicItem, 0)
	combats := make([]MonsterGroup, 0)
	chars := make([]AdventureCharacter, 0)
	return NewAdventureRecord(id, -1, *NewCoins(0, 0, 0, 0, 0), gems, jewellery, combats, magicItems, chars, "", time.Now(), time.Now(), time.Now())
}

func (a Adventure) calculateTotalXP() int {
	totalXp := a.Coins.TotalXPAmount + a.totalGemXp() + a.totalJewelleryXp() + a.totalMagicItemXp() + a.totalCombatXp()
	return totalXp
}

func (a Adventure) calculateXPShares() (int, int) {
	numberOfShares := a.calculateNumberOfShares()
	totalXp := a.Coins.TotalXPAmount + a.totalGemXp() + a.totalJewelleryXp() + a.totalMagicItemXp() + a.totalCombatXp()
	fullShareXP := math.RoundToEven(float64(totalXp) / numberOfShares)
	halfShareXP := math.RoundToEven(fullShareXP / 2.0)
	return int(fullShareXP), int(halfShareXP)
}

func (a Adventure) calculateNumberOfShares() float64 {
	totalShares := 0.0
	for _, char := range a.Characters {
		if char.Halfshare {
			totalShares += 0.5
		} else {
			totalShares += 1.0
		}
	}
	return totalShares
}

func (a Adventure) totalGemXp() (xp int) {
	xp = 0
	for _, gem := range a.Gems {
		xp += gem.TotalXPAmount
	}

	return xp
}

func (a Adventure) totalJewelleryXp() (xp int) {
	xp = 0
	for _, jewellery := range a.Jewellery {
		xp += jewellery.TotalXPAmount
	}
	return xp
}

func (a Adventure) totalMagicItemXp() (xp int) {
	xp = 0
	for _, mi := range a.MagicItems {
		xp += mi.TotalXPAmount
	}

	return xp
}

func (a Adventure) totalCombatXp() (xp int) {
	xp = 0
	for _, c := range a.Combat {
		xp += c.TotalXPAmount
	}

	return xp
}
