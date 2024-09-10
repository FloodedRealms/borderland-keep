package types

type XPSource interface {
	TotalXPAmount() float64
}

type GoldSource interface {
	TotalGoldAmount() float64
}

type ViewableXPSource interface {
	DisplayXPAmount() int
}

type ViewableGPSource interface {
	DisplayGPAmount() int
}

type ViewableMultiMemberXPSource interface {
	DisplayTotalXPAmount() int
}

type ViewableMultiMemberGPSource interface {
	DisplayTotalGPAmount() int
}

type Loot interface {
	XPSource
	GoldSource
	ViewableMetadata
	ViewableXPSource
	ViewableGPSource
	ViewableMultiMemberGPSource
	ViewableMultiMemberXPSource
	DisplayNumber() string
}

type Combat interface {
	ViewableMultiMemberXPSource
	ViewableXPSource
	DisplayNumber() string
}

type ViewableMetadata interface {
	DisplayName() string
}

type Adventure interface {
	ViewableMetadata
	ViewableMultiMemberGPSource
	ViewableMultiMemberXPSource
	DisplayDate() string
	DisplayTotalShares() string
	DisplayFullXPShare() int
	DisplayFullGPShare() int
	DisplayHalfXPShare() int
	DisplayHalfGPShare() int
}

type MagicalLoot interface {
	ViewableGPSource
	ViewableXPSource
	DisplayActualValue() int
	DisplayApparentValue() int
	WasSold() bool
}
