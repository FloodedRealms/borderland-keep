package types

type Jewellery struct {
	Loot *XPSource
}

func NewJewellery(n, d string, v float64, t int) *Jewellery {
	return &Jewellery{NewLoot(n, d, v, t)}
}
