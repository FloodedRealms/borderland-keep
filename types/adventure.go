package types

import (
	"fmt"
	"math"
	"time"
)

type AdventureRecord struct {
	Id             int    `json:"id"`
	CampaignId     int    `json:"campaign_id"`
	Name           string `json:"name"`
	NumberOfShares float64
	TotalXP        float64
	TotalGP        float64
	FullShareXP    float64 `json:"full_share"`
	HalfShareXP    float64 `json:"half_share"`
	FullShareGP    float64
	HalfShareGP    float64
	AdventureDate  time.Time            `json:"adventure_date"`
	GameDays       int                  `json:"duration"`
	Coins          Coins                `json:"coins"`
	Gems           []Gem                `json:"gems"`
	Jewellery      []Jewellery          `json:"jewellery"`
	Combat         []MonsterGroup       `json:"combat"`
	MagicItems     []MagicItem          `json:"magic_items"`
	Characters     []AdventureCharacter `json:"characters"`
}

func NewAdventureRecordById(id int) *AdventureRecord {
	return &AdventureRecord{Id: id}
}

func NewAdventureRecordByNumberOfCharacters(numPlayerCharacters, numHenchmen int) *AdventureRecord {
	c := make([]AdventureCharacter, 0)
	for range numPlayerCharacters {
		c = append(c, *NewAdventureCharacter(false, 0))
	}
	for range numHenchmen {
		c = append(c, *NewAdventureCharacter(true, 0))
	}
	a := &AdventureRecord{
		Characters: c,
		Coins:      Coins{},
		Gems:       []Gem{},
		Jewellery:  []Jewellery{},
		Combat:     []MonsterGroup{},
		MagicItems: []MagicItem{},
	}
	a.CalculateNumberOfShares()
	return a
}

func NewAdventureRecord(id, cid, days int, coins Coins, g []Gem, j []Jewellery, c []MonsterGroup, m []MagicItem, char []AdventureCharacter, n string, date time.Time) *AdventureRecord {
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
	a.CalculateXPShares()
	return a
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

func (a *AdventureRecord) CalculateTotalGP() {
	a.CalculateNumberOfShares()
	c := a.Coins.TotalGoldAmount()
	g := a.totalGemGp()
	j := a.totalJewelleryGp()
	m := a.totalMagicItemGp()
	a.TotalGP = c + g + j + m
}

func (a *AdventureRecord) CalculateGPShares() (fullshare, halfsare int) {
	a.CalculateTotalGP()
	if a.NumberOfShares <= 0 {
		a.FullShareGP = 0.0
		a.HalfShareGP = 0.0
		return
	}

	a.FullShareGP = math.RoundToEven(a.TotalGP / a.NumberOfShares)
	a.HalfShareGP = math.RoundToEven(a.FullShareGP / 2.0)
	return int(a.FullShareGP), int(a.HalfShareGP)
}

func (a *AdventureRecord) CalculateTotalXP() {
	a.CalculateNumberOfShares()

	c := a.Coins.TotalXPAmount()
	g := a.totalGemXp()
	j := a.totalJewelleryXp()
	m := a.totalMagicItemXp()
	com := a.totalCombatXp()
	a.TotalXP = c + g + j + m + com
}

func (a *AdventureRecord) CalculateXPShares() (fullshare, halfsare int) {
	a.CalculateTotalXP()
	if a.NumberOfShares <= 0 {
		a.FullShareXP = 0.0
		a.HalfShareXP = 0.0
		return
	}

	a.FullShareXP = math.RoundToEven(a.TotalXP / a.NumberOfShares)
	a.HalfShareXP = math.RoundToEven(a.FullShareXP / 2.0)
	return int(a.FullShareXP), int(a.HalfShareXP)
}

func (a *AdventureRecord) CalculateNumberOfShares() {
	totalShares := 0.0
	for _, char := range a.Characters {
		if char.Halfshare {
			totalShares += 0.5
		} else {
			totalShares += 1.0
		}
	}
	a.NumberOfShares = totalShares
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

func (a AdventureRecord) totalCoinGp() float64 {
	return a.Coins.TotalGoldAmount()
}

func (a AdventureRecord) totalGemGp() float64 {
	xp := 0.0
	for _, gem := range a.Gems {
		xp += gem.TotalGoldAmount()
	}

	return xp
}

func (a AdventureRecord) totalJewelleryGp() float64 {
	xp := 0.0
	for _, jewellery := range a.Jewellery {
		xp += jewellery.TotalGoldAmount()
	}
	return xp
}

func (a AdventureRecord) totalMagicItemGp() float64 {
	xp := 0
	for _, mi := range a.MagicItems {
		xp += mi.TotalGoldAmount()
	}

	return float64(xp)
}

func (a AdventureRecord) CoinGoldForType(c string) int {

	ans := -10.0
	switch c {
	case "Copper":
		ans = a.Coins.Copper.TotalGoldAmount()
	case "Silver":
		ans = a.Coins.Silver.TotalGoldAmount()
	case "Electrum":
		ans = a.Coins.Electrum.TotalGoldAmount()
	case "Gold":
		ans = a.Coins.Gold.TotalGoldAmount()
	case "Platinum":
		ans = a.Coins.Platinum.TotalGoldAmount()
	}
	return int(ans)
}

// Interface Fullfillment

func (a AdventureRecord) DisplayName() string {
	return a.Name
}

func (a AdventureRecord) DisplayTotalGPAmount() int {
	a.CalculateGPShares()
	return int(a.TotalGP)
}

func (a AdventureRecord) DisplayTotalXPAmount() int {
	a.CalculateXPShares()
	return int(a.TotalXP)
}

func (a AdventureRecord) DisplayDate() string {
	return a.AdventureDate.Format("January 01, 2006")
}

func (a AdventureRecord) DisplayTotalShares() string {
	return fmt.Sprintf("%.1f", a.NumberOfShares)
}

func (a AdventureRecord) DisplayFullGPShare() int {
	a.CalculateGPShares()
	return int(a.FullShareGP)
}

func (a AdventureRecord) DisplayHalfGPShare() int {
	a.CalculateGPShares()
	return int(a.HalfShareGP)
}

func (a AdventureRecord) DisplayFullXPShare() int {
	a.CalculateXPShares()
	return int(a.FullShareXP)
}

func (a AdventureRecord) DisplayHalfXPShare() int {
	a.CalculateXPShares()
	return int(a.HalfShareXP)
}
