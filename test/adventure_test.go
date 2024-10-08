package test

import (
	"testing"

	"github.com/floodedrealms/borderland-keep/types"
)

func fiveShareGroup() []types.AdventureCharacter {
	return []types.AdventureCharacter{
		{Halfshare: true},
		{Halfshare: true},
		{Halfshare: true},
		{Halfshare: true},
		{Halfshare: false},
		{Halfshare: false},
		{Halfshare: false},
	}
}
func fivePointFiveShareGroup() []types.AdventureCharacter {
	return []types.AdventureCharacter{
		{Halfshare: true},
		{Halfshare: true},
		{Halfshare: true},
		{Halfshare: true},
		{Halfshare: true},
		{Halfshare: false},
		{Halfshare: false},
		{Halfshare: false},
	}
}

func coins() types.Coins {
	return types.Coins{
		Copper:   *types.NewCopper(1000),
		Silver:   *types.NewSilver(1000),
		Electrum: *types.NewElectrum(1000),
		Gold:     *types.NewGold(1000),
		Platinum: *types.NewPlatinum(1000),
	}
}

func gems() []types.Gem {
	return []types.Gem{
		{XPValue: 1000, Number: 1},
	}
}

func jewellery() []types.Jewellery {
	return []types.Jewellery{
		{XPValue: 500, Number: 2},
	}
}
func magicItems() []types.MagicItem {
	return []types.MagicItem{
		{ApparentValue: 250},
		{ApparentValue: 250},
		{ApparentValue: 500},
	}
}
func combat() []types.MonsterGroup {
	return []types.MonsterGroup{
		{XPPerOneKill: 10, NumberDefeated: 10},
		{XPPerOneKill: 15, NumberDefeated: 4},
		{XPPerOneKill: 20, NumberDefeated: 2},
		{XPPerOneKill: 800, NumberDefeated: 1},
	}
}
func tenThousandAdventureFiveShares() types.AdventureRecord {
	return types.AdventureRecord{
		Coins:      coins(),
		Gems:       gems(),
		MagicItems: magicItems(),
		Jewellery:  jewellery(),
		Combat:     []types.MonsterGroup{{XPPerOneKill: 90, NumberDefeated: 1}, {XPPerOneKill: 300, NumberDefeated: 1}},
		Characters: fiveShareGroup(),
	}
}
func tenThousandAdventureFivePointFiveShares() types.AdventureRecord {
	return types.AdventureRecord{
		Coins:      coins(),
		Gems:       gems(),
		MagicItems: magicItems(),
		Jewellery:  jewellery(),
		Combat:     []types.MonsterGroup{{XPPerOneKill: 90, NumberDefeated: 1}, {XPPerOneKill: 300, NumberDefeated: 1}},
		Characters: fivePointFiveShareGroup(),
	}
}

func hobGoblinSlayer() types.AdventureRecord {
	return types.AdventureRecord{
		Coins:     *types.NewCoins(0, 150, 0, 649, 0),
		Jewellery: []types.Jewellery{{Name: "Bracelet", Number: 2, XPValue: 135}},
		Combat: []types.MonsterGroup{
			{Name: "Goblin", XPPerOneKill: 10, NumberDefeated: 5},
			{Name: "Hobhoblin", XPPerOneKill: 15, NumberDefeated: 7},
			{Name: "Orc", XPPerOneKill: 10, NumberDefeated: 1},
			{Name: "Gnoll", XPPerOneKill: 20, NumberDefeated: 1},
			{Name: "Orge", XPPerOneKill: 140, NumberDefeated: 1},
		},
		Characters: []types.AdventureCharacter{
			{Halfshare: false},
			{Halfshare: false},
			{Halfshare: false},
			{Halfshare: false},
		},
	}
}

func pissClothesMysters() types.AdventureRecord {
	return types.AdventureRecord{
		Coins: *types.NewCoins(0, 5000, 11000, 34, 2),
		Jewellery: []types.Jewellery{
			{Name: "Crown", Number: 1, XPValue: 4000},
			{Name: "Necklace", Number: 1, XPValue: 300},
			{Name: "Cup", Number: 1, XPValue: 90},
			{Name: "Tapestry", Number: 1, XPValue: 900},
		},
		Combat: []types.MonsterGroup{
			{Name: "Goblin", XPPerOneKill: 10, NumberDefeated: 1},
			{Name: "Hobgoblin", XPPerOneKill: 15, NumberDefeated: 3},
		},
		Characters: []types.AdventureCharacter{
			{Halfshare: false},
			{Halfshare: false},
			{Halfshare: false},
			{Halfshare: false},
			{Halfshare: true},
			{Halfshare: true},
		},
	}
}

