package types

import (
	"log"
	"math"
	"time"
)

type AdventureRecord struct {
	ID             int                  `json:"id"`
	CampaignId     int                  `json:"campaign_id"`
	Name           string               `json:"name"`
	TotalXPAmount  int                  `json:"total_xp"`
	NumberOfShares float64              `json:"total_shares"`
	FullShareXP    int                  `json:"full_share"`
	HalfShareXP    int                  `json:"half_share"`
	AdventureDate  time.Time            `json:"adventure_date"`
	GameDays       int                  `json:"duration"`
	Coins          Coins                `json:"coins"`
	Gems           []Gem                `json:"gems"`
	Jewellery      []Jewellery          `json:"jewellery"`
	Combat         []MonsterGroup       `json:"combat"`
	MagicItems     []MagicItem          `json:"magic_items"`
	Characters     []AdventureCharacter `json:"characters"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

type CreateAdventureRequest struct {
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

func NewAdventureRecord(id, campId, duration int, c Coins, g []Gem, j []Jewellery, mo []MonsterGroup, mi []MagicItem, ch []AdventureCharacter, name string, cdate, udate, adate time.Time) *AdventureRecord {
	newAdventure := AdventureRecord{
		ID:            id,
		CampaignId:    campId,
		Name:          name,
		Coins:         c,
		Gems:          g,
		Jewellery:     j,
		MagicItems:    mi,
		Combat:        mo,
		Characters:    ch,
		CreatedAt:     cdate,
		UpdatedAt:     udate,
		AdventureDate: adate,
		GameDays:      duration,
	}
	newAdventure.NumberOfShares = newAdventure.calculateNumberOfShares()
	newAdventure.TotalXPAmount = newAdventure.calculateTotalXP()
	newAdventure.FullShareXP, newAdventure.HalfShareXP = newAdventure.CalculateXPShares()
	return &newAdventure
}

func NewAdventureRecordById(id int) *AdventureRecord {
	gems := make([]Gem, 0)
	jewellery := make([]Jewellery, 0)
	magicItems := make([]MagicItem, 0)
	combats := make([]MonsterGroup, 0)
	chars := make([]AdventureCharacter, 0)
	return NewAdventureRecord(id, -1, 1, *NewCoins(0, 0, 0, 0, 0), gems, jewellery, combats, magicItems, chars, "", time.Now(), time.Now(), time.Now())
}

func (a AdventureRecord) calculateTotalXP() int {
	totalXp := a.Coins.TotalXPAmount + a.totalGemXp() + a.totalJewelleryXp() + a.totalMagicItemXp() + a.totalCombatXp()
	return totalXp
}

func (a AdventureRecord) CalculateXPShares() (fullshare, halfsare int) {
	numberOfShares := a.calculateNumberOfShares()
	totalXp := a.Coins.TotalXPAmount + a.totalGemXp() + a.totalJewelleryXp() + a.totalMagicItemXp() + a.totalCombatXp()
	if totalXp == 0 {
		return 0, 0
	}
	fullShareXP := math.RoundToEven(float64(totalXp) / numberOfShares)
	halfShareXP := math.RoundToEven(fullShareXP / 2.0)
	return int(fullShareXP), int(halfShareXP)
}

func (a AdventureRecord) calculateNumberOfShares() float64 {
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

func (a AdventureRecord) totalGemXp() (xp int) {
	xp = 0
	for _, gem := range a.Gems {
		xp += gem.TotalXPAmount()
	}

	return xp
}

func (a AdventureRecord) totalJewelleryXp() (xp int) {
	xp = 0
	for _, jewellery := range a.Jewellery {
		xp += jewellery.TotalXPAmount
	}
	return xp
}

func (a AdventureRecord) totalMagicItemXp() (xp int) {
	xp = 0
	for _, mi := range a.MagicItems {
		xp += mi.TotalXPAmount
	}

	return xp
}

func (a AdventureRecord) totalCombatXp() (xp int) {
	xp = 0
	for _, c := range a.Combat {
		xp += c.XPEarned
	}

	return xp
}

type UpdateAdventureRequest struct {
	ID            int                        `json:"id"`
	CampaignID    int                        `json:"campaign_id"`
	Copper        int                        `json:"copper"`
	Silver        int                        `json:"silver"`
	Electrum      int                        `json:"electrum"`
	Gold          int                        `json:"gold"`
	Platinum      int                        `json:"platinum"`
	Gems          []loot                     `json:"gems"`
	Jewellery     []loot                     `json:"jewellery"`
	Combat        []loot                     `json:"combat"`
	MagicItems    []incomingMagicItem        `json:"magic_items"`
	Characters    []UpdateAdventureCharacter `json:"characters"`
	Name          string                     `json:"name"`
	AdventureDate time.Time                  `json:"adventure_date"`
}

func (r UpdateAdventureRequest) GenerateGemList() []Gem {
	gems := make([]Gem, 0)
	for _, gem := range r.Gems {
		newGem := NewGem(gem.Name, gem.Description, gem.XPValueOfOne, gem.NumberOfItem, -1)
		gems = append(gems, *newGem)
	}
	return gems
}
func (r UpdateAdventureRequest) GenerateJewelleryList() []Jewellery {
	gems := make([]Jewellery, 0)
	for _, gem := range r.Jewellery {
		newGem := NewJewellery(gem.Name, gem.Description, gem.XPValueOfOne, gem.NumberOfItem, -1)
		gems = append(gems, *newGem)
	}
	return gems
}
func (r UpdateAdventureRequest) GenerateMagicItemList() []MagicItem {
	gems := make([]MagicItem, 0)
	for _, gem := range r.MagicItems {
		log.Print(gem)
		newItem := NewMagicItem(gem.Name, gem.Description, float64(gem.ApparentValue), gem.ActualValue, -1)
		gems = append(gems, *newItem)
	}
	return gems
}
func (r UpdateAdventureRequest) GenerateCombatList() []MonsterGroup {
	gems := make([]MonsterGroup, 0)
	for _, gem := range r.Combat {
		newGem := NewMonsterGroup(gem.Name, gem.NumberOfItem, -1, gem.XPValueOfOne)
		gems = append(gems, *newGem)
	}
	return gems
}
func (r UpdateAdventureRequest) GenerateCharacterList() []AdventureCharacter {
	characters := make([]AdventureCharacter, 0)
	for _, char := range r.Characters {
		newChar := NewAdventureCharacter(NewCharacterById(char.ID), char.Halfshare, char.XpGained)
		characters = append(characters, *newChar)
	}
	return characters
}
