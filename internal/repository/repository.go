package repository

import (
	"database/sql"

	"github.com/floodedrealms/adventure-archivist/types"
)

type Repository interface {
	//PerformInsertStatement(*types.SQLLiteExportable) (*types.JSONExportable, error)
	CreateCampaign(*types.CampaignRecord) (*types.CampaignRecord, error)
	UpdateCampaign(*types.CampaignRecord) (*types.CampaignRecord, error)
	GetCampaign(int) (*types.CampaignRecord, error)
	DeleteCampaign(*types.CampaignRecord) (bool, error)
	ListCampaigns() ([]*types.CampaignRecord, error)
	ListCampaignsForClient(string) ([]*types.CampaignRecord, error)
	UpdateCampaignPassword(int, types.Password) error

	GetAdventureRecordsForCampaign(*types.CampaignRecord) ([]*types.AdventureRecord, error)
	CreateAdventureRecordForCampaign(*types.AdventureRecord) (*types.AdventureRecord, error)
	GetAdventureRecordById(*types.AdventureRecord) (*types.AdventureRecord, error)
	UpdateCoinsForAdventure(a *types.AdventureRecord, c *types.Coins) (bool, error)
	UpdateAdventureName(*types.AdventureRecord, string) error
	UpdateAdventureDate(*types.AdventureRecord, types.ArcvhistDate) error
	GetCoinsForAdventure(*types.AdventureRecord) (*types.Coins, error)

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
	//GetCharactersForAdventure(types.AdventureRecord) ([]types.CharacterRecord, error)
	UpdateCharacter(types.CharacterRecord) (*types.CharacterRecord, error)
	AddCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter) (bool, error)
	AddHalfshareCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter, int) (bool, error)
	AddFullshareCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter, int) (bool, error)

	RemoveCharacterFromAdventure(*types.AdventureRecord, *types.CharacterRecord) (bool, error)
	ChangeCharacterShares(*types.AdventureRecord, *types.CharacterRecord, bool) (bool, error)
	GetCharacterById(types.CharacterRecord) (*types.CharacterRecord, error)
	UpdateCharacterTotalXP(types.CharacterRecord) error
	GetCharacterXPGains(types.CharacterRecord) ([]int, error)
	GetLevelForXP(types.CampaignRecord, types.CharacterRecord) int
	AddCampaignActivityForCharacter(types.CampaignActivity) error

	SaveApiUser(types.User, bool) error
	GetApiUserById(providedClientId, providedAPIKey string) (*types.APIUser, error)

	ExecuteQuery(q string, params ...interface{}) (sql.Result, error)
	RunQuery(q string, params ...interface{}) (*sql.Rows, error)
}
