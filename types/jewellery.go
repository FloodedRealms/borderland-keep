package types

type Jewellery struct {
	Id       int `json:"id"`
	XPEarned int `json:"jewellery_xp"`
	Loot     Loot
}

func NewJewellery(n, d string, v float64, number, id int) *Jewellery {
	j := &Jewellery{
		Id:   id,
		Loot: *NewLoot(n, d, v, number)}
	j.XPEarned = j.Loot.TotalXPAmount()
	return j
}
