package repository

import (
	"database/sql"
	"errors"

	"github.com/floodedrealms/borderland-keep/internal/util"
	"github.com/floodedrealms/borderland-keep/types"
)

type JSONRepo struct {
	Campaign   []types.CampaignRecord  `json:"campaigns"`
	Adventures []types.AdventureRecord `json:"adventures"`
	Characters []types.CharacterRecord `json:"characters"`
}

func NewJSONRepo(filename string) *JSONRepo {

	var repo JSONRepo
	err := decodeJSONFile(filename, &repo)
	if err != nil {
		panic(err)
	}
	return &repo
}

func (j JSONRepo) ExecuteQuery(stmt string, params ...interface{}) (sql.Result, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) RunQuery(stmt string, params ...interface{}) (*sql.Rows, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) GetCoinsForAdventure(a *types.AdventureRecord) (*types.Coins, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) CreateCampaign(*types.CampaignRecord) (*types.CampaignRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) UpdateCampaign(*types.CampaignRecord) (*types.CampaignRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) GetCampaign(int) (*types.CampaignRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) DeleteCampaign(*types.CampaignRecord) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) ListCampaigns() ([]*types.CampaignRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) ListCampaignsForClient(string) ([]*types.CampaignRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) GetAdventureRecordsForCampaign(*types.CampaignRecord) ([]*types.AdventureRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) CreateAdventureRecordForCampaign(*types.AdventureRecord) (*types.AdventureRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) GetAdventureRecordById(in *types.AdventureRecord) (*types.AdventureRecord, error) {
	for _, a := range j.Adventures {
		if a.Id == in.Id {
			return &a, nil
		}
	}
	return nil, errors.New("id not found")
}

func (j JSONRepo) UpdateCoinsForAdventure(a *types.AdventureRecord, c *types.Coins) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) UpdateAdventureName(*types.AdventureRecord, string) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) UpdateAdventureDate(*types.AdventureRecord, types.ArcvhistDate) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) DeleteGemsForAdventure(a *types.AdventureRecord) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) DeleteJewelleryForAdventure(a *types.AdventureRecord) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) DeleteMagicItemsForAdventure(a *types.AdventureRecord) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) DeleteCombatForAdventure(a *types.AdventureRecord) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) DeleteCharactersForAdventure(a *types.AdventureRecord) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) AddGemToAdventure(*types.AdventureRecord, *types.Gem) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) AddJewelleryToAdventure(*types.AdventureRecord, *types.Jewellery) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) AddMagicItemToAdventure(*types.AdventureRecord, *types.MagicItem) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) AddCombatToAdventure(*types.AdventureRecord, *types.MonsterGroup) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) CreateCharacterForCampaign(*types.CampaignRecord, types.Character) (*types.CharacterRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) GetCharactersForCampaign(*types.CampaignRecord) ([]types.CharacterRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) UpdateCharacter(types.CharacterRecord) (*types.CharacterRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) AddCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) AddHalfshareCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter, int) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) AddFullshareCharacterToAdventure(*types.AdventureRecord, *types.AdventureCharacter, int) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) RemoveCharacterFromAdventure(*types.AdventureRecord, *types.CharacterRecord) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) ChangeCharacterShares(*types.AdventureRecord, *types.CharacterRecord, bool) (bool, error) {
	return false, util.NotYetImplmented()
}

func (j JSONRepo) GetCharacterById(types.CharacterRecord) (*types.CharacterRecord, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) UpdateCharacterTotalXP(types.CharacterRecord) error {
	return util.NotYetImplmented()
}

func (j JSONRepo) GetCharacterXPGains(types.CharacterRecord) ([]int, error) {
	return nil, util.NotYetImplmented()
}

func (j JSONRepo) GetLevelForXP(types.CampaignRecord, types.CharacterRecord) int {
	return 1
}

func (j JSONRepo) AddCampaignActivityForCharacter(types.CampaignActivity) error {
	return util.NotYetImplmented()
}
