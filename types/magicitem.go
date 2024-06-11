package types

type MagicItem struct {
	Id            int `json:"id"`
	TotalXPAmount int `json:"magic_item_xp"`
	Loot          XPSource
	ActualValue   int `json:"actual_value"`
}

type MagicItemRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ApparentValue int    `json:"apparent_value"`
	ActualValue   int    `json:"actual_value"`
}

func NewMagicItem(n, d string, v float64, av, id int) *MagicItem {
	return &MagicItem{
		Id:          id,
		Loot:        *NewLoot(n, d, v, 1),
		ActualValue: av,
	}
}

func (m MagicItem) ApparentValue() int {
	return int(m.Loot.XPValue)
}
