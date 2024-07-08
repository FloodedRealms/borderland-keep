package types

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
