package types

import (
	"log"
	"math"
)

type Gem struct {
	Id            int `json:"id"`
	TotalXPAmount int `json:"gem_xp"`
	Loot          XPSource
}

func NewGem(n, d string, v float64, number, id int) *Gem {
	gem := Gem{
		Id:   id,
		Loot: *NewLoot(n, d, v, number),
	}
	gem.TotalXPAmount = gem.calculateTotalXP()
	return &gem
}

func (g *Gem) Name() string {
	return g.Loot.Name
}

func (g *Gem) calculateTotalXP() int {
	value := g.Loot.XPValue * float64(g.Loot.Number)
	log.Printf("Calculating XP for Gem named %s. Individual worth %f, number of items %d. Got a value of %f", g.Name(), g.Loot.XPValue, g.Loot.Number, value)
	rounded := math.Floor(value)
	log.Printf("Rounded to %f", rounded)
	final := int(rounded)
	log.Printf("Int value to %d", final)
	return final

}
