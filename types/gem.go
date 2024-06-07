package types

type Gem struct {
	Loot XPSource
}

func NewGem(n, d string, v float64, t int) *Gem {
	return &Gem{*NewLoot(n, d, v, t)}
}

func (g *Gem) Name() string {
	return g.Loot.Name
}
