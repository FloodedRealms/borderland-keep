package test

import (
	"testing"

	"github.com/floodedrealms/adventure-archivist/types"
)

func setup_characterTests() *types.CharacterRecord {
	return types.NewCharacter(1, 0, 0, 1, "Testy", "test man")

}

func TestXPGainFunctionWithNoPrimeReq(t *testing.T) {
	char := setup_characterTests()
	char.AddXP(1000)
	wanted := 1000
	if char.CurrentXP != wanted {
		t.Errorf("Character failed to gain correct XP with a Prime Requiset was %d. Wanted %d xp, got %d xp", char.PrimeReqPercent, wanted, char.CurrentXP)
	}

}
func TestXPGainFunctionWithFivePrimeReq(t *testing.T) {
	char := setup_characterTests()
	char.PrimeReqPercent = 5
	char.AddXP(1000)
	wanted := 1050
	if char.CurrentXP != wanted {
		t.Errorf("Character failed to gain correct XP with a Prime Requiset was %d. Wanted %d xp, got %d xp", char.PrimeReqPercent, wanted, char.CurrentXP)
	}

}
func TestXPGainFunctionWithTenPrimeReq(t *testing.T) {
	char := setup_characterTests()
	char.PrimeReqPercent = 10
	char.AddXP(1000)
	wanted := 1100
	if char.CurrentXP != wanted {
		t.Errorf("Character failed to gain correct XP with a Prime Requiset was %d. Wanted %d xp, got %d xp", char.PrimeReqPercent, wanted, char.CurrentXP)
	}

}
