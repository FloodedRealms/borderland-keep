package types

type XPSource interface {
	TotalXPAmount() int
}

type GoldSource interface {
	TotalValueInGold() int
}

type XPSourceContainer interface {
	GenerateCoins() Coins
	GenerateGemList() []Gem
	GenerateJewelleryList() []Jewellery
	GenerateMagicItemList() []MagicItem
	GenerateCombatList() []MonsterGroup
}
type APIResponse interface {
}
type APIObject interface {
	GenerateSuccessfulCreationJSON() APIResponse
}

type SQLLiteExportable interface {
	Id() int
	GenerateUpdateStatement() string
}

type Adventure interface {
	XPSource
	XPSourceContainer
	GenerateCharacterToAdventureList() []AdventureCharacter
	HalfShareXP() int
	FullShareXP() int
}
