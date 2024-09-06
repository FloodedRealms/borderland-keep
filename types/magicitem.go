package types

type MagicItem struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"Description"`
	XPValue     int    `json:"magic_item_xp"`
	GoldValue   int    `json:"actual_value"`
}

func NewMagicItem(n, d string, v, av int) *MagicItem {
	return &MagicItem{
		Name:        n,
		Description: d,
		XPValue:     v,
		GoldValue:   av,
	}
}

func (m MagicItem) TotalGoldAmount() int {
	return m.GoldValue
}

func (m MagicItem) TotalXPAmount() float64 {
	return float64(m.XPValue)
}
