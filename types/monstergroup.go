package types

type MonsterGroup struct {
	Id       int `json:"id"`
	XPEarned int `json:"combat_xp"`
	XP       loot
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

	mon.XPEarned = mon.XP.TotalXPAmount()
	return &mon

}

func (g *MonsterGroup) calculateTotalXP() int {
	return g.XP.TotalXPAmount()

}
