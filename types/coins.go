package types

import (
	"fmt"
	"math"
)

type Coins struct {
	Copper   int
	Silver   int
	Electrum int
	Gold     int
	Platinum int
}

func copperToGold(cp int) float64 {
	return float64(cp) / 100.0
}
func silverToGold(sp int) float64 {
	return float64(sp) / 10.0
}
func electrumToGold(ep int) float64 {
	return float64(ep) / 2.0
}
func platinumToGold(pp int) float64 {
	return float64(pp) * 5.0
}

func NewCoins(c, s, e, g, p int) *Coins {
	return &Coins{
		Copper:   c,
		Silver:   s,
		Electrum: e,
		Gold:     g,
		Platinum: p,
	}
}

func (c *Coins) GoldValue() float64 {
	return copperToGold(c.Copper) + silverToGold(c.Silver) + electrumToGold(c.Electrum) + float64(c.Gold) + platinumToGold(c.Platinum)
}

func (c *Coins) XPValue() int {
	return int(math.Floor(c.GoldValue()))
}

func (c *Coins) Summary() string {
	return fmt.Sprintf("Recovered %d Copper, %d Silver, %d Electrum, %d Gold, and %d Platinum.", c.Copper, c.Silver, c.Electrum, c.Gold, c.Platinum)
}
