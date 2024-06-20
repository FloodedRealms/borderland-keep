package services

import (
	"context"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type CharacterService interface {
	CreateCharacterForCampaign(*types.CampaignRecord, *types.CreateCharacterRecordRequest) (*types.CharacterRecord, error)
	UpdateCharacter(int, *types.UpdateCharacterRecordRequest) (*types.CharacterRecord, error)
	ManageCharactersForAdventure(adventure types.AdventureRecord, character *types.CharacterRecord, operation, halfshare string) (bool, error)
	GetCharactersForCampaign(campaign *types.CampaignRecord) ([]types.CharacterRecord, error)
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

func (s CharacterServiceImpl) CreateCharacterForCampaign(campaign *types.CampaignRecord, charToInsert *types.CreateCharacterRecordRequest) (*types.CharacterRecord, error) {
	character := types.NewCharacterFromCreateRequest(-1, *charToInsert)
	return s.repo.CreateCharacterForCampaign(campaign, character)
}

func (s CharacterServiceImpl) UpdateCharacter(id int, updateReq *types.UpdateCharacterRecordRequest) (*types.CharacterRecord, error) {
	characterToUpdate := types.NewCharacterFromUpdateRequest(id, *updateReq)
	return s.repo.UpdateCharacter(characterToUpdate)
}

func (s CharacterServiceImpl) ManageCharactersForAdventure(ad types.AdventureRecord, char *types.CharacterRecord, operation, halfshare string) (bool, error) {
	isGettingHalfshare := false
	if halfshare != "false" {
		isGettingHalfshare = true
	}
	fullShareXP, halfShareXP := ad.CalculateXPShares()
	switch operation {
	case "add":
		var adventureCharacter *types.AdventureCharacter
		if isGettingHalfshare {
			adventureCharacter = types.NewAdventureCharacter(char, isGettingHalfshare, halfShareXP)
		} else {
			adventureCharacter = types.NewAdventureCharacter(char, isGettingHalfshare, fullShareXP)
		}
		return s.repo.AddCharacterToAdventure(&ad, adventureCharacter)
	case "remove":
		return s.repo.RemoveCharacterFromAdventure(&ad, char)
	case "change-shares":
		return s.repo.ChangeCharacterShares(&ad, char, isGettingHalfshare)
	default:
		return false, util.UnknownCharacterOperation(operation)
	}
}

func (s CharacterServiceImpl) GetCharactersForCampaign(campaign *types.CampaignRecord) ([]types.CharacterRecord, error) {
	characterList, err := s.repo.GetCharactersForCampaign(campaign)
	if err != nil {
		return nil, err
	}
	return characterList, nil
}

func (s CharacterServiceImpl) UpdateTotalCharacterXP(char types.CharacterRecord, xpGained int) error {
	currentChar := s.repo.GetCharacterById(char)
	currentChar.AddXP(xpGained)
	_, err := s.repo.UpdateCharacter(char)
	return err
}
