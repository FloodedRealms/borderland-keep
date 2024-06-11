package types

import (
	"fmt"
	"math"
)

type XPSource struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	XPValue     float64 `json:"value"`
	Number      int     `json:"total"`
}

type LootCategory string

const (
	CoinLoot      LootCategory = "coins"
	GemLoot       LootCategory = "gem"
	JewelleryLoot LootCategory = "jewellery"
	MagicItemLoot LootCategory = "magicItem"
	Combat        LootCategory = "combat"
)

func NewLoot(name, desc string, xpValue float64, numberOfSource int) *XPSource {
	source := &XPSource{
		Name:        name,
		Description: desc,
		XPValue:     xpValue,
		Number:      numberOfSource,
	}
	return source
}

func (g *XPSource) CalculateTotalXPValue() int {
	value := g.XPValue * float64(g.Number)
	rounded := math.Floor(value)
	final := int(rounded)
	return final
}

func (g *XPSource) GoldValue() float64 {
	return math.Floor((g.XPValue * float64(g.Number)))

}

func (g *XPSource) Summary() string {
	if g.Description == "" {
		return fmt.Sprintf("Recovered a %s. It is worth %f gold.", g.Name, g.XPValue)
	}
	return fmt.Sprintf("Recovered a %s worth %f gold. It %s.", g.Name, g.XPValue, g.Description)

}
