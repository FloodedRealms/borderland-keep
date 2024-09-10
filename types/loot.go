package types

import (
	"math"
	"strconv"
)

type GenericLootType string

const (
	CoinLoot      GenericLootType = "coin"
	GemLoot       GenericLootType = "gem"
	JewelleryLoot GenericLootType = "jewellery"
	CombatLoot    GenericLootType = "combat"
	MagicItemLoot GenericLootType = "magicItem"
)

type GenericLoot struct {
	Id          int             `json:"id"`
	LootType    GenericLootType `json:"type"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Number      int             `json:"number"`
	XPValue     float64         `json:"xp_value"`
	GoldValue   float64         `json:"gold_value"`
}

// Basic Coins and Gems and Jewellery opaque types on genericLoot
type Copper GenericLoot
type Silver GenericLoot
type Electrum GenericLoot
type Gold GenericLoot
type Platinum GenericLoot
type Gem GenericLoot
type Jewellery GenericLoot

func (c Copper) TotalGoldAmount() float64 {
	return GenericLoot(c).TotalGoldAmount()
}
func (c Copper) TotalXPAmount() float64 {
	return GenericLoot(c).TotalXPAmount()
}
func (c Silver) TotalGoldAmount() float64 {
	return GenericLoot(c).TotalGoldAmount()
}
func (c Silver) TotalXPAmount() float64 {
	return GenericLoot(c).TotalXPAmount()
}
func (c Electrum) TotalGoldAmount() float64 {
	return GenericLoot(c).TotalGoldAmount()
}
func (c Electrum) TotalXPAmount() float64 {
	return GenericLoot(c).TotalXPAmount()
}
func (c Gold) TotalGoldAmount() float64 {
	return GenericLoot(c).TotalGoldAmount()
}
func (c Gold) TotalXPAmount() float64 {
	return GenericLoot(c).TotalXPAmount()
}
func (c Platinum) TotalGoldAmount() float64 {
	return GenericLoot(c).TotalGoldAmount()
}
func (c Platinum) TotalXPAmount() float64 {
	return GenericLoot(c).TotalXPAmount()
}
func (c Gem) TotalGoldAmount() float64 {
	return GenericLoot(c).TotalGoldAmount()
}
func (c Gem) TotalXPAmount() float64 {
	return GenericLoot(c).TotalXPAmount()
}
func (c Jewellery) TotalGoldAmount() float64 {
	return GenericLoot(c).TotalGoldAmount()
}
func (c Jewellery) TotalXPAmount() float64 {
	return GenericLoot(c).TotalXPAmount()
}

/*
	 For now, we are hardcoding the ACKS assumptions into the system.
		* That means that each coin demarcation is an XP source with an XPValue equal to its
		* value in gold according to ACKS and a Number equal to the normal coins included in the adventure.
		* These values are:
		* Copper = 1/100
		* Silver = 1/10
		* Electrum = 1/2
		* Gold = 1
		* Platinum = 5

		TODO: Front ends should be able to change this as they desire, i.e. a front end written to handle 7voz would be on the silver standard,
		so silver has an XP value of 1. A front end for d&d 5e would have all coin XPValues at 0, since treasure does not award XP
*/

func NewCopper(a int) *Copper {
	return &Copper{
		LootType:    CoinLoot,
		Name:        "Copper",
		Description: "Small Coin",
		Number:      a,
		XPValue:     .01,
		GoldValue:   .01,
	}
}
func NewSilver(a int) *Silver {
	return &Silver{
		LootType:    CoinLoot,
		Description: "Don't spend it all in one place.",
		Name:        "Silver",
		Number:      a,
		XPValue:     .1,
		GoldValue:   .1,
	}
}
func NewElectrum(a int) *Electrum {
	return &Electrum{
		LootType:    CoinLoot,
		Description: "Best Coin",
		Name:        "Electrum",
		Number:      a,
		XPValue:     .5,
		GoldValue:   .5,
	}
}
func NewGold(a int) *Gold {
	return &Gold{
		LootType:    CoinLoot,
		Description: "THE standard",
		Name:        "Gold",
		Number:      a,
		XPValue:     1,
		GoldValue:   1,
	}
}

func NewPlatinum(a int) *Platinum {
	return &Platinum{
		LootType:    CoinLoot,
		Description: "Big Spender",
		Name:        "Platinum",
		Number:      a,
		XPValue:     5,
		GoldValue:   5,
	}
}
func NewGem(number int, name, d string, xpv, gv float64) *Gem {
	return &Gem{
		LootType:    GemLoot,
		Description: d,
		Name:        name,
		Number:      number,
		XPValue:     xpv,
		GoldValue:   gv,
	}
}
func NewJewellery(number int, name, d string, xpv, gv float64) *Jewellery {
	return &Jewellery{
		LootType:    JewelleryLoot,
		Description: d,
		Name:        name,
		Number:      number,
		XPValue:     xpv,
		GoldValue:   gv,
	}
}

func (c GenericLoot) TotalXPAmount() float64 {
	return c.XPValue * float64(c.Number)
}

func (c GenericLoot) TotalGoldAmount() float64 {
	return math.Round(c.GoldValue * float64(c.Number))
}

func (c GenericLoot) DisplayName() string {
	return c.Name
}
func (c GenericLoot) DisplayGPAmount() int {
	return int(c.GoldValue)
}
func (c GenericLoot) DisplayXPAmount() int {
	return int(c.XPValue)
}

func (c GenericLoot) DisplayTotalGPAmount() int {
	return int(c.TotalGoldAmount())
}
func (c GenericLoot) DisplayTotalXPAmount() int {
	return int(c.TotalXPAmount())
}

func (c GenericLoot) DisplayNumber() string {
	return strconv.Itoa(c.Number)
}

type Coins struct {
	Copper   Copper   `json:"copper"`
	Silver   Silver   `json:"silver"`
	Electrum Electrum `json:"electrum"`
	Gold     Gold     `json:"gold"`
	Platinum Platinum `json:"platinum"`
}

func NewCoins(c, s, e, g, p int) *Coins {

	copper := NewCopper(c)
	silver := NewSilver(s)
	electrum := NewElectrum(e)
	gold := NewGold(g)
	platinum := NewPlatinum(p)

	newCoin := Coins{
		Copper:   *copper,
		Silver:   *silver,
		Electrum: *electrum,
		Gold:     *gold,
		Platinum: *platinum,
	}
	return &newCoin
}

func (c Coins) TotalXPAmount() float64 {
	return c.Copper.TotalXPAmount() + c.Silver.TotalXPAmount() + c.Electrum.TotalXPAmount() + c.Gold.TotalXPAmount() + c.Platinum.TotalXPAmount()
}
func (c Coins) TotalGoldAmount() float64 {
	return c.Copper.TotalGoldAmount() + c.Silver.TotalGoldAmount() + c.Electrum.TotalGoldAmount() + c.Gold.TotalGoldAmount() + c.Platinum.TotalGoldAmount()
}
