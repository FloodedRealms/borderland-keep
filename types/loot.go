package types

import (
	"fmt"
	"math"
)

type LootCategory string

const (
	CoinLoot      LootCategory = "coins"
	GemLoot       LootCategory = "gem"
	JewelleryLoot LootCategory = "jewellery"
	MagicItemLoot LootCategory = "magicItem"
	Combat        LootCategory = "combat"
)

type Loot struct {
	Category       LootCategory `json:"type_of_loot"`
	Name           string       `json:"name"`
	Description    string       `json:"description"`
	XPValueOfOne   float64      `json:"xp_value"`
	GoldValueOfOne float64      `json:"gold_value"`
	NumberOfItem   int          `json:"total"`
}

func NewLoot(name, desc string, xpValue float64, numberOfSource int) *Loot {
	source := &Loot{
		Name:         name,
		Description:  desc,
		XPValueOfOne: xpValue,
		NumberOfItem: numberOfSource,
	}
	return source
}

func (g *Loot) TotalXPAmount() int {
	value := g.XPValueOfOne * float64(g.NumberOfItem)
	rounded := math.Floor(value)
	final := int(rounded)
	return final
}

func (g *Loot) GoldValue() float64 {
	return math.Floor((g.GoldValueOfOne * float64(g.NumberOfItem)))

}

func (g *Loot) Summary() string {
	switch g.Category {
	case CoinLoot:
		return fmt.Sprintf("Recovered %d %s coins. This is worth %f gold and %d XP.", g.NumberOfItem, g.Name, g.GoldValue(), g.TotalXPAmount())
	case Combat:
		return fmt.Sprintf("Defeated %d %s, earning %d XP.", g.NumberOfItem, g.Name, g.TotalXPAmount())
	default:
		return fmt.Sprintf("Recovered a %s. It is worth %f gold and %d XP.", g.Name, g.GoldValue(), g.TotalXPAmount())
	}
}
