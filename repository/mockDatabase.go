package repository

import (
	"time"

	"github.com/floodedrealms/adventure-archivist/types"
)

type MockDB struct {
	throwError     bool
	testCampaign   *types.CampaignRecord
	testAdventure  *types.AdventureRecord
	testCharacters []*types.CharacterRecord
}

func NewMockDatabase(shouldError bool) *MockDB {
	testCampaign := types.NewCampaign(1)
	testAdventure := types.NewAdventureRecord(1, 1, 5, *types.NewCoins(1000, 1000, 1000, 1000, 1000), []types.Gem{}, []types.Jewellery{}, []types.MonsterGroup{}, []types.MagicItem{}, []types.AdventureCharacter{}, "Test Adventure", time.Now(), time.Now(), types.ArcvhistDate(time.Now()))
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
/*
func (m *MockDB) CreateCampaign(*types.CampaignRecord) (*types.CampaignRecord, error) {
	if m.throwError {
		return error.N
	}
}*/

/*	UpdateCampaign(*types.CampaignRecord) (*types.CampaignRecord, error)
	GetCampaign(int) (*types.CampaignRecord, error)
	DeleteCampaign(*types.CampaignRecord) (bool, error)
	ListCampaigns() ([]*types.CampaignRecord, error)
	ListCampaignsForClient(string) ([]*types.CampaignRecord, error)

	GetAdventureRecordsForCampaign(*types.CampaignRecord) ([]*types.AdventureRecord, error)
	CreateAdventureRecordForCampaign(*types.CreateAdventureRequest) (*types.AdventureRecord, error)
	GetAdventureRecordById(*types.AdventureRecord) (*types.AdventureRecord, error)
	UpdateCoinsForAdventure(a *types.AdventureRecord, c *types.Coins) (bool, error)

	DeleteGemsForAdventure(a *types.AdventureRecord) error
	DeleteJewelleryForAdventure(a *types.AdventureRecord) error
	DeleteMagicItemsForAdventure(a *types.AdventureRecord) error
	DeleteCombatForAdventure(a *types.AdventureRecord) error
	DeleteCharactersForAdventure(a *types.AdventureRecord) error

	AddGemToAdventure(*types.AdventureRecord, *types.Gem) (bool, error)
	AddJewelleryToAdventure(*types.AdventureRecord, *types.Jewellery) (bool, error)
	AddMagicItemToAdventure(*types.AdventureRecord, *types.MagicItem) (bool, error)
	AddCombatToAdventure(*types.AdventureRecord, *types.MonsterGroup) (bool, error)

	CreateCharacterForCampaign(*types.CampaignRecord, types.Character) (*types.CharacterRecord, error)
	GetCharactersForCampaign(*types.CampaignRecord) ([]types.CharacterRecord, error)
	UpdateCharacter(types.Character) (*types.CharacterRecord, error)
	AddCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter) (bool, error)
	AddHalfshareCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter, int) (bool, error)
	AddFullshareCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter, int) (bool, error)

	RemoveCharacterFromAdventure(*types.AdventureRecord, *types.CharacterRecord) (bool, error)
	ChangeCharacterShares(*types.AdventureRecord, *types.CharacterRecord, bool) (bool, error)
	GetCharacterById(types.CharacterRecord) *types.CharacterRecord

	SaveApiUser(types.User, bool) error
	GetApiUserById(providedClientId, providedAPIKey string) (*types.APIUser, error) */
