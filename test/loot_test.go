package test

import (
	"testing"

	"github.com/floodedrealms/adventure-archivist/types"
)

func TestCoins(t *testing.T) {
	var tests = []struct {
		name            string
		coinType        string
		numberOfCoins   int
		expectedXPValue float64
	}{
		{"100 Copper Should be 1 XP", "copper", 100, 1.0},
		{"1000 Copper Should be 10 XP", "copper", 1000, 10.0},
		{"1500 Copper Should be 15 XP", "copper", 1500, 15.0},

		{"100 Silver Should be 10 XP", "silver", 100, 10.0},
		{"1000 Silver Should be 100 XP", "silver", 1000, 100.0},
		{"1500 Silver Should be 150 XP", "silver", 1500, 150.0},

		{"100 Electrum Should be 50 XP", "electrum", 100, 50.0},
		{"1000 Electrum Should be 500 XP", "electrum", 1000, 500.0},
		{"1500 Electrum Should be 750 XP", "electrum", 1500, 750.0},

		{"100  Gold Should be 100 XP", "gold", 100, 100.0},
		{"1000 Gold Should be 1000 XP", "gold", 1000, 1000.0},
		{"1500 Gold Should be 1500 XP", "gold", 1500, 1500.0},

		{"100 Platinum Should be 500 XP", "platinum", 100, 500.0},
		{"1000 Platinum Should be 5000 XP", "platinum", 1000, 5000.0},
		{"1500 Platinum Should be 7500 XP", "platinum", 1500, 7500.0},
	}
	for _, test := range tests {
		t.Run(test.name, getCoinTestFunc(test.coinType, test.expectedXPValue, test.numberOfCoins))

	}
}

func getCoinTestFunc(coin string, expectedXPValue float64, numberOfCoins int) func(*testing.T) {
	switch coin {
	case "copper":
		c := types.NewCopper(numberOfCoins)
		return func(t *testing.T) {
			xp := c.TotalXPAmount()
			if xp != expectedXPValue {
				t.Errorf("Got %f XP for %d %s, want %f xp", xp, numberOfCoins, coin, expectedXPValue)
			}
		}
	case "silver":
		c := types.NewSilver(numberOfCoins)
		return func(t *testing.T) {
			xp := c.TotalXPAmount()
			if xp != expectedXPValue {
				t.Errorf("Got %f XP for %d %s, want %f xp", xp, numberOfCoins, coin, expectedXPValue)
			}
		}
	case "electrum":
		c := types.NewElectrum(numberOfCoins)
		return func(t *testing.T) {
			xp := c.TotalXPAmount()
			if xp != expectedXPValue {
				t.Errorf("Got %f XP for %d %s, want %f xp", xp, numberOfCoins, coin, expectedXPValue)
			}
		}
	case "gold":
		c := types.NewGold(numberOfCoins)
		return func(t *testing.T) {
			xp := c.TotalXPAmount()
			if xp != expectedXPValue {
				t.Errorf("Got %f XP for %d %s, want %f xp", xp, numberOfCoins, coin, expectedXPValue)
			}
		}
	case "platinum":
		c := types.NewPlatinum(numberOfCoins)
		return func(t *testing.T) {
			xp := c.TotalXPAmount()
			if xp != expectedXPValue {
				t.Errorf("Got %f XP for %d %s, want %f xp", xp, numberOfCoins, coin, expectedXPValue)
			}
		}
	}
	return func(t *testing.T) {
		t.Errorf("Coint type %s not testable!", coin)
	}
}
