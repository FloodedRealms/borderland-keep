package services

import (
	"context"

	"github.com/floodedrealms/adventure-archivist/repository"
	"github.com/floodedrealms/adventure-archivist/types"
	"github.com/floodedrealms/adventure-archivist/util"
)

type CharacterService interface {
	CreateCharacterForCampaign(campaign *types.Campaign) (*types.Character, error)
	ManageCharactersForAdventure(adventure *types.Adventure, character *types.Character, operation, halfshare string) (bool, error)
	GetCharactersForCampaign(campaign *types.Campaign) ([]types.Character, error)
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

func (s CharacterServiceImpl) CreateCharacterForCampaign(campaign *types.Campaign) (*types.Character, error) {
	return s.repo.CreateCharacterForCampaign(campaign)
}

func (s CharacterServiceImpl) ManageCharactersForAdventure(ad *types.Adventure, char *types.Character, operation, halfshare string) (bool, error) {
	isGettingHalfshare := false
	if halfshare != "false" {
		isGettingHalfshare = true
	}
	switch operation {
	case "add":
		return s.repo.AddCharacterToAdventure(ad, char, isGettingHalfshare)
	case "remove":
		return s.repo.RemoveCharacterFromAdventure(ad, char)
	case "change-shares":
		return s.repo.ChangeCharacterShares(ad, char, isGettingHalfshare)
	default:
		return false, util.UnknownCharacterOperation(operation)
	}
}

func (s CharacterServiceImpl) GetCharactersForCampaign(campaign *types.Campaign) ([]types.Character, error) {
	characterList, err := s.repo.GetCharactersForCampaign(campaign)
	if err != nil {
		return nil, err
	}
	return characterList, nil
}
