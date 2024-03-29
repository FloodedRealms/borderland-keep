package types

import "fmt"

type MagicItem struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	ApparentValue int    `json:"apparent_value"`
}

func NewMagicItem(n, d string, v int) *MagicItem {
	return &MagicItem{
		Name:          n,
		Description:   d,
		ApparentValue: v,
	}
}

func (m *MagicItem) TotalXPValue() int {
	return m.ApparentValue
}

func (g *MagicItem) Summary() string {
	if g.Description == "" {
		return fmt.Sprintf("Recovered a %s. It looks worth %d gold.", g.Name, g.ApparentValue)
	}
	return fmt.Sprintf("Recovered a %s, apparently worth %d gold. It %s.", g.Name, g.ApparentValue, g.Description)

}
