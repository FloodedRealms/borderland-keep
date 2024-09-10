package types

type MagicItem struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"Description"`
	ApparentValue int    `json:"magic_item_xp"`
	ActualValue   int    `json:"actual_value"`
	Sold          bool
}

func NewMagicItem(n, d string, v, av int, sold bool) *MagicItem {
	return &MagicItem{
		Name:          n,
		Description:   d,
		ApparentValue: v,
		ActualValue:   av,
		Sold:          sold,
	}
}

func (m MagicItem) TotalGoldAmount() int {
	if m.Sold {
		return m.ActualValue + m.ApparentValue
	}
	return m.ApparentValue
}

func (m MagicItem) TotalXPAmount() float64 {
	if m.Sold {
		return float64(m.ActualValue + m.ApparentValue)
	}
	return float64(m.ApparentValue)
}

func (m MagicItem) DisplayXPAmount() int {
	return int(m.TotalXPAmount())
}

func (m MagicItem) DisplayGPAmount() int {
	return m.TotalGoldAmount()
}

func (m MagicItem) DisplayActualValue() int {
	return m.ActualValue
}

func (m MagicItem) DisplayApparentValue() int {
	return m.ApparentValue
}

func (m MagicItem) WasSold() bool {
	return m.Sold
}
