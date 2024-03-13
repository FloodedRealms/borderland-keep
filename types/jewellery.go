package types

import "fmt"

type Jewellery struct {
	Name        string
	Description string
	Value       int
}

func NewJewellery(n, d string, v int) *Jewellery {
	return &Jewellery{
		Name:        n,
		Description: d,
		Value:       v,
	}
}

func (j *Jewellery) XPValue() int {
	return j.Value
}

func (j *Jewellery) GoldValue() float64 {
	return float64(j.Value)
}

func (g *Jewellery) Summary() string {
	if g.Description == "" {
		return fmt.Sprintf("Recovered a %s. It is worth %d gold.", g.Name, g.Value)
	}
	return fmt.Sprintf("Recovered a %s worth %d gold. It %s.", g.Name, g.Value, g.Description)

}
