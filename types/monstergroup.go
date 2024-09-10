package types

import "strconv"

type MonsterGroup struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	XPPerOneKill   int    `json:"xp"`
	NumberDefeated int    `json:"defeated"`
}

type MonsterGroupRequest struct {
	MonsterName    string `json:"monster_name"`
	XPPerMonster   int    `json:"xp"`
	NumberDefeated int    `json:"number_defeated"`
}

func NewMonsterGroup(name, d string, numberDefeated, xp int) *MonsterGroup {
	mon := MonsterGroup{
		Name:           name,
		Description:    d,
		XPPerOneKill:   xp,
		NumberDefeated: numberDefeated,
	}
	return &mon

}

func (g MonsterGroup) TotalXPAmount() float64 {
	return float64(g.XPPerOneKill * g.NumberDefeated)

}

func (g MonsterGroup) DisplayXPAmount() int {
	return int(g.XPPerOneKill)
}

func (g MonsterGroup) DisplayTotalXPAmount() int {
	return int(g.TotalXPAmount())
}

func (g MonsterGroup) DisplayNumber() string {
	return strconv.Itoa(g.NumberDefeated)
}
