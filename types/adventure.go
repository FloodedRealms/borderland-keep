package types

import (
	"encoding/json"
	"math"
)

type AdventureRecord struct {
	Id            int                  `json:"id"`
	CampaignId    int                  `json:"campaign_id"`
	Name          string               `json:"name"`
	FullShareXP   int                  `json:"full_share"`
	HalfShareXP   int                  `json:"half_share"`
	AdventureDate ArcvhistDate         `json:"adventure_date"`
	GameDays      int                  `json:"duration"`
	Coins         Coins                `json:"coins"`
	Gems          []Gem                `json:"gems"`
	Jewellery     []Jewellery          `json:"jewellery"`
	Combat        []MonsterGroup       `json:"combat"`
	MagicItems    []MagicItem          `json:"magic_items"`
	Characters    []AdventureCharacter `json:"characters"`
}

func NewAdventureRecordById(id int) *AdventureRecord {
	return &AdventureRecord{Id: id}
}

func NewAdventureRecord(id, cid, days int, coins Coins, g []Gem, j []Jewellery, c []MonsterGroup, m []MagicItem, char []AdventureCharacter, n string, date ArcvhistDate) *AdventureRecord {
	a := &AdventureRecord{
		Id:            id,
		CampaignId:    cid,
		Name:          n,
		AdventureDate: date,
		GameDays:      days,
		Coins:         coins,
		Gems:          g,
		Jewellery:     j,
		Combat:        c,
		MagicItems:    m,
		Characters:    char,
	}
	f, h := a.CalculateXPShares()
	a.FullShareXP = f
	a.HalfShareXP = h
	return a
}

func (a *AdventureRecord) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	var incomingRequest struct {
		Id            int                  `json:"id"`
		CampaignId    int                  `json:"campaign_id"`
		Name          string               `json:"name"`
		AdventureDate ArcvhistDate         `json:"adventure_date"`
		GameDays      int                  `json:"duration"`
		Coins         Coins                `json:"coins"`
		Gems          []Gem                `json:"gems"`
		Jewellery     []Jewellery          `json:"jewellery"`
		Combat        []MonsterGroup       `json:"combat"`
		MagicItems    []MagicItem          `json:"magic_items"`
		Characters    []AdventureCharacter `json:"characters"`
	}
	if err := json.Unmarshal(data, &incomingRequest); err != nil {
		return nil
	}
	a.Id = incomingRequest.Id
	a.CampaignId = incomingRequest.CampaignId
	a.Name = incomingRequest.Name
	a.AdventureDate = incomingRequest.AdventureDate
	a.GameDays = incomingRequest.GameDays
	a.Coins = incomingRequest.Coins
	a.Gems = incomingRequest.Gems
	a.Jewellery = incomingRequest.Jewellery
	a.Combat = incomingRequest.Combat
	a.MagicItems = incomingRequest.MagicItems
	a.Characters = incomingRequest.Characters
	f, h := a.CalculateXPShares()
	a.FullShareXP = f
	a.HalfShareXP = h
	return nil
}

func (a AdventureRecord) TotalXPAmount() int {
	coins := a.totalCoinXp()
	gems := a.totalGemXp()
	jewells := a.totalJewelleryXp()
	magic := a.totalMagicItemXp()
	combat := a.totalCombatXp()
	totalXp := coins + gems + jewells + magic + combat
	return int(math.RoundToEven(totalXp))
}

func (a AdventureRecord) CalculateXPShares() (fullshare, halfsare int) {
	numberOfShares := a.CalculateNumberOfShares()
	c := a.Coins.TotalXPAmount()
	g := a.totalGemXp()
	j := a.totalJewelleryXp()
	m := a.totalMagicItemXp()
	com := a.totalCombatXp()
	totalXp := c + g + j + m + com
	if totalXp == 0 {
		return 0, 0
	}
	fullShareXP := math.RoundToEven(float64(totalXp) / numberOfShares)
	halfShareXP := math.RoundToEven(fullShareXP / 2.0)
	return int(fullShareXP), int(halfShareXP)
}

func (a AdventureRecord) CalculateNumberOfShares() float64 {
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

func (a AdventureRecord) totalCoinXp() float64 {
	return a.Coins.TotalXPAmount()
}

func (a AdventureRecord) totalGemXp() float64 {
	xp := 0.0
	for _, gem := range a.Gems {
		xp += gem.TotalXPAmount()
	}

	return xp
}

func (a AdventureRecord) totalJewelleryXp() float64 {
	xp := 0.0
	for _, jewellery := range a.Jewellery {
		xp += jewellery.TotalXPAmount()
	}
	return xp
}

func (a AdventureRecord) totalMagicItemXp() float64 {
	xp := 0.0
	for _, mi := range a.MagicItems {
		xp += mi.TotalXPAmount()
	}

	return xp
}

func (a AdventureRecord) totalCombatXp() float64 {
	xp := 0.0
	for _, c := range a.Combat {
		xp += c.TotalXPAmount()
	}

	return xp
}

func (a AdventureRecord) CoinXPForType(c string) int {

	ans := -10.0
	switch c {
	case "Copper":
		ans = a.Coins.Copper.TotalXPAmount()
	case "Silver":
		ans = a.Coins.Silver.TotalXPAmount()
	case "Electrum":
		ans = a.Coins.Electrum.TotalXPAmount()
	case "Gold":
		ans = a.Coins.Gold.TotalXPAmount()
	case "Platinum":
		ans = a.Coins.Platinum.TotalXPAmount()
	}
	return int(ans)
}

/*
func (r AdventureRecord) GenerateGemList() []GenericLoot {
	gems := make([]GenericLoot, 0)
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
*/
