package types

import "fmt"

type MonsterGroup struct {
	MonsterName    string `json:"monster_name"`
	NumberDefeated int    `json:"number_defeated"`
	XP             int    `json:"xp_per_monster"`
	TotalXP        int    `json:"total_xp"`
}

func NewMonsterGroup(n string, d, v int) *MonsterGroup {
	group := &MonsterGroup{
		MonsterName:    n,
		NumberDefeated: d,
		XP:             v,
	}
	group.TotalXP = group.NumberDefeated * group.XP
	return group
}

func (m *MonsterGroup) XPValue() int {
	return m.TotalXP
}

func (m *MonsterGroup) Summary() string {
	if m.NumberDefeated == 1 {
		return fmt.Sprintf("Defeated a vile %s. This brings %d XP.", m.MonsterName, m.TotalXP)
	}
	return fmt.Sprintf("Defeated a group of vile %ss. This brings %d XP.", m.MonsterName, m.TotalXP)
}
