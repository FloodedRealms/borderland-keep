package types

type Jewellery struct {
	Id            int `json:"id"`
	TotalXPAmount int `json:"jewellery_xp"`
	Loot          loot
}

func NewJewellery(n, d string, v float64, number, id int) *Jewellery {
	return &Jewellery{
		Id:   id,
		Loot: *NewLoot(n, d, v, number)}
}
