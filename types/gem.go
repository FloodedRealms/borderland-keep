package types

type Gem struct {
	Id       int `json:"id"`
	XPEarned int `json:"gem_xp"`
	Loot     Loot
}

func NewGem(n, d string, v float64, number, id int) *Gem {
	gem := Gem{
		Id:   id,
		Loot: *NewLoot(n, d, v, number),
	}
	gem.XPEarned = gem.TotalXPAmount()
	return &gem
}

func (g *Gem) Name() string {
	return g.Loot.Name
}

func (g *Gem) TotalXPAmount() int {
	return g.Loot.TotalXPAmount()

}
