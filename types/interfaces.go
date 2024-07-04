package types

type XPSource interface {
	TotalXPAmount() float64
}

type GoldSource interface {
	TotalValueInGold() int
}

type Loot interface {
	XPSource
	GoldSource
}
