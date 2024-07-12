package services

import (
	"context"

	"github.com/floodedrealms/adventure-archivist/internal/repository"
	"github.com/floodedrealms/adventure-archivist/internal/util"
	"github.com/floodedrealms/adventure-archivist/types"
)

type CharacterService interface {
	CreateCharacterForCampaign(*types.CampaignRecord, *types.CharacterRecord) (*types.CharacterRecord, error)
	UpdateCharacter(int, *types.CharacterRecord) (*types.CharacterRecord, error)
	//ManageCharactersForAdventure(adventure types.AdventureRecord, character *types.CharacterRecord, operation, halfshare string) (bool, error)
	GetCharactersForCampaign(campaign *types.CampaignRecord) ([]types.CharacterRecord, error)
	//AddAdventureXPToCharacters(types.AdventureRecord) ([]types.CharacterRecord, error)
}

type CharacterServiceImpl struct {
	repo   repository.Repository
	logger *util.Logger
	ctx    context.Context
}

func NewCharacterService(repo repository.Repository, logger *util.Logger, ctx context.Context) *CharacterServiceImpl {
	return &CharacterServiceImpl{
		repo:   repo,
		logger: logger,
		ctx:    ctx,
	}
}

func (s CharacterServiceImpl) CreateCharacterForCampaign(campaign *types.CampaignRecord, charToInsert *types.CharacterRecord) (*types.CharacterRecord, error) {
	return s.repo.CreateCharacterForCampaign(campaign, charToInsert)
}

func (s CharacterServiceImpl) UpdateCharacter(id int, char *types.CharacterRecord) (*types.CharacterRecord, error) {
	return s.repo.UpdateCharacter(*char)
}

/*func (s CharacterServiceImpl) ManageCharactersForAdventure(ad types.AdventureRecord, char *types.CharacterRecord, operation, halfshare string) (bool, error) {
	isGettingHalfshare := false
	if halfshare != "false" {
		isGettingHalfshare = true
	}
	fullShareXP, halfShareXP := ad.CalculateXPShares()
	switch operation {
	case "add":
		var adventureCharacter *types.AdventureCharacter
		if isGettingHalfshare {
			adventureCharacter = types.NewAdventureCharacter(isGettingHalfshare, char.ID)
		} else {
			adventureCharacter = types.NewAdventureCharacter(isGettingHalfshare, char.ID)
		}
		return s.repo.AddCharacterToAdventure(&ad, adventureCharacter)
	case "remove":
		return s.repo.RemoveCharacterFromAdventure(&ad, char)
	case "change-shares":
		return s.repo.ChangeCharacterShares(&ad, char, isGettingHalfshare)
	default:
		return false, util.UnknownCharacterOperation(operation)
	}
}*/

func (s CharacterServiceImpl) GetCharactersForCampaign(campaign *types.CampaignRecord) ([]types.CharacterRecord, error) {
	characterList, err := s.repo.GetCharactersForCampaign(campaign)
	if err != nil {
		return nil, err
	}
	for i, c := range characterList {
		allXp, err := s.repo.GetCharacterXPGains(c)
		if err != nil {
			return nil, err
		}
		characterList[i].CurrentXP = s.sumXP(allXp)
		characterList[i].Level = s.repo.GetLevelForXP(*campaign, characterList[i])
	}

	return characterList, nil
}

func (s CharacterServiceImpl) GetPossibleCharactersForAdventure(aId int) ([]types.CharacterRecord, error) {
	chars := make([]types.CharacterRecord, 0)
	stmtString := ("SELECT c.id, c.name from adventures a " +
		"JOIN \"characters\" c ON c.campaign_id = a.campaign_id " +
		"WHERE a.id = ?;")

	result, err := s.repo.RunQuery(stmtString, aId)
	if err != nil {
		return chars, err
	}
	for result.Next() {
		c := types.CharacterRecord{}
		result.Scan(&c.Id, &c.Name)
		chars = append(chars, c)
	}
	return chars, nil
}

func (s CharacterServiceImpl) UpdateTotalCharacterXP(char types.CharacterRecord, xpGained int) (*types.CharacterRecord, error) {
	currentChar, err := s.repo.GetCharacterById(char)
	if err != nil {
		return nil, err
	}
	currentChar.AddXP(xpGained)
	c, err := s.repo.UpdateCharacter(char)
	if err != nil {
		return nil, err
	}
	return c, err
}

func (s CharacterServiceImpl) sumXP(a []int) int {
	xp := 0
	for _, x := range a {
		xp += x
	}
	return xp
}
