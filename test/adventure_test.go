package test

import (
	"testing"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
)

func setup_adventureTests() repository.Repository {
	return repository.NewJSONRepo("/json-data.json")

}

func TestXPCalculationForAdventures(t *testing.T) {
	repo := setup_adventureTests()
	a1, _ := repo.GetAdventureRecordById(types.NewAdventureRecordById(1))
	testAdventureXP(314, 157, *a1, t)
	a2, _ := repo.GetAdventureRecordById(types.NewAdventureRecordById(2))
	testAdventureXP(2278, 1138, *a2, t)
	a3, _ := repo.GetAdventureRecordById(types.NewAdventureRecordById(3))
	testAdventureXP(244, 122, *a3, t)
	a4, _ := repo.GetAdventureRecordById(types.NewAdventureRecordById(4))
	testAdventureXP(381, 181, *a4, t)

}

func testAdventureXP(expectedFullShare, expectedHalfShare int, a types.AdventureRecord, t *testing.T) {
	givenFullshare, givenHalfshare := a.CalculateXPShares()
	if givenFullshare != expectedFullShare {
		t.Errorf("Fullshare Calculation incorrect for adventure %s. \n\tWanted: %d, Got: %d", a.Name, expectedFullShare, givenFullshare)
	}
	if givenHalfshare != expectedHalfShare {
		t.Errorf("Halfshare Calculation incorrect for adventure %s. \n\tWanted: %d, Got: %d", a.Name, expectedHalfShare, givenHalfshare)
	}
}
