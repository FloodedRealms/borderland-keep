package types

type MagicItem struct {
	Loot        *XPSource
	ActualValue int `json:"actual_value"`
}

func NewMagicItem(n, d string, v float64, av int) *MagicItem {
	return &MagicItem{
		Loot:        NewLoot(n, d, v, 1),
		ActualValue: av,
	}
}

func (m MagicItem) ApparentValue() int {
	return int(m.Loot.XPValue)
}