func TestXPCalculationsForAdventures(t *testing.T) {

	var tests = []struct {
		name                   string
		expectedTotalXp        int
		expectedFullXpShare    int
		expectedHalfXpShare    int
		expextedNumberOfShares float64
		adventure              types.AdventureRecord
	}{
		{"Adventure with 1000 xp of each coin and 5 shares", 6610, 1322, 661, 5.0, types.AdventureRecord{Coins: coins(), Characters: fiveShareGroup()}},
		{"Adventure with 1000 xp of each coin and 5.5 shares", 6610, 1202, 601, 5.5, types.AdventureRecord{Coins: coins(), Characters: fivePointFiveShareGroup()}},
		{"Adventure with 1000 XP of Gems and 5 shares", 1000, 200, 100, 5.0, types.AdventureRecord{Gems: gems(), Characters: fiveShareGroup()}},
		{"Adventure with 1000 XP of Gems and 5.5 shares", 1000, 182, 91, 5.5, types.AdventureRecord{Gems: gems(), Characters: fivePointFiveShareGroup()}},
		{"Adventure with 1000 XP of Jewellery and 5 shares", 1000, 200, 100, 5.0, types.AdventureRecord{Jewellery: jewellery(), Characters: fiveShareGroup()}},
		{"Adventure with 1000 XP of Jewellery and 5.5 shares", 1000, 182, 91, 5.5, types.AdventureRecord{Jewellery: jewellery(), Characters: fivePointFiveShareGroup()}},
		{"Adventure with 1000 XP of Magic Items and 5 shares", 1000, 200, 100, 5.0, types.AdventureRecord{MagicItems: magicItems(), Characters: fiveShareGroup()}},
		{"Adventure with 1000 XP of Magic Items and 5.5 shares", 1000, 182, 91, 5.5, types.AdventureRecord{MagicItems: magicItems(), Characters: fivePointFiveShareGroup()}},
		{"Adventure with 1000 XP of Combat and 5 shares", 1000, 200, 100, 5.0, types.AdventureRecord{Combat: combat(), Characters: fiveShareGroup()}},
		{"Adventure with 1000 XP of Combat and 5.5 shares", 1000, 182, 91, 5.5, types.AdventureRecord{Combat: combat(), Characters: fivePointFiveShareGroup()}},
		{"Adventure with 10000 XP Total and 5 shares", 10000, 2000, 1000, 5.0, tenThousandAdventureFiveShares()},
		{"Adventure with 10000 XP Total and 5.5 shares", 10000, 1818, 909, 5.5, tenThousandAdventureFivePointFiveShares()},
		{"[Hob]Goblin Slayer", 1259, 315, 158, 4.0, hobGoblinSlayer()},
		{"Piss Clothes", 11389, 2278, 1139, 5.0, pissClothesMysters()},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.adventure.CalculateNumberOfShares()
			shares := tt.adventure.NumberOfShares
			totalXp := tt.adventure.TotalXPAmount()
			fullShare, halfShare := tt.adventure.CalculateXPShares()
			if shares != tt.expextedNumberOfShares {
				t.Errorf("Expected %f shares, got %f", tt.expextedNumberOfShares, shares)
			}
			if totalXp != tt.expectedTotalXp {
				t.Errorf("Expected %d total XP, got %d", tt.expectedTotalXp, totalXp)
			}
			if halfShare != tt.expectedHalfXpShare {
				t.Errorf("Expected %d for half share, got %d", tt.expectedHalfXpShare, halfShare)
			}
			if fullShare != tt.expectedFullXpShare {
				t.Errorf("Expected %d total XP, got %d", tt.expectedFullXpShare, fullShare)
			}

		})
	}

}

