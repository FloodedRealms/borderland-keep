package types

type MonsterGroup struct {
	XP XPSource
}

func NewMonsterGroup(name string, numberDefeated int, xpValue float64) *MonsterGroup {
	return &MonsterGroup{
		XP: *NewLoot(name, "vicous monsters", xpValue, numberDefeated),
	}
}
