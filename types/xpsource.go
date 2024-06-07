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
	GemLoot       LootCategory = "gem"
	JewelleryLoot LootCategory = "jewellery"
	MagicItemLoot LootCategory = "magicItem"
)

func NewLoot(name, desc string, xpValue float64, numberOfSource int) *XPSource {
	return &XPSource{
		Name:        name,
		Description: desc,
		XPValue:     xpValue,
		Number:      numberOfSource,
	}
}

func (g *XPSource) TotalXPValue() int {
	return int(g.XPValue)
}

func (g *XPSource) GoldValue() float64 {
	return math.Floor((g.XPValue * float64(g.Number)))

}

func (g *XPSource) Summary() string {
	if g.Description == "" {
		return fmt.Sprintf("Recovered a %s. It is worth %d gold.", g.Name, g.XPValue)
	}
	return fmt.Sprintf("Recovered a %s worth %d gold. It %s.", g.Name, g.XPValue, g.Description)

}
