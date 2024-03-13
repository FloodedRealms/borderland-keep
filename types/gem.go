package types

import "fmt"

type Gem struct {
	Name        string
	Description string
	Value       int
}

func NewGem(n, d string, v int) *Gem {
	return &Gem{
		Name:        n,
		Description: d,
		Value:       v,
	}
}

func (g *Gem) XPValue() int {
	return g.Value
}

func (g *Gem) GoldValue() float64 {
	return float64(g.Value)
}

func (g *Gem) Summary() string {
	if g.Description == "" {
		return fmt.Sprintf("Recovered a %s. It is worth %d gold.", g.Name, g.Value)
	}
	return fmt.Sprintf("Recovered a %s worth %d gold. It %s.", g.Name, g.Value, g.Description)

}
