package repository

import (
	"time"

	"github.com/floodedrealms/borderland-keep/types"
)

type MockDB struct {
	throwError     bool
	testCampaign   *types.CampaignRecord
	testAdventure  *types.AdventureRecord
	testCharacters []*types.CharacterRecord
}

func NewMockDatabase(shouldError bool) *MockDB {
	testCampaign := types.NewCampaign(1)
	testAdventure := types.NewAdventureRecord(1, 1, 5, *types.NewCoins(1000, 1000, 1000, 1000, 1000), []types.Gem{}, []types.Jewellery{}, []types.MonsterGroup{}, []types.MagicItem{}, []types.AdventureCharacter{}, "Test Adventure", types.ArcvhistDate(time.Now()))
	character1 := types.NewCharacter(1, 0, 5, 1, "Billy the Test", "Figher")
	character2 := types.NewCharacter(2, 0, 5, 1, "Willy the Test", "Figher")
	character3 := types.NewCharacter(3, 0, 5, 1, "Chilly the Test", "Figher")
	return &MockDB{
		throwError:     shouldError,
		testCampaign:   testCampaign,
		testAdventure:  testAdventure,
		testCharacters: []*types.CharacterRecord{character1, character2, character3},
	}

}

// PerformInsertStatement(*types.SQLLiteExportable) (*types.JSONExportable, error)
