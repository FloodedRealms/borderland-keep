package types

type MagicItem struct {
	Id          int `json:"id"`
	XPEarned    int `json:"magic_item_xp"`
	Loot        Loot
	ActualValue int `json:"actual_value"`
}

type MagicItemRequest struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ApparentValue int    `json:"apparent_value"`
	ActualValue   int    `json:"actual_value"`
}

func NewMagicItem(n, d string, v float64, av, id int) *MagicItem {
	mi := &MagicItem{
		Id:          id,
		Loot:        *NewLoot(n, d, v, 1),
		ActualValue: av,
	}
	mi.XPEarned = mi.Loot.TotalXPAmount()
	return mi
}

func (m *MagicItem) ApparentValue() int {
	return int(m.Loot.XPValueOfOne)
}

type IncomingMagicItem struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ApparentValue int    `json:"apparent_value"`
	ActualValue   int    `json:"actual_value"`
}
