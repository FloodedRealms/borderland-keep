package types

import (
	"log"
	"math"
)

type MonsterGroup struct {
	Id            int `json:"id"`
	TotalXPAmount int `json:"combat_xp"`
	XP            XPSource
}

type MonsterGroupRequest struct {
	MonsterName    string `json:"monster_name"`
	XPPerMonster   int    `json:"xp"`
	NumberDefeated int    `json:"number_defeated"`
}

func NewMonsterGroup(name string, numberDefeated, id int, xpValue float64) *MonsterGroup {
	mon := MonsterGroup{
		Id: id,
		XP: *NewLoot(name, "vicous monsters", xpValue, numberDefeated),
	}

	mon.TotalXPAmount = mon.XP.CalculateTotalXPValue()
	return &mon

}

func (g *MonsterGroup) calculateTotalXP() int {
	value := g.XP.XPValue * float64(g.XP.Number)
	log.Printf("Calculating XP for Gem named %s. Individual worth %f, number of items %d. Got a value of %f", g.XP.Name, g.XP.XPValue, g.XP.Number, value)
	rounded := math.Floor(value)
	log.Printf("Rounded to %f", rounded)
	final := int(rounded)
	log.Printf("Int value to %d", final)
	return final

}