func TestSimpleCalculatorDetailOnFirstOpen(t *testing.T) {
	adventure := *types.NewAdventureRecordByNumberOfCharacters(0, 0)
	expectedTotalXP := 0
	actualTotalXP := adventure.DisplayTotalXPAmount()
	if expectedTotalXP != actualTotalXP {
		t.Errorf("Expected %d shares, got %d shares instead.", expectedTotalXP, actualTotalXP)
	}

	expectedShares := 0.0
	actualShares := adventure.NumberOfShares
	if expectedShares != actualShares {
		t.Errorf("Expected %f shares, got %f shares instead.", expectedShares, actualShares)
	}

	expectedDisplayShares := "0.0"
	actualDisplayShares := adventure.DisplayTotalShares()
	if expectedDisplayShares != actualDisplayShares {
		t.Errorf("Expected %s shares, got %s shares instead.", expectedDisplayShares, actualDisplayShares)
	}

	expectFullXPShare := 0
	actualFullXPShare := adventure.DisplayFullXPShare()
	if expectFullXPShare != actualFullXPShare {
		t.Errorf("Expect a full XP share to be %d, got %d", expectFullXPShare, actualFullXPShare)
	}

	expectHalfXPShare := 0
	actualHalfXPShare := adventure.DisplayHalfXPShare()
	if expectHalfXPShare != actualHalfXPShare {
		t.Errorf("Expect a full XP share to be %d, got %d", expectHalfXPShare, actualHalfXPShare)
	}

	expectFullGPShare := 0
	actualFullGPShare := adventure.DisplayFullGPShare()
	if expectFullGPShare != actualFullGPShare {
		t.Errorf("Expect a full GP share to be %d, got %d", expectFullGPShare, actualFullGPShare)
	}

	expectHalfGPShare := 0
	actualHalfGPShare := adventure.DisplayHalfGPShare()
	if expectHalfGPShare != actualHalfGPShare {
		t.Errorf("Expect a full GP share to be %d, got %d", expectHalfGPShare, actualHalfGPShare)
	}

}

func TestTotalsWithZeroShares(t *testing.T) {
	adventure := *types.NewAdventureRecordByNumberOfCharacters(0, 0)
	adventure.Coins = *types.NewCoins(0, 0, 0, 400, 0)
	expectedTotalXP := 400
	actualTotalXP := adventure.DisplayTotalXPAmount()
	if expectedTotalXP != actualTotalXP {
		t.Errorf("Expected %d XP Total, got %d XP Total instead.", expectedTotalXP, actualTotalXP)
	}

	expectedTotalGp := 400
	actualTotalGp := adventure.DisplayTotalGPAmount()
	if expectedTotalGp != actualTotalGp {
		t.Errorf("Expected %d GP Total, got %d GP total instead.", expectedTotalGp, actualTotalGp)
	}

	expectedShares := 0.0
	actualShares := adventure.NumberOfShares
	if expectedShares != actualShares {
		t.Errorf("Expected %f shares, got %f shares instead.", expectedShares, actualShares)
	}

	expectedDisplayShares := "0.0"
	actualDisplayShares := adventure.DisplayTotalShares()
	if expectedDisplayShares != actualDisplayShares {
		t.Errorf("Expected %s shares, got %s shares instead.", expectedDisplayShares, actualDisplayShares)
	}

	expectFullXPShare := 0
	actualFullXPShare := adventure.DisplayFullXPShare()
	if expectFullXPShare != actualFullXPShare {
		t.Errorf("Expect a full XP share to be %d, got %d", expectFullXPShare, actualFullXPShare)
	}

	expectHalfXPShare := 0
	actualHalfXPShare := adventure.DisplayHalfXPShare()
	if expectHalfXPShare != actualHalfXPShare {
		t.Errorf("Expect a full XP share to be %d, got %d", expectHalfXPShare, actualHalfXPShare)
	}

	expectFullGPShare := 0
	actualFullGPShare := adventure.DisplayFullGPShare()
	if expectFullGPShare != actualFullGPShare {
		t.Errorf("Expect a full GP share to be %d, got %d", expectFullGPShare, actualFullGPShare)
	}

	expectHalfGPShare := 0
	actualHalfGPShare := adventure.DisplayHalfGPShare()
	if expectHalfGPShare != actualHalfGPShare {
		t.Errorf("Expect a full GP share to be %d, got %d", expectHalfGPShare, actualHalfGPShare)
	}

}
